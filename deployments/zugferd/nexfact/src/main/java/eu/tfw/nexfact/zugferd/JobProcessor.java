// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package eu.tfw.nexfact.zugferd;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Map;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import org.apache.pdfbox.Loader;
import org.mustangproject.validator.ZUGFeRDValidator;
import org.mustangproject.ZUGFeRD.ZUGFeRDExporterFromA3;
import org.mustangproject.ZUGFeRD.ZUGFeRDImporter;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDDocumentNameDictionary;
import org.apache.pdfbox.pdmodel.PDEmbeddedFilesNameTreeNode;
import org.apache.pdfbox.pdmodel.common.filespecification.PDComplexFileSpecification;
import org.apache.pdfbox.pdmodel.common.filespecification.PDEmbeddedFile;
import org.mustangproject.ZUGFeRD.ValidationLogVisualizer;
import org.tinylog.Logger;
import org.w3c.dom.Document;

/*
JobProcessor prcess one Job
Job Parameters are taken from a job.json-file: see JobConfig.class
every Paths MUST be in the scope of container paths
*/
public class JobProcessor {

    private Path jobArtefacts;
            
    public Status process(JobConfig job, Path jobDir) throws Exception {

        jobArtefacts = jobDir.toAbsolutePath();

        Logger.debug("Job Mode: {}", job.getQueue());
        
        switch (job.getQueue()) {
            case generate_zf -> {
                return combine(job);
            }
            case combine -> {
                return combine(job);
            }
            case validate -> {
                return validate(job);
            }
            case extract -> {
                return extract(job);
            }
            default -> throw new IllegalArgumentException("Unsupported mode");
        }
    }

    // Combine takes a pdf-file and merge a factur-x.xml an attachments in it
    // input:  job.pdfPath, job.fxPath, evtl. job.attachments
    // output: hybrid-zugferd-pdf: artefactsRoot/jobID/jobID.zf.pdf
    private Status combine(JobConfig job) {
        // Default Werte setzen
        Status status = Status.ok(job.getQueue().name());
        Path inputPdf = jobArtefacts.resolve((job.getId()+".pdf")).toAbsolutePath();
        Path inputXml = jobArtefacts.resolve((job.getId()+".xml")).toAbsolutePath();

        // Update mit job-werten
        if (job.getPdfPath() != null && !job.getPdfPath().isEmpty()) {
            inputPdf = jobArtefacts.resolve(job.getPdfPath()).toAbsolutePath();
        }
        if (job.getFxPath() != null && !job.getFxPath().isEmpty()) {
            inputXml = jobArtefacts.resolve(job.getFxPath()).toAbsolutePath();
        }

        String baseName = getBaseName(inputPdf.getFileName().toString());
        String outputPdf = jobArtefacts.resolve((baseName + ".zf.pdf")).toString();
        
         if (!Files.exists(inputPdf)) {
            return status.addError("PDF_NOT_FOUND", ("Input PDF: " +  inputPdf.toString() + " not found"));
        } 
        if (!Files.exists(inputXml)) {
            return status.addError("XML_NOT_FOUND", "Input XML: " +  inputXml.toString() + " not found");
        }

        Logger.info("trying to combine {} and {} into {}", inputPdf.toString(), inputXml.toString(), outputPdf);
        try {
            // Sollte bereits PDF/A-3 sein
            // 1. Exporter mit PDF/A-3 initialisieren
            ZUGFeRDExporterFromA3 exporter = new ZUGFeRDExporterFromA3();
            exporter.setProducer("NexFact ZFprocessor")
                    .setCreator("NexFact")
                    .setZUGFeRDVersion(2) // Factur-X / ZUGFeRD 2.x
                    .load(Files.readAllBytes(inputPdf));
                    
            // 2. Optionale Anhänge (0:n) hinzufügen Reihenfolge ist wichtig!! attachments VOR xml
            List<String> attachmentPaths = job.getAttachments(); // Annahme: Existiert in JobConfig
            if (attachmentPaths != null && !attachmentPaths.isEmpty()) {
                Logger.info("trying to add additional attachments into {}", outputPdf);
                for (String pathStr : attachmentPaths) {
                    String aPathString = jobArtefacts.resolve(pathStr).toAbsolutePath().toString();
                    File fileToAttach = new File(aPathString);

                    if (!fileToAttach.exists()) {
                        status.getWarnings().add(new StatusMessage("ATTACHMENT_MISSING", "Skipped missing attachment: " + pathStr));
                        continue; // Nächstes File probieren
                    }

                    byte[] fileData = Files.readAllBytes(Paths.get(aPathString));
                    String mimeType = determineMimeType(aPathString);
                    
                    // Anhang in den PDF/A-3 Container injizieren
                    exporter.attachFile(fileToAttach.getName(), fileData, mimeType, "Supplement");
                    Logger.info("added attachments {}", fileToAttach.getName());
                }
            }
   
            // 3. Das bereits validierte XML anhängen
            byte[] xmlData = Files.readAllBytes(inputXml);
            exporter.setXML(xmlData);

            // 4. Das fertige Factur-X/ZUGFeRD Hybrid-PDF speichern
            exporter.export(outputPdf);
            status.getMeta().put("output_file", outputPdf);
            status.getMeta().put("attachments_processed", attachmentPaths == null ? 0 : attachmentPaths.size());

        } catch (IOException e) {
            status.setStatus(Status.STATUS_ERROR);
            status.addError("COMBINE_FAILED", "Failed to combine PDF and XML: " + e.getMessage());
        }

        return status;
    }
    
