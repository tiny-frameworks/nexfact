// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"codeberg.org/tiny-frameworks/nexutils/errors"
)

type StatusKind string

const (
	StatusOK      StatusKind = "success"
	StatusError   StatusKind = "error"
	StatusRunning StatusKind = "running"
)

type QueueKind string

const (
	QueueExtract  = "extract"
	QueueValidate = "validate"
	QueueCombine  = "combine" // ZUGFeRD
	QueuePar      = "generate_par"
	QueueFacturX  = "generate_fx"
	QueueZugferd  = "generate_zf"
	QueuePdf      = "generate_pdf"
)

type RenderJob struct {
	ID        string     `json:"id"`
	Status    StatusKind `json:"status"`
	Queue     QueueKind  `json:"queue"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`

	InvoiceID   string   `json:"invoice_id,omitempty"` // Falls schon bekannt
	Attachments []string `json:"attachments,omitempty"`

	PdfPath string `json:"pdf_path,omitempty"`
	FXPath  string `json:"fx_path,omitempty"`
	ParPath string `json:"par_path,omitempty"`
	// OutputPath string `json:"output_path"`
}

func New() *RenderJob {

	var id string
	// 1. Format timestamp (Go's magic Referencedate: 2006-01-02 15:04:05)
	timestamp := time.Now().Format("20060102-150405")

	// 2. Generate short random-part (2 Bytes = 4 Hex)
	// we use crypto/rand, since this is safer with concurrency (goroutines)
	// than math/rand.
	randomBytes := make([]byte, 2)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback: It is extremely unlikely that this will happen.)
		id = fmt.Sprintf("%s-%04d", timestamp, time.Now().UnixNano()%10000)
	} else {
		randomSuffix := hex.EncodeToString(randomBytes)
		id = fmt.Sprintf("%s-%s", timestamp, randomSuffix)
	}

	// 3. Assembly
	return &RenderJob{
		ID:        id,
		Timestamp: timestamp,
		Status:    StatusRunning,
	}
}

func (rJob *RenderJob) Result(doneFile string) *errors.Error {

	// Read status JSON
	if _, err := os.Stat(doneFile); err != nil {
		return errors.New(
			errors.ReadError,
			fmt.Sprintf("cannot find file: %v", doneFile),
			"api.job.result",
		)
	}

	jobBytes, err := os.ReadFile(doneFile)
	if err != nil {
		return errors.Wrap(
			errors.ReadError,
			fmt.Sprintf("cannot read file: %v", doneFile),
			"api.job.result",
			err,
		)
	}

	//var jc RenderJob
	if err := json.Unmarshal(jobBytes, &rJob); err != nil {
		return errors.Wrap(
			errors.RuleViolation,
			fmt.Sprintf("unable to unmarshal file: %v", doneFile),
			"api.job.result",
			err,
		)
	}

	if rJob.Status == StatusOK {
		return nil
	} else {
		return errors.New(
			errors.Inconsistent,
			rJob.Message,
			"core.job.jobresult.error",
		)
	}
}
