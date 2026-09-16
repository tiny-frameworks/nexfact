// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package pdf

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

type PdfWriter struct {
	writer.AbstractWriter
}

func New(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (*PdfWriter, error) {

	w := &PdfWriter{}
	w.Engine = engine
	w.TimeOut = timeout
	w.Heartbeat = heartbeat
	w.Provider = provider

	return w, nil
}

// automatic registration in api.registry at startup
func init() {
	writer.Register(writer.PdfKind, func(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (writer.Writer, error) {
		return New(provider, engine, timeout, heartbeat)
	})
}

func (w *PdfWriter) Render(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput, finalize bool) (*writer.RenderOutput, error) {

	if err := w.Prepare(rJob, rInput); err != nil {
		return nil, err
	}

	// build additional paths for pdf generating in artefactsdir
	ottPath := filepath.Join(w.JobArtefactsPath, rJob.ID+".ott")
	rInput.ParPath = filepath.Join(w.JobArtefactsPath, rJob.ID+".par")

	if err := writer.CopyFile(rInput.OttPath, ottPath); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("can not move file %s", rInput.OttPath),
			"engines.pdf.render",
			err,
		)
	}

	switch w.Provider {
	case writer.NativeProvider:
		if err := w.renderNative(ctx, rJob, rInput); err != nil {
			return nil, err
		}
	case writer.ContainerProvider:
		if err := w.renderContainer(ctx, rJob, rInput); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(
			errors.InvalidValue,
			fmt.Sprintf("Unknown Provider: %s", w.Provider),
			"engines.pdf.render",
		)
	}

	sourcePath := filepath.Join(w.JobArtefactsPath, rJob.ID+".pdf")
	return w.Finalize(rJob, rInput, sourcePath, finalize)
}
