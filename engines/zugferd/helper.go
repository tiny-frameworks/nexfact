// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package zugferd

import (
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
)

func (w *ZugferdWriter) buildJobArtefacts(rJob *job.RenderJob, rInput *writer.RenderInput) error {

	// input-pdf -> jobID.pdf
	source := rInput.PdfPath
	if source != "" {
		pdfPath := filepath.Join(w.JobArtefactsPath, rJob.ID+".pdf")
		err := writer.CopyFile(source, pdfPath)
		if err != nil {
			return err
		}
	}

	// input-xml -> jobID.xml
	source = rInput.XMLPath
	if source != "" {
		fxPath := filepath.Join(w.JobArtefactsPath, rJob.ID+".xml")
		err := writer.CopyFile(source, fxPath)
		if err != nil {
			return err
		}
	}

	// input-attachments -> jobArtefactsDir
	for _, attach := range rInput.Attachments {
		attachBase := filepath.Base(attach)
		attachJob := filepath.Join(w.JobArtefactsPath, attachBase)
		writer.CopyFile(attach, attachJob)

		//attachContainer := filepath.Join("/", "data", "artefacts", rJob.ID, attachBase)
		rJob.Attachments = append(rJob.Attachments, attachBase)
	}

	if err := w.WriteJob(rJob); err != nil {
		return err
	}
	return nil
}