    // Kleine Hilfsmethode für die wichtigsten Mime-Types in Rechnungsanhängen
    private String determineMimeType(String path) {
        String lowerPath = path.toLowerCase();
        if (lowerPath.endsWith(".pdf")) return "application/pdf";
        if (lowerPath.endsWith(".xml")) return "application/xml";
        if (lowerPath.endsWith(".csv")) return "text/csv";
        if (lowerPath.endsWith(".txt")) return "text/plain";
        if (lowerPath.endsWith(".png")) return "image/png";
        if (lowerPath.endsWith(".jpg") || lowerPath.endsWith(".jpeg")) return "image/jpeg";
        // Fallback for unknow (allowed from PDF/A-3)
        return "application/octet-stream"; 
    }

    // Validate takes a pdfFile or a factur-x.xml -File
    // input: job.pdfPath or job.fxPath
    //     if pdf then validate extracts the factur-x.xml then validate the factur-x.xml
    // output: artefactsRoot/jobID/jobID.result.pdf
    private Status validate(JobConfig job) throws Exception {
        
        // paths: Test for pdf or xml. if both exists: xml will win
        Status status = Status.ok(job.getQueue().name());
        Path fileToValidate = null;
        // test for pdf
        Path inputPdf = jobArtefacts.resolve((job.getId()+".pdf")).toAbsolutePath();
        if (job.getPdfPath() != null && !job.getPdfPath().isEmpty()) {
            inputPdf = jobArtefacts.resolve(job.getPdfPath()).toAbsolutePath();
        }
        if (Files.exists(inputPdf)) {
            fileToValidate = inputPdf;
        } 
         
        // test for xml
        Path inputXml = jobArtefacts.resolve((job.getId()+".xml")).toAbsolutePath();
        if (job.getFxPath() != null && !job.getFxPath().isEmpty()) {
            inputXml = jobArtefacts.resolve(job.getFxPath()).toAbsolutePath();
        }

        if (Files.exists(inputXml)) {
            fileToValidate = inputXml;
        } 
        
        if (fileToValidate == null) {
            status.setStatus(Status.STATUS_ERROR);
            return status.addError("FILE_NOT_FOUND", "No PDF or XML found");
        }

        // --- Mustang Validator instanziieren ---
        ZUGFeRDValidator validator = new ZUGFeRDValidator();

        // Den Validator direkt mit dem PDF füttern.
        // Er extrahiert das XML automatisch und checkt PDF/A-Konformität + XML-Schema + Schematron.
        String validationReportXml = validator.validate(fileToValidate.toString());

        // --- Den zurückgegebenen XML-Report parsen ---
        DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
        DocumentBuilder builder = factory.newDocumentBuilder();
        Document doc = builder.parse(new ByteArrayInputStream(validationReportXml.getBytes(StandardCharsets.UTF_8)));

        ValidationLogVisualizer vlvi = new ValidationLogVisualizer();
        String resultFileName = getBaseName(fileToValidate.toString()) + ".result.pdf";
        vlvi.toPDF(validationReportXml, resultFileName);

        return status;
    }
    
