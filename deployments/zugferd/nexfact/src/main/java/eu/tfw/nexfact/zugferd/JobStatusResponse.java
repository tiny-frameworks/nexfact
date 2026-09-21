package eu.tfw.nexfact.zugferd;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL) // Entspricht dem Go 'omitempty'
public class JobStatusResponse {
    private String status;
    private String message;
    private String timestamp;

    // --- Getter und Setter ---
    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }

    public String getMessage() { return message; }
    public void setMessage(String message) { this.message = message; }

    public String getTimestamp() { return timestamp; }
    public void setTimestamp(String timestamp) { this.timestamp = timestamp; }
}
