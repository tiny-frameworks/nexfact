// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package orchestrator

import (
	"path/filepath"
	"strings"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
)

func (o *Orchestrator) validate(rJob *job.RenderJob, finalize bool) (renderoutput *writer.RenderOutput, err error) {

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

	// outputpath is either assembled by InvoiceID or facturx-source or pdf-source
	resultBase := o.ZugferdMaster.Invoice.InvoiceID
	if resultBase == "" {
		tmpBase := o.ZugferdMaster.Options.FacturXSource
		if tmpBase == "" {
			tmpBase = o.ZugferdMaster.Options.PDFSource
		}
		resultBase = strings.TrimSuffix(filepath.Base(tmpBase), filepath.Ext(tmpBase))
	}
	renderInput.OutputPath = filepath.Join(config.SellerParams.SellerRoot, "outputs", resultBase+".result.pdf")

	return zugferdWriter.Render(o.ctx, rJob, &renderInput, finalize)

}