    // Extract takes a pdfFile and extracts the factur-x.xml and attachments
    // input: job.pdfPath
    // output: artefactsRoot/jobID/jobID.xml []artefactsRoot/jobID/attachment
    private Status extract(JobConfig job) throws Exception {
        // Default Werte setzen
        Status status = Status.ok(job.getQueue().name());
        Path inputPdf = jobArtefacts.resolve((job.getId()+".pdf")).toAbsolutePath();

        // Update mit job-werten
        if (job.getPdfPath() != null && !job.getPdfPath().isEmpty()) {
            inputPdf = jobArtefacts.resolve(job.getPdfPath()).toAbsolutePath();
        }
        String pdfFile = inputPdf.toString();
        Logger.error("Input PDF {} not found", pdfFile);
        if (!Files.exists(inputPdf)) {
            return status.addError("PDF_NOT_FOUND", "Input PDF not found");
        }

        // ---------------------------------------------------------
        // 1. Das Rechnungs-XML extrahieren (via Mustang API)
        // ---------------------------------------------------------
        ZUGFeRDImporter importer = new ZUGFeRDImporter(pdfFile);
        byte[] xmlData = importer.getRawXML();
        
        if (xmlData != null && xmlData.length > 0) {
            Files.write(Paths.get(jobArtefacts.toString(), "factur-x.xml"), xmlData);
            status.getMeta().put("output_file", "ZUGFeRD XML successful extracted.");
        }

        // ---------------------------------------------------------
        // 2. Weitere Attachments extrahieren (via Apache PDFBox)
        // ---------------------------------------------------------
        try (PDDocument document = Loader.loadPDF(new File(pdfFile))) {
            
            PDDocumentNameDictionary namesDictionary = new PDDocumentNameDictionary(document.getDocumentCatalog());
            PDEmbeddedFilesNameTreeNode efTree = namesDictionary.getEmbeddedFiles();
            
            if (efTree != null) {
                Map<String, PDComplexFileSpecification> names = efTree.getNames();
                if (names != null) {
                    
                    // initialize attachmentsList
                    if (job.getAttachments() == null) {
                        job.setAttachments(new ArrayList<>());
                    }
                    
                    // now extract attachments
                    for (Map.Entry<String, PDComplexFileSpecification> entry : names.entrySet()) {
                        String filename = entry.getKey();
                        
                        // Das ZUGFeRD XML überspringen, das haben wir oben schon komfortabel via Mustang geholt
                        // aber vorher noch in job-file schreiben
                        if ("factur-x.xml".equals(filename) || "zugferd-invoice.xml".equals(filename)) {
                            job.setFxPath(filename);
                            continue;
                        }

                        // Anhang auslesen und speichern
                        PDComplexFileSpecification fileSpec = entry.getValue();
                        PDEmbeddedFile embeddedFile = fileSpec.getEmbeddedFile();
                        
                        if (embeddedFile != null) {
                            byte[] fileBytes = embeddedFile.toByteArray();

                            String attachName = (Paths.get(filename)).getFileName().toString();
                            Path attachPath = Paths.get(jobArtefacts.toString(), attachName);
                            Files.write(attachPath, fileBytes, StandardOpenOption.CREATE_NEW);
                            status.getMeta().put(filename, "Additional attachments extracted");

                            job.getAttachments().add(attachPath.toString());
                        }
                    }
                }
            }
        }
        return status;
    }
    
    private String getBaseName(String inputFilename) {
        int dot = inputFilename.lastIndexOf('.');
        if (dot <= 0) {
            return inputFilename;
        } else {
            return inputFilename.substring(0, dot);        
        }
    }
}

