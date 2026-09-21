// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package zugferd

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

func (w *ZugferdWriter) renderContainer(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput) (rOutput *writer.RenderOutput, err error) {
	hbFile := filepath.Join(rInput.RootPath, "logs", "nexfact.heartbeat")
	if _, err := writer.CheckHeartbeat(hbFile, w.Heartbeat*2); err != nil {
		return nil, err
	}

	if err := os.Rename(w.JobTmpFile, w.JobJsonFile); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot rename jobfile %s in %s", w.JobTmpFile, w.JobJsonFile),
			"engines.zugfer.renderWithContainerProvider",
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
			return nil, errors.Wrap(
				errors.ContainerTimeout,
				fmt.Sprintf("container timeout %v reached", w.TimeOut),
				"engines.pdf.renderWithContainerProvider",
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
					return nil, nerr // A genuine error, whatever it is.
				}
			} else {
				outPath := ""
				switch rJob.Queue {
				case job.QueueValidate:
					outPath = filepath.Join(w.JobArtefactsPath, rJob.ID+".result.pdf")
				case job.QueueExtract:
					outPath = filepath.Join(w.JobArtefactsPath, "factur-x.xml")
				default:
					outPath = filepath.Join(w.JobArtefactsPath, rJob.ID+".zf.pdf")
				}

				rOutput = &writer.RenderOutput{
					JobID:      rJob.ID,
					OutputPath: outPath,
				}

				return rOutput, nil
			}

		}
	}
	return nil, errors.New(
		errors.ZugferdContainerError,
		fmt.Sprintf("unknow error, should never happen because endless-loop"),
		"engines.pdf.renderWithContainerProvider",
	)
}
