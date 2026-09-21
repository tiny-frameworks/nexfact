// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package eu.tfw.nexfact.zugferd;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Status {

    public static String STATUS_ERROR = "error";
    public static String STATUS_SUCCESS = "success";
    public static String STATUS_RUNNING = "running";
      
    private String status;          // success | error | running
    private String mode;
    private String output;

    private List<StatusMessage> warnings = new ArrayList<>();
    private List<StatusMessage> errors = new ArrayList<>();

    private Map<String, Object> meta = new HashMap<>();

    public Status() {}

    public static Status ok(String mode) {
        Status s = new Status();
        s.status = STATUS_SUCCESS;
        s.mode = mode;
        return s;
    }

    public static Status error(String mode) {
        Status s = new Status();
        s.status = STATUS_ERROR;
        s.mode = mode;
        return s;
    }

    // Getter / Setter

    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }

    public String getMode() { return mode; }
    public void setMode(String mode) { this.mode = mode; }

    public String getOutput() { return output; }
    public void setOutput(String output) { this.output = output; }

    public List<StatusMessage> getWarnings() { return warnings; }
    public void setWarnings(List<StatusMessage> warnings) { this.warnings = warnings; }

    public List<StatusMessage> getErrors() { return errors; }
    public void setErrors(List<StatusMessage> errors) { this.errors = errors; }

    public Map<String, Object> getMeta() { return meta; }
    public void setMeta(Map<String, Object> meta) { this.meta = meta; }
 
    public Status addError(String code, String message) {
        this.status = Status.STATUS_ERROR;
        this.errors.add(new StatusMessage(code, message));
        return this;
    }
}
