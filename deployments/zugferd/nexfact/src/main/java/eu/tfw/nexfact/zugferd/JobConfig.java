// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package eu.tfw.nexfact.zugferd;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class JobConfig {

    // not all attributes are used, but must be there for synchronize calling process
    private String id;
    private String status;
    private JobQueue queue;
    private String message;
    private String timestamp;
    
    private List<String> attachments;
    private String invoiceID;
    private String pdfPath;
    private String fxPath;
    private String parPath;
    private String outputPath;
       
    public JobConfig() {
    }

    @JsonProperty("id")
    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

     @JsonProperty("status")
    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    @JsonProperty("message")
    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    @JsonProperty("timestamp")
    public String getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(String timestamp) {
        this.timestamp = timestamp;
    }

    @JsonProperty("queue")
    public JobQueue getQueue() {
        return queue;
    }

    public void setQueue(JobQueue queue) {
        this.queue = queue;
    }

    @JsonProperty("attachments")
    public List<String> getAttachments() {
        return attachments;
    }

    public void setAttachments(List<String> attachments) {
        this.attachments = attachments;
    }

    @JsonProperty("invoice_id")
    public String getInvoiceID() {
        return invoiceID;
    }

    public void setInvoiceID(String id) {
        this.invoiceID = id;
    }
    
    @JsonProperty("pdf_path")
    public String getPdfPath() {
        return pdfPath;
    }

    public void setPdfPath(String pdf) {
        this.pdfPath = pdf;
    }

    @JsonProperty("fx_path")
    public String getFxPath() {
        return fxPath;
    }

    public void setFxPath(String xml) {
        this.fxPath = xml;
    }
    
    @JsonProperty("par_path")
    public String getParPath() {
        return parPath;
    }

    public void setParPath(String parPath) {
        this.parPath = parPath;
    }

    @JsonProperty("output_path")
    public String getOutputPath() {
        return outputPath;
    }

    public void setOutputPath(String output) {
        this.outputPath = output;
    }

}
