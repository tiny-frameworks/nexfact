// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package zugferd

import (
	"context"
	"fmt"

	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

type ZugferdWriter struct {
	writer.AbstractWriter
	Attachments []string
	SourcePDF   string
	SourceXMl   string
	TargetPDF   string
}

func New(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (*ZugferdWriter, error) {

	if timeout < time.Second {
		return nil, errors.New(
			errors.InvalidValue,
			fmt.Sprintf("timeout too small: %v", timeout),
			"engines.zugferd.New(...)",
		)
	}

	w := &ZugferdWriter{}
	w.Heartbeat = heartbeat
	w.TimeOut = timeout
	w.Engine = engine
	w.Provider = provider

	return w, nil
}

// automatic registration in api.registry at startup
func init() {
	writer.Register(writer.ZugferdKind, func(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (writer.Writer, error) {
		return New(provider, engine, timeout, heartbeat)
	})
}

func (w *ZugferdWriter) Render(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput, finalize bool) (rOutput *writer.RenderOutput, err error) {

	if err := w.Prepare(rJob, rInput); err != nil {
		return nil, err
	}

	// build additional paths in artefactsJobDir
	if err := w.buildJobArtefacts(rJob, rInput); err != nil {
		return nil, err
	}

	switch w.Provider {
	case writer.NativeProvider: // w.Engine = mustang-cli.jar
		if rOutput, err = w.renderNative(ctx, rJob, rInput); err != nil {
			return nil, err
		}
	case writer.ContainerProvider: // w.Engine := container artefactsRoot
		if rOutput, err = w.renderContainer(ctx, rJob, rInput); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(
			errors.InvalidValue,
			fmt.Sprintf("Unknown Provider: %s", w.Provider),
			"engines.zugferd.render",
		)
	}

	return w.Finalize(rJob, rInput, rOutput.OutputPath, finalize)
}
