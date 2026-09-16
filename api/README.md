
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

### 3. `api/nexerrors`

Implements centralized, structured error handling tailored for the nexgate ecosystem.

```go
type Error struct {
    Code    Code
    Message string
    Path    string
    Cause   error
}

```

**Key Functions:**

* `New(code Code, msg string, path string) *Error`: Creates a new structured error.
* `Wrap(code Code, msg string, path string, cause error) *Error`: Wraps an existing error into the nexgate error domain.
* `GetExitCode(err error) int`: Traverses the error chain and unwraps the appropriate OS exit code.
* `LogError(log *slog.Logger, err error)`: Unpacks `nexerrors.Error` attributes (`code`, `path`, `cause`) and logs them as structured `slog` key-value pairs.

---

### 4. `api/logger`

Provides a thread-safe MultiHandler logger on top of Go's native `log/slog`. It forwards log entries simultaneously to the console and a rotating log file.

**Usage:**

```go
import "codeberg.org/tiny-frameworks/nexfact/api/logger"

func main() {
    // Initialize logging output to file
    logger.SetupLogging("./main.log")

    logger.Logger.Info("Application initialized successfully")
}

```

---

## Design Principle

Packages inside `api/` must **never depend on higher-level packages** (like `orchestrator` or `engines`). They serve strictly as lightweight, dependency-free contracts to avoid circular dependencies in the monorepo.

---


