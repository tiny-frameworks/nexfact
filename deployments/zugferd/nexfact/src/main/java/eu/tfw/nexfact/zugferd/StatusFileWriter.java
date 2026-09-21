package eu.tfw.nexfact.zugferd;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import java.io.File;
import java.io.IOException;

public class StatusFileWriter {

    private final ObjectMapper mapper;

    public StatusFileWriter() {
        this.mapper = new ObjectMapper();

        // Macht das JSON schön lesbar (optional, aber hilfreich fürs Debugging)
        this.mapper.enable(SerializationFeature.INDENT_OUTPUT); 
    }

    public void writeStatusFile(Status internalStatus, File outputFile) {
        try {
            mapper.writeValue(outputFile, internalStatus);
            System.out.println("Status file written to: " + outputFile.getAbsolutePath());
        } catch (IOException e) {
            // Falls das Schreiben des Status-Files fehlschlägt (z.B. keine Rechte)
            System.err.println("CRITICAL: Failed to write status file! " + e.getMessage());
        }
    }
}