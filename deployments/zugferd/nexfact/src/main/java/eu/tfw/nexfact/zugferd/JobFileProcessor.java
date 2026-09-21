// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package eu.tfw.nexfact.zugferd;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import org.tinylog.Logger;

/**
 *
 * @author georg
 */
public abstract class JobFileProcessor {

        protected final String inDir;
        private final ObjectMapper mapper;
        protected final Path artefactsRoot;
        
    public JobFileProcessor(String inDir, String artefacts) {
        this.inDir = inDir;
        this.artefactsRoot = Paths.get(artefacts);
        this.mapper = new ObjectMapper();
        this.mapper.enable(SerializationFeature.INDENT_OUTPUT);
    }
    
    protected void processJobFileSafe(File jobFile) {
        // 1. ATOMIC LOCK: Wir benennen die Datei um. 
        // Das klappt unter Linux nur, wenn kein anderer Prozess (Go) mehr exklusiv darauf schreibt.
        File processingFile = new File(jobFile.getAbsolutePath() + ".processing");
        if (!jobFile.renameTo(processingFile)) {
            // Datei ist noch gelockt oder wird gerade geschrieben. 
            // Wir überspringen sie in diesem Durchlauf einfach und versuchen es in 500ms nochmal.
            return; 
        }

        JobConfig job = null;
        Status status;
 
        try {
            // Wir lesen nun von der umbenannten .processing Datei
            job = mapper.readValue(processingFile, JobConfig.class);
            job.setStatus(Status.STATUS_RUNNING);
            
            Logger.debug("nach marshalling: {}", job.getQueue());
            
            // build jobDirectory
            Path jobDir = artefactsRoot.resolve(job.getId()).toAbsolutePath();
            try {
                Files.createDirectories(jobDir);
            } catch (IOException ex) {
                
            }
            
            JobProcessor processor = new JobProcessor();
            status = processor.process(job, jobDir);
            job.setStatus(status.getStatus());
        } catch (Exception ex) {
            String mode = (job != null && job.getQueue() != null)
                    ? job.getQueue().name()
                    : "UNKNOWN";

            status = Status.error(mode);
            status.getErrors().add(
                    new StatusMessage("PROCESSING_ERROR", ex.getMessage())
            );
            job.setStatus(Status.STATUS_ERROR);
        }

        // 2. STATUS SCHREIBEN & AUFRÄUMEN
        try {
            // Wir schreiben das JobFile, StatusFile und *_result.pdf
            mapper.writeValue(processingFile, job); 

            String resultName = jobFile.toString();            
            
            StatusFileWriter sfw = new StatusFileWriter();
            File statusFile = new File(resultName.replace(".json", ".status"));
            sfw.writeStatusFile(status, statusFile);

            File doneFile = new File(resultName.replace(".json", ".done"));
            processingFile.renameTo(doneFile);
            
        } catch (IOException e) {
            System.err.println("Failed to write status file for " + jobFile.getName() + ": " + e.getMessage());
            Logger.error("Failed to write status file for {}: {}", jobFile.getName(), e.getMessage());
            
        } finally {
            // Löschen der .processing Datei, egal ob Erfolg oder Fehler
            //processingFile.delete(); 
        }
    }
}
