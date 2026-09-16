// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package orchestrator

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/internal/preprocessor"
	"codeberg.org/tiny-frameworks/nexutils/errors"
	"codeberg.org/tiny-frameworks/nexutils/p2p/rpc"
)

// webServiceHandler bundles the dependencies of all RPC methods.
type webServiceHandler struct {
	o   *Orchestrator
	lgr *slog.Logger
}

func newWebServiceHandler(o *Orchestrator, lgr *slog.Logger) *webServiceHandler {
	return &webServiceHandler{
		o:   o,
		lgr: lgr,
	}
}

func (h *webServiceHandler) handleProcess(p *rpc.Peer, req rpc.JsonRPCrequest) (any, *rpc.JsonRPCerror) {
	var jsonInvoice string

	if err := req.UnmarshalParams(&jsonInvoice); err != nil {
		h.lgr.Error("nexgate.process: invalid parameters", "err", err, "peer", p.RemoteAddr())
		nerr := errors.Wrap(errors.InvalidValue,
			"nexgate.process: invalid parameters",
			"orchestrator.setupWebService", err)
		return nil, rpc.NewRPCErrorFromNexError(rpc.InvalidParams, nerr)
	}

	h.lgr.Info("nexgate.process", "peer", p.RemoteAddr())

	// prepare staging directory
	tempDir, _ := os.MkdirTemp("", "nexgate_stage_*")
	defer os.RemoveAll(tempDir)

	// Preprocessing (map-based, writes files and transforms paths)
	orchestratorJSON, err := preprocessor.PreProcess(tempDir, []byte(jsonInvoice))

	if err != nil {
		return nil, rpc.NewRPCErrorFromNexError(rpc.InternalError, err.(*errors.Error))
	}

	// Hand over to core orchestrator
	out, err := h.o.RunJsonJob(orchestratorJSON)

	if err != nil {
		nerr := errors.Wrap(errors.InvalidValue,
			"nexgate.process: invalid parameters",
			"orchestrator.setupWebService", err)
		return nil, rpc.NewRPCErrorFromNexError(rpc.InternalError, nerr)
	}

	outBytes, err := os.ReadFile(out.OutputPath)
	if err != nil {
		nerr := errors.Wrap(errors.ReadError,
			"nexgate.process: ungültige Parameter",
			"orchestrator.setupWebService", err)
		h.lgr.Error("Error reading outputDatei", "path", out.OutputPath, "err", err)
		return nil, rpc.NewRPCErrorFromNexError(rpc.InvalidParams, nerr)
	}

	mimeType := http.DetectContentType(outBytes)
	return writer.JobRPCResponse{
		JobID: out.JobID,
		OutputFile: writer.OutputFileType{
			Name:     filepath.Base(out.OutputPath),
			MimeType: mimeType,
			Size:     int64(len(outBytes)),
			Data:     outBytes},
	}, nil

}

func (h *webServiceHandler) handleEcho(p *rpc.Peer, req rpc.JsonRPCrequest) (any, *rpc.JsonRPCerror) {
	var params string
	req.UnmarshalParams(&params)
	h.lgr.Info("[Echo] Received echo request from partner", "params", params, "peer", p.RemoteAddr())
	return "Echo: " + params, nil
}
