
## nexfact-engines
<sup>the *engines module*, member of **GSF-nexfact ZUGFeRD** family.<sup>

---

The `engines` module provides concrete `Writer` implementations responsible for executing render jobs within the **NexFact** ecosystem. It bridges the gap between high-level orchestration (`orchestrator`) and low-level processing workers (LibreOffice, Mustang, XML generators).

---

## Core Responsibilities

* **Directory Setup:** Preparing job-specific workspace structures inside `nexfact-env` (`in/`, `artefacts/`, `logs/`).
* **Job Contract Generation:** Serializing the `RenderJob` payload into the expected target files (e.g., `job.json`).
* **Execution & Monitoring:** Dispatching tasks to concrete workers based on the configured provider (`native` vs. `container`) and enforcing timeouts and heartbeats.
* **Finalization:** Validating generated artifacts, collecting metrics, and moving final outputs (PDF, Factur-X, hybrid ZUGFeRD) to the target output directory.

---

## Available Engine Package Types

The module contains specialized sub-packages for each rendering target:

| Engine Package | Target Format | Supported Providers | Underlying Technology |
| --- | --- | --- | --- |
| `engines/base` | Pure XML / PAR | `native` | Internal Go templates / XML serialization |
| `engines/pdf` | Layout PDF | `native`, `container` | LibreOffice Basic Macro (`loEngine.odt`) |
| `engines/zugferd` | Hybrid ZUGFeRD / Factur-X PDF | `native`, `container` | Mustang Project CLI (`mustang-cli.jar`) |

---

## Architecture & Lifecycle

All engines extend `writer.AbstractWriter` to reuse shared lifecycle methods (`Prepare`, `buildJobArtefacts`, `Finalize`).

```text
                  +-----------------------+
                  |  orchestrator.Render  |
                  +-----------+-----------+
                              |
                              v
                  +-----------------------+
                  |   zugferd.Render()    |
                  +-----------+-----------+
                              |
                +-------------+-------------+
                |                           |
                v                           v
     +--------------------+       +--------------------+
     | w.Prepare(rJob)    |       | w.buildArtefacts() |
     +----------+---------+       +---------+----------+
                |                           |
                +-------------+-------------+
                              |
                              v
                +---------------------------+
                | switch w.Provider         |
                +-------------+-------------+
                              |
             +----------------+----------------+
             |                                 |
             v                                 v
  +--------------------+             +--------------------+
  |  renderNative()    |             |  renderContainer() |
  |  (Direct Subproc)  |             |  (File Locking)    |
  +----------+---------+             +---------+----------+
             |                                 |
             +----------------+----------------+
                              |
                              v
                  +-----------------------+
                  |  w.Finalize(rJob, ..) |
                  +-----------------------+

```

---

## Engine Registration Pattern

Engines automatically register factory constructors with the central `api/writer` registry upon application startup via Go's `init()` function. This enables dynamic instantiation by the orchestrator based on configuration flags without hardcoupling dependencies.

### Implementation Example (`engines/zugferd`):

```go
package zugferd

import (
    "time"
    "codeberg.org/tiny-frameworks/nexfact/api/writer"
)

// Automatic self-registration in api/writer registry at startup
func init() {
    writer.Register(writer.ZugferdKind, func(provider string, engine string, timeout time.Duration, heartbeat time.Duration) (writer.Writer, error) {
        return New(provider, engine, timeout, heartbeat)
    })
}

```

---

## Usage Example (Orchestrator Integration)

The orchestrator dynamically creates the required engine writer via the registry and invokes the `Render` method:

```go
func (o *Orchestrator) generateZugferd(rJob *job.RenderJob, finalize bool) (*writer.RenderOutput, error) {
    provider := config.SystemParams.Engine.ZugferdProvider // "native" | "container"
    engine   := config.SystemParams.Engine.ZugferdEngine   // "mustang-cli.jar" or container image
    timeout  := config.SystemParams.Engine.CmdTimeout

    // 1. Factory lookup from registry
    zugferdWriter, err := o.CreateWriter(provider, engine, timeout, writer.ZugferdKind)
    if err != nil {
        return nil, err
    }

    // 2. Build input parameters
    renderInput := writer.RenderInput{
        ArtefactsPath: filepath.Join(config.SellerParams.SellerRoot, "data", "artefacts"),
        RootPath:      filepath.Join(config.SellerParams.SellerRoot, "data", "zugferds"),
        OutputPath:    filepath.Join(config.SellerParams.SellerRoot, "outputs", o.ZugferdMaster.Invoice.InvoiceID+".zf.pdf"),
    }

    // 3. Render job (handles native vs container execution transparently)
    return zugferdWriter.Render(o.ctx, rJob, &renderInput, finalize)
}

```

---


