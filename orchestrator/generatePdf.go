// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package orchestrator

import (
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
)

func (o *Orchestrator) generatePdf(rJob *job.RenderJob, finalize bool) (*writer.RenderOutput, error) {

	engine := config.SystemParams.Engine.PdfEngine     // loEngine.odt
	provider := config.SystemParams.Engine.PdfProvider // native | container | webservice
	kind := writer.PdfKind
	timeout := config.SystemParams.Engine.CmdTimeout
	heartbeat := config.SystemParams.Engine.HeartbeatFrequency
	pdfWriter, err := writer.CreateWriter(kind, provider, engine, timeout, heartbeat)
	if err != nil {
		return nil, err
	}

	renderInput := writer.RenderInput{}
	renderInput.OttPath = o.ZugferdMaster.Paths.OTTfile
	renderInput.ArtefactsPath = filepath.Join(config.SellerParams.SellerRoot, "data", "artefacts")
	renderInput.RootPath = filepath.Join(config.SellerParams.SellerRoot, "data", "pdfs")
	renderInput.OutputPath = filepath.Join(config.SellerParams.SellerRoot, "outputs", o.ZugferdMaster.Invoice.InvoiceID+".pdf")

	return pdfWriter.Render(o.ctx, rJob, &renderInput, finalize)
}
