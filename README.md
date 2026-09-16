
<sup>🌍 **Language:** 🇩🇪 [German →](README.de.md)</sup>

---
|[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)| |
|----|----|
|![GSF-NexFact-Logo](nexfact.png)| ***GSF-NexFact***<br>A modular system for generating legally compliant electronic invoices (ZUGFeRD/Factur-X). It combines the flexibility of Go with the layout power of LibreOffice and the standards compliance of the Mustang project.|

<sup>***GSF*** stands for ***Go Small Frameworks*** — minimalist tools for robust applications.</sup>

---

### Repository

This project is developed and maintained on Codeberg.

**Canonical repository:**

> https://codeberg.org/tiny-frameworks/nexfact

GitHub is an automated mirror of the Codeberg repository:

**GitHub mirror:**

> https://github.com/tiny-frameworks/nexfact

---

## GSF-NexFact

GSF-NexFact is a modular system designed for flexible generating, analyzing, and validating **ZUGFeRD and Factur-X e-invoices** from JSON data.  
It transforms inputs into **XML, standard PDFs, or hybrid ZUGFeRD invoices** (including attachments) using a **file-based architecture** with clear conventions (*Convention over Configuration*).

---

## Architecture & Core Principle: Convention over Configuration

The entire system relies on a strictly file-based data exchange using a defined directory
structure (`NEXFACT-env`). The Go Orchestrator manages state transitions,
while specialized workers (running natively or within isolated containers)
monitor and process the directories. A standardized `RenderJob` file acts as
the bridge between Go, the LibreOffice Basic macro, and the Mustang engine (Java).

```
            +-----------------------+
            |  nexfact-orchestrator |
            +-----------+-----------+
                        | (Handles orchestration, transforms
                        |  input-json into internal MasterZugFerd-Struct)
                        v
            +-----------------------+
            |      nexfact-env      | <--- Common state (Directories)
            +-----------+-----------+
                        |
       +----------------+-------------+
       v                              v
+------------------------+     +-----------------------+
|   PDF via LibreOffice  |     |  ZUGFeRD via Mustang  |
|  native or container   |     |  native or container  | (Workers)
+------------------------+     +-----------------------+

```

---

## Repository Overview, Components (Modules), Documentation

This repository is organized as a **multi-module repository**:

- [`api`](./api/README.md): Core interfaces, types, structs, centralized error handling (e.g., `RenderJob`, `Error`, `Writer`)
- [`engines`](./engines/README.md): Processing engines (`writer`), synchronization &amp; worker control 
- [`orchestrator`](./orchestrator/README.md): Control unit, entry point for the Go
- [`examples`](./examples/main.go): Ready-to-run code examples
- [`env`](./env/README.md): Directory structure & configurations
- [`deployments`](./deployments/README.md): Container repos for Docker/Podman (PDF/Mustang)


### Helper Modules

- [`cache`](https://codeberg.org/tiny-frameworks/nexutils/cache): Minimalistic caching module
- [`logger`](https://codeberg.org/tiny-frameworks/nexutils/logger): Multi-logger (Console/File) based on `log/slog`

---

## Quick Start & Integration (Go API)

A simple `main.go` snippet serving as an entry point. This setup utilizes default values:
- The `env` structure is located right next to the `main.go` in the current working directory.
- The `default` seller profile is used (found in `env/sellers/default`).
- The provider type is set to `native` (configured in `env/system/systm.yaml` via `zugferd_provider=native` and `pdf_provider=native`).

You can find the documented, working example at [examples](./examples/main.go). 
You can find the containerized counterpart to this `main.go` at [deployments](./deployments.README.md). 
You can find a detailed and complex template in the [orchestrator](./orchestrator/orchestrator_suite_test.go).


```go
package main

import (
	"errors"
	"os"
	"os/exec"
	"runtime"

	"codeberg.org/tiny-frameworks/nexfact/api/logger"
	"codeberg.org/tiny-frameworks/nexfact/api/nexerrors"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator"
)

func main() {

	// Required only for the local monorepo example:
	os.Setenv(orchestrator.EnvRootKey, "../env")

	// In production, the orchestrator searches for the env-directory via:
	// 1: orchestrator.NewDefault():
	//   1.1: environment variable: nexfact_ENVROOT, if empty:
	//   1.2: default directory position (env-dir in working directory)
	// 2: orchestrator.New(env-path)

	orch, err := orchestrator.NewDefault()
	if err != nil {
		nexerrors.LogError(logger.Logger, err)
		os.Exit(nexerrors.GetExitCode(err))
	}

	// 2. Execute job (expects JSON payload for the toolchain)
	jsonPayload := jsonJob()
	result, err := orch.RunJsonJob(jsonPayload)
	if err != nil {
		nexerrors.LogError(logger.Logger, err)
		os.Exit(nexerrors.GetExitCode(err))
	}

	logger.Logger.Info("Job succeeded", "JobID:", result.JobID, "OutputPath", result.OutputPath)
	if err := openPDF(result.OutputPath); err != nil {
		nexerrors.LogError(logger.Logger, err)
		os.Exit(nexerrors.GetExitCode(err))
	}
	os.Exit(0)
}

func openPDF(filePath string) error {
	// go-code to open the generated pdf-file	
	// see: examples/main.go
}

func jsonJob() []byte {
	// bytestream for a minimal sample zugferd-invoice
	// see: examples/main.go
}

```

---

## Organizational & Standards

* **Copyright:** © 2026 Georg Hagn.
* **Repository:** `codeberg.org/tiny-frameworks/nexfact`
* **License:** Apache License, Version 2.0.

*GSF-NexFact is an independent open-source project and is not affiliated with any corporation of a similar name.*

---

## Contact & Support

For inquiries, architectural discussions, or security issues, please contact:

📧 **georghagn [at] tiny-frameworks.io**

---
