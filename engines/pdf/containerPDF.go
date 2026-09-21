// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package pdf

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

// WATCH MODE (Container must already running)
func (w *PdfWriter) renderContainer(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput) error {
	hbFile := filepath.Join(rInput.RootPath, "logs", "nexfact.heartbeat")
	if _, err := writer.CheckHeartbeat(hbFile, w.Heartbeat*2); err != nil {
		return err
	}

	if err := os.Rename(w.JobTmpFile, w.JobJsonFile); err != nil {
		return errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot rename jobfile %s in %s", w.JobTmpFile, w.JobJsonFile),
			"engines.pdf.renderWithContainerProvider",
			err,
		)
	}

	renderCtx, cancel := context.WithTimeout(ctx, w.TimeOut)
	defer cancel()

	ticker := time.NewTicker(writer.PollIntervall)
	defer ticker.Stop()

	for {
		select {
		case <-renderCtx.Done():
			return errors.Wrap(
				errors.ContainerTimeout,
				fmt.Sprintf("container timeout %v reached", w.TimeOut),
				"nexfact-engines.pdf.renderWithContainerProvider",
				renderCtx.Err(),
			)

		case <-ticker.C:
			// 1. Check if the status file has appeared.
			var nerr *errors.Error
			if nerr = rJob.Result(w.JobResultFile); nerr != nil {
				if nerr.Code == errors.ReadError {
					// Fallback: If the file system stalls for a fraction of a millisecond,
					// we simply try again on the next tick (timeout).
					continue
				} else {
					return nerr // A genuine error, whatever it is.
				}
			} else {
				return nil
			}
		}
	}
	return nil
}
