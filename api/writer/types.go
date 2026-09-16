// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package writer

import (
	"time"
)

type RenderInput struct {
	Provider      string
	BytesToRender []byte
	ParPath       string
	OttPath       string
	PdfPath       string
	XMLPath       string
	ArtefactsPath string
	RootPath      string
	Attachments   []string
	OutputPath    string
}

type RenderOutput struct {
	JobID      string
	OutputPath string
	ParPath    string
	StatusPath string
	ODTPath    string
	FXPath     string
	PDFPath    string
}

// AbstractWriter is Base-Struct for spezialized writer, should never be instantiated
type AbstractWriter struct {
	Provider   string
	Engine     string
	Executable string
	TimeOut    time.Duration
	Heartbeat  time.Duration

	JobArtefactsPath string

	JobJsonFile   string
	JobResultFile string
	JobTmpFile    string
}

// Working Env Definition
type WriterKind string

const (
	NativeProvider    string = "native"
	ContainerProvider string = "container"
	BaseProvider      string = "native"
	//WSProvider        string = "webservice"

	PdfEngine     string = "loEngine.odt"
	ZugferdEngine string = "mustang-cli.jar"
	BaseEngine    string = "none"

	BaseKind    WriterKind = "base"
	PdfKind     WriterKind = "pdf"
	ZugferdKind WriterKind = "zugferd"

	PdfExecutable     string = "soffice"
	ZugferdExecutable string = "java"
)

const (
	FileMask      = 0644
	DirectoryMask = 0755
	PollIntervall = 200 * time.Millisecond
)

// ----------------------------
// OutputFileType beschreibt eine generierte Datei oder einen Anhang.
type OutputFileType struct {
	Name     string `json:"name"`      // "factur-x.xml", "rechnung.pdf", ...
	MimeType string `json:"mime_type"` // "application/pdf", "text/xml", ...
	Size     int64  `json:"size"`      // Filesize in Bytes

	// Entweder Direct Payload (für normale/kleine Dateien)...
	Data []byte `json:"data,omitempty"` // Wird im JSON automatisch Base64-kodiert

	// ...oder Referenz für große Dateien (wird später genutzt)
	URL string `json:"url,omitempty"` // Download-URL oder Pfad
}

// JobRPCResponse bildet das Gesamtergebnis des Jobs ab.
type JobRPCResponse struct {
	JobID       string           `json:"job_id"`
	OutputFile  OutputFileType   `json:"output_file"`           // Hauptergebnis (PDF, XML, etc.)
	Attachments []OutputFileType `json:"attachments,omitempty"` // Optionale Zusatzdokumente/Belege
}
