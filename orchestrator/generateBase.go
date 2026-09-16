// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package orchestrator

import (
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexfact/api/job"

	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/generator"
)

func (o *Orchestrator) generatePar(rJob *job.RenderJob, finalize bool) (*writer.RenderOutput, error) {

	var parBytes []byte
	var parwriter writer.Writer

	engine := ".par"
	provider := writer.BaseProvider // par is always native
	kind := writer.BaseKind
	timeout := config.SystemParams.Engine.CmdTimeout
	heartbeat := config.SystemParams.Engine.HeartbeatFrequency
	parwriter, err := writer.CreateWriter(kind, provider, engine, timeout, heartbeat)
	if err != nil {
		return nil, err
	}
	parTemplate := config.SellerParams.Templates.ParTemplate
	parGenerator, err := generator.NewPar(parTemplate)
	if err != nil {
		return nil, err
	}

	if parBytes, err = parGenerator.Generate(o.ZugferdMaster); err != nil {
		return nil, err
	}

	renderInput := writer.RenderInput{}
	renderInput.BytesToRender = parBytes
	renderInput.ArtefactsPath = filepath.Join(config.SellerParams.SellerRoot, "data", "artefacts")
	renderInput.RootPath = filepath.Join(config.SellerParams.SellerRoot, "data", "base")
	renderInput.OutputPath = filepath.Join(config.SellerParams.SellerRoot, "outputs", o.ZugferdMaster.Invoice.InvoiceID+".par")

	return parwriter.Render(o.ctx, rJob, &renderInput, finalize)
}

func (o *Orchestrator) generateFacturX(rJob *job.RenderJob, finalize bool) (*writer.RenderOutput, error) {

	var xmlBytes []byte
	var fxwriter writer.Writer

	engine := ".xml"
	provider := writer.BaseProvider // facturX is always native
	kind := writer.BaseKind
	timeout := config.SystemParams.Engine.CmdTimeout
	heartbeat := config.SystemParams.Engine.HeartbeatFrequency
	fxwriter, err := writer.CreateWriter(kind, provider, engine, timeout, heartbeat)
	fxTemplate := config.SellerParams.Templates.XmlTemplate

	fxGenerator, err := generator.NewFX(fxTemplate)
	if err != nil {
		return nil, err
	}

	if xmlBytes, err = fxGenerator.Generate(o.ZugferdMaster); err != nil {
		return nil, err
	}

	renderInput := writer.RenderInput{}
	renderInput.BytesToRender = xmlBytes
	renderInput.ArtefactsPath = filepath.Join(config.SellerParams.SellerRoot, "data", "artefacts")
	renderInput.RootPath = filepath.Join(config.SellerParams.SellerRoot, "data", "base")
	renderInput.OutputPath = filepath.Join(config.SellerParams.SellerRoot, "outputs", o.ZugferdMaster.Invoice.InvoiceID+".xml")

	return fxwriter.Render(o.ctx, rJob, &renderInput, finalize)
}
