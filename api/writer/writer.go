// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

type Writer interface {
	Render(ctx context.Context, rJob *job.RenderJob, rInput *RenderInput, finalize bool) (*RenderOutput, error)
	Prepare(job *job.RenderJob, input *RenderInput) error
	Finalize(rJob *job.RenderJob, rInput *RenderInput, sourcePath string, lastStep bool) (*RenderOutput, error)
}

func (a *AbstractWriter) Prepare(rJob *job.RenderJob, rInput *RenderInput) error {

	// build working path
	// atomic write: WriteFile is not atomic, but rename is therefore write
	// jobFile as tmpFile. Caller must rename
	a.JobArtefactsPath = filepath.Join(rInput.ArtefactsPath, rJob.ID)
	a.JobJsonFile = filepath.Join(rInput.RootPath, "in", rJob.ID+".json")
	a.JobResultFile = filepath.Join(rInput.RootPath, "in", rJob.ID+".done")
	a.JobTmpFile = a.JobJsonFile + ".tmp"

	if _, err := os.Stat(a.JobArtefactsPath); os.IsNotExist(err) {
		// We only create it if it absolutely does not exist.
		if err := os.Mkdir(a.JobArtefactsPath, DirectoryMask); err != nil {
			return errors.Wrap(
				errors.WriteError,
				"cannot create ArtefactsPath",
				"writer.AbstractWriter.prepare",
				err,
			)
		}
	} else {
		if err != nil {
			// If another error occurs (e.g., "Permission Denied" due to the Podman mount)
			return errors.Wrap(
				errors.WriteError,
				"ArtefactsPath access conflict (Podman lock?)",
				"writer.AbstractWriter.prepare",
				err,
			)
		}
	}

	if err := a.WriteJob(rJob); err != nil {
		return err
	}

	return nil
}

func (a *AbstractWriter) WriteJob(rJob *job.RenderJob) error {

	data, err := json.MarshalIndent(rJob, "", "  ")
	if err != nil {
		return errors.Wrap(
			errors.WriteError,
			"cannot marshall job json file",
			"writer.AbstractWriter.prepare",
			err,
		)
	}

	if err := os.WriteFile(a.JobTmpFile, data, FileMask); err != nil {
		return errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("cannot write job json file %s to in directory", a.JobTmpFile),
			"writer.AbstractWriter.prepare",
			err,
		)
	}

	return nil
}

func (a *AbstractWriter) Finalize(rJob *job.RenderJob, rInput *RenderInput, sourcePath string, lastStep bool) (*RenderOutput, error) {
	if !lastStep {
		return nil, nil
	}

	if err := CopyFile(sourcePath, rInput.OutputPath); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("can not move file %s", sourcePath),
			"writer.base.writer.render.Finalize()",
			err,
		)
	}

	// move attachments to output folder
	if rJob.Queue == job.QueueExtract {
		for _, attach := range rJob.Attachments {
			attachBase := filepath.Base(attach)
			attachArte := filepath.Join(a.JobArtefactsPath, attachBase)
			outDir := path.Dir(rInput.OutputPath)
			attachOutp := filepath.Join(outDir, attachBase)
			if err := CopyFile(attachArte, attachOutp); err != nil {
				return nil, errors.Wrap(
					errors.WriteError,
					fmt.Sprintf("can not copy file %s to %s", attachArte, attachOutp),
					"writer.base.writer.render.Finalize()",
					err,
				)
			}
		}
	}

	if err := os.RemoveAll(a.JobArtefactsPath); err != nil {
		return nil, errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("can not removeAll from %s ", a.JobArtefactsPath),
			"writer.base.writer.render.Finalize()",
			err,
		)
	}
	return &RenderOutput{
		OutputPath: rInput.OutputPath,
		JobID:      rJob.ID,
	}, nil
}

// Helper-func CopyFile: copies a file from src to dst
func CopyFile(src, dst string) error {

	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening the sourcefile: %w", err)
	}
	defer sourceFile.Close()

	// 2. create targetfile (overwrites existing files!)
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating the targetfile: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("error copying to targetfile: %w", err)
	}

	destFile.Sync()

	return nil
}

// func Check verifies the heartbeat file.
// It returns 'true' if everything is okay, and 'false', along with an error, if the engine is stuck.
func CheckHeartbeat(heartbeatFile string, timeout time.Duration) (bool, error) {
	// Read content from file
	data, err := os.ReadFile(heartbeatFile)
	if err != nil {
		return false, errors.Wrap(
			errors.ReadError,
			fmt.Sprintf("could not read Heartbeat-File: %s", heartbeatFile),
			"api.writer.CheckHeartbeat", err)
	}

	// parse timestamp as UTC
	lastHeartbeat, err := time.Parse(time.RFC3339, strings.TrimSpace(string(data)))
	if err != nil {
		return false, errors.Wrap(
			errors.ReadError,
			fmt.Sprintf("could not parse lastHeartbeat: %v", lastHeartbeat),
			"api.writer.CheckHeartbeat", err)
	}

	// calculate difference to UTC-NOW
	diff := time.Since(lastHeartbeat.UTC())
	if diff > timeout {
		return false, errors.New(
			errors.ContainerTimeout,
			fmt.Sprintf("heartbeat timeout reached: last=%v, diff=%v", lastHeartbeat, diff),
			"api.writer.CheckHeartbeat")
	}
	return true, nil
}

// Registry for writers:
// The shared blueprint for all factory functions
type FactoryFunc func(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (Writer, error)

// The central repository for all available engines
var registry = make(map[WriterKind]FactoryFunc)

// Register is called by the engines to sign up.
func Register(kind WriterKind, factory FactoryFunc) {
	registry[kind] = factory
}

// CreateWriter is now the universal factory that the orchestrator calls.
func CreateWriter(kind WriterKind, provider string, engine string, timeout time.Duration, heartbeat time.Duration) (Writer, error) {
	factory, exists := registry[kind]
	if !exists {
		return nil, errors.New(
			errors.UnsupportedWriterKind,
			fmt.Sprintf("unsupportet writerkind found: %v", kind),
			"api.writer.CreateWriter")
	}
	return factory(provider, engine, timeout, heartbeat)
}
