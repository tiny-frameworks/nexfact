// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

// commonWriter is responsible for *.par and factur-x.xml
type BaseWriter struct {
	writer.AbstractWriter
}

func New(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (*BaseWriter, error) {

	w := &BaseWriter{}
	w.Engine = engine
	w.TimeOut = timeout
	w.Heartbeat = heartbeat
	w.Provider = provider

	return w, nil
}

// automatic registration in api.registry at startup
func init() {
	writer.Register(writer.BaseKind, func(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (writer.Writer, error) {
		return New(provider, engine, timeout, heartbeat)
	})
}

func (w *BaseWriter) Render(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput, finalize bool) (*writer.RenderOutput, error) {
	if len(rInput.BytesToRender) == 0 {
		return nil, errors.New(
			errors.EmptyInput,
			"inputBytes to render size = 0",
			"engine.base.render",
		)
	}

	if err := w.Prepare(rJob, rInput); err != nil {
		return nil, err
	}

	// build filePaths for par or factur-x
	//fileExt := filepath.Ext(job.OutputPath)
	sourcePath := filepath.Join(w.JobArtefactsPath, rJob.ID+w.Engine)

	if err := os.Rename(w.JobTmpFile, w.JobJsonFile); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot rename job json file: %s in %s", w.JobTmpFile, w.JobJsonFile),
			"writer.mustang.conmerge.render",
			err,
		)
	}

	if err := os.WriteFile(sourcePath, rInput.BytesToRender, writer.FileMask); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("can no write file %s", sourcePath),
			"writer.base.writer.render",
			err,
		)
	}

	return w.Finalize(rJob, rInput, sourcePath, finalize)
}
