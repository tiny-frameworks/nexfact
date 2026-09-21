// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package zugferd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

func (w *ZugferdWriter) renderNative(
	ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput) (rOutput *writer.RenderOutput, err error) {

	w.Executable = writer.ZugferdExecutable

	switch rJob.Queue {
	case job.QueueValidate: // w.Engine := mustang-cli.jar
		rOutput, err = w.nativeValidate(ctx, rJob, rInput)

	case job.QueueCombine: // w.Engine := mustang-cli.jar
		rOutput, err = w.nativeCombine(ctx, rJob, rInput)

	case job.QueueZugferd: // w.Engine := mustang-cli.jar
		rOutput, err = w.nativeZugferd(ctx, rJob, rInput)

	case job.QueueExtract:
		rOutput, err = w.nativeExtract(ctx, rJob, rInput)

	default:
		err = errors.New(
			errors.InvalidValue,
			fmt.Sprintf("Unknown Queue/Action: %s", rJob.Queue),
			"engines.zugferd.renderWithNativeProvider",
		)
	}
	if err == nil {
		os.Remove(w.JobJsonFile)
		os.Remove(w.JobTmpFile)
	}
	return rOutput, err
}

func (w *ZugferdWriter) nativeZugferd(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput) (rOutput *writer.RenderOutput, err error) {

	rInput.PdfPath = filepath.Join(w.JobArtefactsPath, rJob.ID+".pdf")
	rInput.XMLPath = filepath.Join(w.JobArtefactsPath, rJob.ID+".xml")

	return w.nativeCombine(ctx, rJob, rInput)
}

func (w *ZugferdWriter) nativeCombine(ctx context.Context, rJob *job.RenderJob, rInput *writer.RenderInput) (rOutput *writer.RenderOutput, err error) {

	if err := os.Rename(w.JobTmpFile, w.JobJsonFile); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot rename jobfile %s in %s", w.JobTmpFile, w.JobJsonFile),
			"engines.zugferd.renderNative.nativeZugferd",
			err,
		)
	}

	// 1. Create an isolated temp directory
	tmpDir, err := os.MkdirTemp("", "nexfact-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir) // Cleanup after the run

	pdfTmp := filepath.Join(tmpDir, filepath.Base(rInput.PdfPath))
	writer.CopyFile(rInput.PdfPath, pdfTmp)

	xmlTmp := filepath.Join(tmpDir, filepath.Base(rInput.XMLPath))
	writer.CopyFile(rInput.XMLPath, xmlTmp)

	zfTmp := filepath.Join(tmpDir, filepath.Base(rJob.ID+".zf.pdf"))

	// Argumente als Slice definieren
	args := []string{
		"-jar", w.Engine,
		"--action", "combine",
		"--source", filepath.Base(pdfTmp),
		"--source-xml", filepath.Base(xmlTmp),
		"--format", "zf",
		"--profile", "e",
		"--version", "2",
	}

	// Dynamically attach attachments as INDIVIDUAL arguments
	if len(rInput.Attachments) == 0 {
		args = append(args, "--no-additional-attachments")
	} else {
		for _, attach := range rInput.Attachments {
			attachBase := filepath.Base(attach)
			tmpAttach := filepath.Join(tmpDir, attachBase)
			writer.CopyFile(attach, tmpAttach)
			args = append(args, "--attachment", filepath.Base(tmpAttach))
		}
	}

	// Append the output path to the very end
	args = append(args, "--out", filepath.Base(zfTmp))

	// execute
	ctx, cancel := context.WithTimeout(ctx, w.TimeOut)
	defer cancel()

	// Call the command using the "unpack" operator (...)
	cmd := exec.CommandContext(ctx, w.Executable, args...)
	cmd.Dir = tmpDir // Redirect the command's working directory to our temp directory.

	_, stderr := cmd.CombinedOutput()

	if stderr != nil {
		return nil, errors.Wrap(
			errors.JavaEngineError,
			stderr.Error(),
			"writer.native.mustang.render [mustang-cli.jar]",
			stderr,
		)
	}

	// save tmpOutput to origin artefacts directory, temp Dir will be deleted
	targetPath := filepath.Join(w.JobArtefactsPath, filepath.Base(rInput.OutputPath))
	rOutput = &writer.RenderOutput{
		JobID:      rJob.ID,
		OutputPath: targetPath,
	}
	return rOutput, writer.CopyFile(zfTmp, targetPath)
}

