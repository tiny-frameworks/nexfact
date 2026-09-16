// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package pdf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

func (w *PdfWriter) renderNative(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput) error {

	if err := os.Rename(w.JobTmpFile, w.JobJsonFile); err != nil {
		return errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot rename jobfile %s in %s", w.JobTmpFile, w.JobJsonFile),
			"engines.pdf.renderWithNativeProvider",
			err,
		)
	}
	renderCtx, cancel := context.WithTimeout(ctx, w.TimeOut)
	defer cancel()

	w.Executable = writer.PdfExecutable
	macroCall := fmt.Sprintf("macro://./Standard.FileEngine.GenerateDoc(\"%s\")", rInput.ParPath)
	cmd := exec.CommandContext(renderCtx, w.Executable,
		"--headless",
		"--nologo",
		"--nodefault",
		"--norestore",
		"--nolockcheck",
		w.Engine,
		macroCall,
	)

	_, stderr := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return errors.Wrap(
			errors.AmbiguousValue,
			fmt.Sprintf("LibreOffice was frozen and was killed by Go!"),
			"engines.pdf.renderWithNativeProvider",
			ctx.Err(),
		)
	}

	if stderr != nil {
		// check if LO has written a .pdf and a .status file
		statusResult := strings.TrimSuffix(rInput.ParPath, ".par") + ".status"
		fmt.Println(statusResult)

	}

	os.Rename(w.JobJsonFile, w.JobResultFile)
	return nil
}
