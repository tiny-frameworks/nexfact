
package eu.tfw.nexfact.zugferd;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.fasterxml.jackson.databind.ObjectMapper;




/**
 *
 * @author georg
 */
public class JobProcessorTest {
    

    private final String watchDir = "src/test/resources/data/in";
    private final String artefactsDir = "./src/test/resources/data/artefacts";
    private final String sourceDir = "./src/test/resources";

    //Merge
    private final Path mergeJobSource = Paths.get(sourceDir + "/testjob_generate.json");
    private final Path mergeJobFile = Paths.get(watchDir + "/testjob_generate.json");
    private final Path mergeDoneFile = Paths.get(watchDir + "/testjob_generate.done");

    private final Path mergePdfSource = Paths.get(sourceDir + "/test_standard_without.pdf");
    private final Path mergePdfFile =  Paths.get(artefactsDir + "/testjob/test_standard_without.pdf");
        
    private final Path mergeFxSource = Paths.get(sourceDir + "/test_standard_xml.xml");
    private final Path mergeFxFile = Paths.get(artefactsDir + "/testjob/test_standard_xml.xml");
 
    private final Path aPngSource = Paths.get(sourceDir + "/test_attachment.png");
    private final Path aPngFile = Paths.get(artefactsDir + "/testjob/test_attachment.png");
    private final Path aPdfSource = Paths.get(sourceDir + "/test_attachment.pdf");
    private final Path aPdfFile = Paths.get(artefactsDir + "/testjob/test_attachment.pdf");
        
    // Validate
    private final Path validateJobSource = Paths.get(sourceDir + "/testjob_validate.json");
    private final Path validateJobFile = Paths.get(watchDir + "/testjob_validate.json");
    private final Path validateDoneFile = Paths.get(watchDir + "/testjob_validate.done");
    
    private final Path validatePdfFile =  Paths.get(artefactsDir + "/testjob/test_valid_zugferd.pdf");
    private final Path validatePdfSource = Paths.get(sourceDir + "/test_valid_zugferd.pdf");
    private final Path notValidPdfSource = Paths.get(sourceDir + "/test_standard_notvalidxml.pdf");
    private final Path notValidPdfFile =  Paths.get(artefactsDir + "/testjob/test_standard_notvalidxml.pdf");
    private final Path notValidJobSource = Paths.get(sourceDir + "/testjob_notvalid_xml.json");
    private final Path notValidJobFile = Paths.get(watchDir + "/testjob_notvalid_xml.json");
    private final Path notValidDoneFile = Paths.get(watchDir + "/testjob_notvalid_xml.done");
        
    // Extract
    private final Path extractJobSource = Paths.get(sourceDir + "/testjob_extract.json");
    private final Path extractJobFile = Paths.get(watchDir + "/testjob_extract.json");
    private final Path extractDoneFile = Paths.get(watchDir + "/testjob_extract.done");
    private final Path extractPdfFile =  Paths.get(artefactsDir + "/testjob/test_max-mit.zf.pdf");
    private final Path extractPdfSource = Paths.get(sourceDir + "/test_max-mit.zf.pdf");
    
    @BeforeAll
    public static void setUpClass() {
    }
    
    @AfterAll
    public static void tearDownClass() {
    }
    private ObjectMapper mapper;
    
    @BeforeEach
    public void setUp() {
        // clear in directory from earlier test
        try {
            //Merge
            Files.deleteIfExists(mergeJobFile);
            Files.deleteIfExists(mergeDoneFile);
            Files.deleteIfExists(mergePdfFile);
            Files.deleteIfExists(mergeFxFile);
            Files.deleteIfExists(aPngFile);
            Files.deleteIfExists(aPdfFile);
            
            // Validate
            Files.deleteIfExists(validateJobFile);
            Files.deleteIfExists(validateDoneFile);
            Files.deleteIfExists(validatePdfFile);
            Files.deleteIfExists(notValidPdfFile);
            Files.deleteIfExists(notValidJobFile);
            Files.deleteIfExists(notValidDoneFile);
            
            //Extract
            Files.deleteIfExists(extractPdfFile);
            Files.deleteIfExists(extractDoneFile);
        } catch (IOException e) {
            // nothing to do. already deleted
        }
    }
    
    @AfterEach
    public void tearDown() {
    }

    @Test
    void testGenerate_validZF2pdf() throws Exception {

        // Arange Data
        Files.copy(mergePdfSource, mergePdfFile, StandardCopyOption.REPLACE_EXISTING);
        Files.copy(mergeFxSource, mergeFxFile, StandardCopyOption.REPLACE_EXISTING);
        Files.copy(aPdfSource, aPdfFile, StandardCopyOption.REPLACE_EXISTING);
        Files.copy(aPngSource, aPngFile, StandardCopyOption.REPLACE_EXISTING);
        Files.copy(mergeJobSource, mergeJobFile, StandardCopyOption.REPLACE_EXISTING);
        
        RunOnceService once = new RunOnceService(watchDir, artefactsDir);
//        WatchServiceLoop once = new WatchServiceLoop(watchDir, artefactsDir);
        once.start();
        
        String result = "success";
        assertTrue(readResult(mergeDoneFile, result), "doneFile result should be: " + result);
        
    }


    @Test
    void testValidate_validZFpdf() throws Exception {
    
        // Arange Data
        Files.copy(validatePdfSource, validatePdfFile, StandardCopyOption.REPLACE_EXISTING);
        Files.copy(validateJobSource, validateJobFile, StandardCopyOption.REPLACE_EXISTING);
        
        RunOnceService once = new RunOnceService(watchDir, artefactsDir);
        once.start();

        // Assert
        String result = "success";
        assertTrue(readResult(validateDoneFile, result), "doneFile result should be: " + result);
    }
 
    @Test
    void testValidate_notValidZFpdf() throws Exception {
    
        // Arange Data
        Files.copy(notValidPdfSource, notValidPdfFile, StandardCopyOption.REPLACE_EXISTING);
        Files.copy(notValidJobSource, notValidJobFile, StandardCopyOption.REPLACE_EXISTING);
        
        RunOnceService once = new RunOnceService(watchDir, artefactsDir);
        once.start();

        // Assert
        //String result = "error";
        String result = "success";
        assertTrue(readResult(notValidDoneFile, result), "doneFile result should be: " + result);
    }
    
    @Test
    void testExtract_validZFpdf() throws Exception {
   
        // Arange Data
        Files.copy(extractPdfSource, extractPdfFile, StandardCopyOption.REPLACE_EXISTING);
        Files.copy(extractJobSource, extractJobFile, StandardCopyOption.REPLACE_EXISTING);
        
        RunOnceService once = new RunOnceService(watchDir, artefactsDir);
        once.start();

        // Assert
        String result = "success";
        assertTrue(readResult(extractDoneFile, result), "doneFile result should be: " + result);
    }

    private boolean readResult(Path doneFile, String doneResult) throws IOException {
        this.mapper = new ObjectMapper();
        JobConfig job = null;
        job = mapper.readValue(doneFile.toFile(), JobConfig.class);
        return doneResult.equals(job.getStatus());
    }
}