func (w *ZugferdWriter) nativeExtract(
	ctx context.Context, rJob *job.RenderJob,
	rInput *writer.RenderInput) (rOutput *writer.RenderOutput, err error) {

	if err := os.Rename(w.JobTmpFile, w.JobJsonFile); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot rename jobfile %s in %s", w.JobTmpFile, w.JobJsonFile),
			"engines.zugferd.renderNative.nativeExtract",
			err,
		)
	}

	// 1. Create an isolated temp directory
	tmpDir, err := os.MkdirTemp("", "nexfact-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir) // Cleanup after the run

	pdfTmp := filepath.Join(tmpDir, filepath.Base(rInput.PdfPath))
	fxBase := filepath.Base(rInput.OutputPath)
	writer.CopyFile(rInput.PdfPath, pdfTmp)

	cmd := exec.CommandContext(ctx, w.Executable,
		"-jar", w.Engine,
		"--action", "extract",
		"--source", filepath.Base(pdfTmp),
		"--out", fxBase,
	)

	cmd.Dir = tmpDir // Redirect the command's working directory to our temp directory.

	ctx, cancel := context.WithTimeout(ctx, w.TimeOut)
	defer cancel()

	_, stderr := cmd.CombinedOutput()
	if stderr != nil {
		return nil, errors.Wrap(
			errors.JavaEngineError,
			stderr.Error(),
			"writer.native.zugferd.nativeExtract [mustang-cli.jar]",
			stderr,
		)
	}

	// save tmpOutput to origin artefacts directory, temp Dir will be deleted
	targetPath := filepath.Join(w.JobArtefactsPath, filepath.Base(rInput.OutputPath))
	rOutput = &writer.RenderOutput{
		JobID:      rJob.ID,
		OutputPath: targetPath,
	}
	return rOutput, writer.CopyFile(filepath.Join(tmpDir, fxBase), targetPath)
}

func (w *ZugferdWriter) nativeValidate(ctx context.Context, rJob *job.RenderJob,
	rInput *writer.RenderInput) (rOutput *writer.RenderOutput, err error) {

	if err := os.Rename(w.JobTmpFile, w.JobJsonFile); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot rename jobfile %s in %s", w.JobTmpFile, w.JobJsonFile),
			"engines.zugferd.renderNative.nativeValidate",
			err,
		)
	}

	// 1. Create an isolated temp directory
	tmpDir, err := os.MkdirTemp("", "nexfact-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir) // Cleanup after the run

	if rInput.PdfPath != "" {
		pdfTmp := filepath.Join(tmpDir, filepath.Base(rInput.PdfPath))
		writer.CopyFile(rInput.PdfPath, pdfTmp)
		return w.privateNativeValidate(ctx, pdfTmp, tmpDir, rJob)
	}

	if rInput.XMLPath != "" {
		xmlTmp := filepath.Join(tmpDir, filepath.Base(rInput.XMLPath))
		writer.CopyFile(rInput.XMLPath, xmlTmp)
		return w.privateNativeValidate(ctx, xmlTmp, tmpDir, rJob)
	}

	return nil, errors.New(
		errors.EmptyInput,
		"either a pdf Source or a xml source must be set!",
		"engines.nativeValidate",
	)
}

func (w *ZugferdWriter) privateNativeValidate(
	ctx context.Context, fileToValidate string, tmpDir string,
	rJob *job.RenderJob) (rOutput *writer.RenderOutput, err error) {

	cmd := exec.CommandContext(ctx, w.Executable,
		"-jar", w.Engine,
		"--action", "validate",
		"--source", filepath.Base(fileToValidate),
		"--log-as-pdf")

	cmd.Dir = tmpDir // Redirect the command's working directory to our temp directory.

	ctx, cancel := context.WithTimeout(ctx, w.TimeOut)
	defer cancel()

	stdout, stderr := cmd.CombinedOutput()
	if stderr != nil {
		if !strings.Contains(string(stdout), "<summary status=") {
			return nil, errors.Wrap(
				errors.JavaEngineError,
				stderr.Error(),
				"writer.native.mustang.render [mustang-cli.jar]",
				stderr,
			)
		}
	}

	// save tmpOutput to origin artefacts directory, temp Dir will be deleted
	fileName := filepath.Base(fileToValidate)
	tmpBase := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName)) + "_result.pdf"
	jobBase := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName)) + ".result.pdf"
	tmpResult := filepath.Join(tmpDir, tmpBase)

	targetPath := filepath.Join(w.JobArtefactsPath, jobBase)
	rOutput = &writer.RenderOutput{
		JobID:      rJob.ID,
		OutputPath: targetPath,
	}

	return rOutput, writer.CopyFile(tmpResult, targetPath)
}
