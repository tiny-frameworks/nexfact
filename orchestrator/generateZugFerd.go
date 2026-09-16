// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package orchestrator

import (
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
)

func (o *Orchestrator) generateZugferd(rJob *job.RenderJob, finalize bool) (*writer.RenderOutput, error) {

	engine := config.SystemParams.Engine.ZugferdEngine     // mustang-cli.jar
	provider := config.SystemParams.Engine.ZugferdProvider // native | container | webservice
	kind := writer.ZugferdKind
	timeout := config.SystemParams.Engine.CmdTimeout
	heartbeat := config.SystemParams.Engine.HeartbeatFrequency
	zugferdWriter, err := writer.CreateWriter(kind, provider, engine, timeout, heartbeat)
	if err != nil {
		return nil, err
	}

	renderInput := writer.RenderInput{}
	renderInput.ArtefactsPath = filepath.Join(config.SellerParams.SellerRoot, "data", "artefacts")
	renderInput.RootPath = filepath.Join(config.SellerParams.SellerRoot, "data", "zugferds")
	renderInput.OutputPath = filepath.Join(config.SellerParams.SellerRoot, "outputs", o.ZugferdMaster.Invoice.InvoiceID+".zf.pdf")

	return zugferdWriter.Render(o.ctx, rJob, &renderInput, finalize)
}

func (o *Orchestrator) combine(rJob *job.RenderJob, finalize bool) (*writer.RenderOutput, error) {

	engine := config.SystemParams.Engine.ZugferdEngine     // mustang-cli.jar
	provider := config.SystemParams.Engine.ZugferdProvider // native | container | webservice
	kind := writer.ZugferdKind
	timeout := config.SystemParams.Engine.CmdTimeout
	heartbeat := config.SystemParams.Engine.HeartbeatFrequency
	zugferdWriter, err := writer.CreateWriter(kind, provider, engine, timeout, heartbeat)
	if err != nil {
		return nil, err
	}

	renderInput := writer.RenderInput{}
	renderInput.ArtefactsPath = filepath.Join(config.SellerParams.SellerRoot, "data", "artefacts")
	renderInput.RootPath = filepath.Join(config.SellerParams.SellerRoot, "data", "zugferds")

	renderInput.PdfPath = o.ZugferdMaster.Options.PDFSource
	renderInput.XMLPath = o.ZugferdMaster.Options.FacturXSource
	renderInput.OutputPath = filepath.Join(config.SellerParams.SellerRoot, "outputs", o.ZugferdMaster.Invoice.InvoiceID+".zf.pdf")

	for _, attach := range o.ZugferdMaster.Options.Attachments {
		if attPath, err := filepath.Abs(attach); err == nil {
			renderInput.Attachments = append(renderInput.Attachments, attPath)
		}
	}

	return zugferdWriter.Render(o.ctx, rJob, &renderInput, finalize)
}
