// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/input"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/invoice"
	"codeberg.org/tiny-frameworks/nexutils/errors"

	// These lines activate the engines in the background:
	_ "codeberg.org/tiny-frameworks/nexfact/engines/base"
	_ "codeberg.org/tiny-frameworks/nexfact/engines/pdf"
	_ "codeberg.org/tiny-frameworks/nexfact/engines/zugferd"
)

const EnvRootKey = "NEXFACT_ENVROOT"

type Orchestrator struct {
	ctx           context.Context
	ZugferdMaster *invoice.ZUGFeRDmaster
}

func New(envRoot string) (*Orchestrator, error) {

	envPath, err := filepath.Abs(envRoot)
	if err != nil {
		return nil, errors.Wrap(
			errors.InternalError,
			"can not find default working directory",
			"orchestrator.NewDefault()",
			err,
		)
	}

	sysErr := config.LoadSystem(envPath)
	if sysErr != nil {
		return nil, sysErr
	}

	return &Orchestrator{
		ctx: context.Background(),
	}, nil

}

func NewDefault() (*Orchestrator, error) {

	if envRoot := os.Getenv(EnvRootKey); envRoot != "" {
		return New(envRoot)
	}
	if wd, err := os.Getwd(); err != nil {
		return nil, errors.Wrap(
			errors.InternalError,
			"can not find default working directory",
			"orchestrator.NewDefault",
			err,
		)
	} else {
		envRoot := filepath.Join(wd, "env")
		return New(envRoot)
	}
}

func (o *Orchestrator) RunJsonJob(json []byte) (renderOutput *writer.RenderOutput, err error) {

	const finalize bool = true
	var inputJob *input.InputJob

	if inputJob, err = input.Parse(json); err != nil {
		return nil, err
	}

	if err = input.EnrichSeller(inputJob); err != nil {
		return nil, err
	}
	// normalizes input.invoice to ZUGFeRDmaster
	zugferdMaster := invoice.New()
	if err := zugferdMaster.NormalizeAndValidate(inputJob); err != nil {
		return nil, err
	} else {
		o.ZugferdMaster = zugferdMaster
	}

	renderJob := job.New()
	renderJob.InvoiceID = o.ZugferdMaster.Invoice.InvoiceID

	queue := strings.ToLower(inputJob.Options.Queue)
	switch queue {
	case job.QueueFacturX:
		renderJob.Queue = job.QueueFacturX
		return o.generateFacturX(renderJob, finalize)

	case job.QueuePar:
		renderJob.Queue = job.QueuePar
		return o.generatePar(renderJob, finalize)

	case job.QueuePdf: // generate pure pdf rom jsondata
		renderJob.Queue = job.QueuePdf
		if rOutput, err := o.generatePar(renderJob, !finalize); err != nil {
			return rOutput, err
		}
		return o.generatePdf(renderJob, finalize)

	case job.QueueZugferd: // generate Zugferd Invoice from jsondata
		renderJob.Queue = job.QueueZugferd
		if rOutput, err := o.generatePar(renderJob, !finalize); err != nil {
			return rOutput, err
		}
		if rOutput, err := o.generateFacturX(renderJob, !finalize); err != nil {
			return rOutput, err
		}
		if rOutput, err := o.generatePdf(renderJob, !finalize); err != nil {
			return rOutput, err
		}
		return o.generateZugferd(renderJob, finalize)

	case job.QueueValidate:
		renderJob.Queue = job.QueueValidate
		return o.validate(renderJob, finalize)

	case job.QueueExtract:
		renderJob.Queue = job.QueueExtract
		return o.extract(renderJob, finalize)

	case job.QueueCombine:
		renderJob.Queue = job.QueueCombine
		return o.combine(renderJob, finalize)

	default:
		return nil, errors.New(
			errors.RuleViolation,
			fmt.Sprintf("Invoice jobtype %s is not valid", renderJob.Queue),
			"orchestrator.orchestrator.RunJsonJob",
		)
	}
}
