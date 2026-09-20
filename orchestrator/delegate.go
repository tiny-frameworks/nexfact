// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/internal/preprocessor"
	"codeberg.org/tiny-frameworks/nexutils/errors"
	"codeberg.org/tiny-frameworks/nexutils/p2p/rpc"
)

type OrchestratorDelegate struct {
	rpc.DefaultNexDelegate
	log *slog.Logger
	o   *Orchestrator
}

type WebServiceManager struct {
	rpc.DefaultNexDelegate
	log *slog.Logger
	o   *Orchestrator
}

func NewDelegate(orch *Orchestrator, addr string, heartbeat time.Duration, lgr *slog.Logger) *rpc.Node {

	oNode := rpc.NewNode(rpc.Options{
		Addr:              addr,
		HeartbeatInterval: heartbeat,
		Delegate: &OrchestratorDelegate{
			log: lgr,
			o:   orch,
		},
	})

	return oNode
}

func (d *OrchestratorDelegate) OnRequest(ctx context.Context, peer *rpc.Peer, method string, params []byte) (any, error) {
	switch method {
	case "nexfact.process":
		return d.privateProcess(peer, params)

	case "nexfact.echo":
		return d.privateEcho(peer, params)

	default:
		return d.DefaultNexDelegate.OnRequest(ctx, peer, method, params)
	}
}

func (d *OrchestratorDelegate) privateEcho(peer *rpc.Peer, params []byte) (any, error) {
	var echo string

	if err := json.Unmarshal(params, &echo); err != nil {
		return nil, &rpc.JsonRPCerror{Code: rpc.InvalidParams, Message: rpc.StdError[rpc.InvalidParams]}
	}
	d.log.Info("[Echo] Received echo request from partner", "params", params, "peer", peer.RemoteAddr())
	return "Echo: " + echo, nil
}

func (d *OrchestratorDelegate) privateProcess(peer *rpc.Peer, params []byte) (any, error) {
	var jsonInvoice string

	if err := json.Unmarshal(params, &jsonInvoice); err != nil {
		d.log.Error("nexfact.process: invalid parameters", "err", err, "peer", peer.RemoteAddr())
		nerr := errors.Wrap(errors.InvalidValue,
			"nexfact.process: invalid parameters",
			"delegate.privateProcess", err)
		return nil, rpc.NewRPCErrorFromNexError(rpc.InvalidParams, nerr)
	}
	d.log.Info("nexfact.process", "peer", peer.RemoteAddr())

	// 1. Staging Directory anlegen & Fehler abfangen
	tempDir, err := os.MkdirTemp("", "nexfact_stage_*")
	if err != nil {
		d.log.Error("failed to create temp staging directory", "err", err)
		nerr := errors.Wrap(errors.ExecutionError,
			"failed to create temporary staging directory",
			"delegate.privateProcess", err)
		return nil, rpc.NewRPCErrorFromNexError(rpc.InternalError, nerr)
	}
	defer os.RemoveAll(tempDir)

	// 2. Preprocessing mit typsicherer errors.As Conversion
	oJSON, err := preprocessor.PreProcess(tempDir, []byte(jsonInvoice))
	if err != nil {
		var nerr *errors.Error
		if !errors.As(err, &nerr) {
			nerr = errors.Wrap(errors.ExecutionError,
				"preprocessing failed",
				"delegate.privateProcess", err)
		}
		return nil, rpc.NewRPCErrorFromNexError(rpc.InternalError, nerr)
	}

	// 3. Orchestrator-Ausführung
	out, err := d.o.RunJsonJob(oJSON)
	if err != nil {
		var nerr *errors.Error
		if !errors.As(err, &nerr) {
			nerr = errors.Wrap(errors.ExecutionError,
				"orchestrator job execution failed",
				"delegate.privateProcess", err)
		}
		return nil, rpc.NewRPCErrorFromNexError(rpc.InternalError, nerr)
	}

	// 4. Ausgabedatei einlesen (Serverseitiger Fehler -> InternalError)
	outBytes, err := os.ReadFile(out.OutputPath)
	if err != nil {
		d.log.Error("failed to read output file", "path", out.OutputPath, "err", err)
		nerr := errors.Wrap(errors.ReadError,
			fmt.Sprintf("cannot read output file: %s", out.OutputPath),
			"delegate.privateProcess", err)
		return nil, rpc.NewRPCErrorFromNexError(rpc.InternalError, nerr)
	}

	mimeType := http.DetectContentType(outBytes)
	return writer.JobRPCResponse{
		JobID: out.JobID,
		OutputFile: writer.OutputFileType{
			Name:     filepath.Base(out.OutputPath),
			MimeType: mimeType,
			Size:     int64(len(outBytes)),
			Data:     outBytes,
		},
	}, nil
}
