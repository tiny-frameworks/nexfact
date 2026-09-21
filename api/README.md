
## nexfact/api
<sup>the *api module*, part of **GSF-nexfact ZUGFeRD**, member of the **tiny-frameworks** family<sup>

---

The `api` module forms the foundational layer of **nexfact**. It defines global structs, interfaces, domain models, structured error types, and shared utilities (such as logging) to guarantee consistency and reusability across all nexgate modules.

---

## Package Overview

### 1. `api/job`
Defines the core data contracts and job state containers used throughout the toolchain.

* **`RenderJob`**: Central object tracking job lifecycle, IDs, and processing state.
* **`RenderInput` / `RenderOutput`**: Standardized input parameters and resulting artifacts (e.g., PDF path, status).

---

### 2. `api/writer`

Defines the uniform `Writer` interface implemented by all specialized engines (PDF conversion, ZUGFeRD/Factur-X injection)
and manages their registration.

**Core Interface:**
```go
type Writer interface {
    Write(ctx context.Context, job *job.RenderJob) (*job.RenderOutput, error)
}

```
*Allows seamless swapping between `native` and `containerized` execution drivers.*

**Engine Registry:**
The package includes a thread-safe `Registry` that maps provider keys (e.g., `"native"`, `"container"`, `"podman"`)
to their respective `Writer` implementations. This decouples the orchestrator from concrete engine implementations.

* **`Register(name string, w Writer)`**: Registers a new engine writer under a unique provider key.
* **`Get(name string) (Writer, error)`**: Retrieves the requested writer instance based on configuration flags (e.g., from `system.yaml`).

*Usage Example:*

```go
// Engines register themselves during initialization
writer.Register("native", nativePdfEngine)

// Orchestrator retrieves the engine dynamically at runtime
pdfWriter, err := writer.Get("native")

```

---

## Design Principle

Packages inside `api/` must **never depend on higher-level packages** (like `orchestrator` or `engines`). They serve strictly as lightweight, dependency-free contracts to avoid circular dependencies in the monorepo.

---


