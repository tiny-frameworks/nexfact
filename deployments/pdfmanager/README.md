[![Codeberg Repository](https://img.shields.io/badge/Codeberg-nexfact-blue?logo=codeberg&logoColor=white)](https://codeberg.org/tiny-frameworks/nexfact)
[![GHCR](https://img.shields.io/badge/Container-ghcr.io-blue?logo=github)](https://github.com/tiny-frameworks/nexfact/pkgs/container/nexfact-pdfmanager)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://codeberg.org/tiny-frameworks/nexfact/src/branch/main/LICENSE)

## nexfact-pdfmanager
<sup>the *podman container pdfmanager*, member of **GSF-nexfact ZUGFeRD** family.</sup>

---

This deployment provides a headless **LibreOffice** container runtime designed for automated PDF generation
from document templates (`.ott`) within the **GSF-nexfact** ecosystem.

A specialized Python manager script (`nexfact_manager.py`) orchestrates LibreOffice, passes job parameters
securely, and executes the LibreOffice Basic macro container (`loContainerEngine.odt`).

---

## Directory Structure

```text
pdfmanager/
├── resources/
│   ├── lo_user_profile/       # Optimized LibreOffice user profile for headless operation
│   ├── loContainerEngine.odt  # LibreOffice macro container for PDF rendering
│   └── nexfact_manager.py     # Python manager script monitoring the job queue
├── Dockerfile                 # Container image specification (compatible with Podman & Docker)
├── podman_build.sh            # Convenience build script
└── LICENSE

```

---

## Quick Start / Container

Das OCI-Container-Image wird über die GitHub Container Registry bereitgestellt:

```bash
podman pull ghcr.io/tiny-frameworks/nexfact-pdfmanager:latest
```

---

## Components & Configuration

* **`lo_user_profile/`**: Pre-configured LibreOffice profile stripped of UI overhead to ensure fast startup times and lock-free headless rendering.
* **`nexfact_manager.py`**: Monitors the mounted input directory, handles atomic file locking (`job.json` $\rightarrow$ `job.process` $\rightarrow$ `job.done`), and triggers the macro execution.

### Supported Environment Variables

| Variable | Default | Description |
| --- | --- | --- |
| `NEXFACT_MODE` | `watch` | Execution mode (`watch` for continuous polling or `once` for single job execution) |
| `WATCH_DIR` | `/data/in` | Target directory inside the container monitored for incoming `job.json` files |

---

## Build & Run Instructions

The container expects the `pdfs` section of a tenant's `nexfact-env` directory structure to be mounted into `/data`.

### 1. Build Container Image

```bash
podman build -t nexfact-pdfmanager .

```

### 2. Run Container (Watcher Mode)

Execute this command from within the target tenant directory (e.g., `nexfact-env/sellers/default/`):

```bash
podman run -d --name nexfact-pdfwatcher \
  -e NEXFACT_MODE=watch \
  -e WATCH_DIR=/data/in \
  -v ./data/pdfs/in:/data/in:Z \
  -v ./data/artefacts:/data/artefacts:Z,shared \
  -v ./data/pdfs/logs:/data/logs:Z \
  -v ./templates:/data/templates:Z,ro \
  nexfact-pdfmanager

```

Important: The ./templates volume mount is mandatory. The containerized LibreOffice engine relies
on the .ott templates stored in the tenant's templates/ directory to generate the PDF layout.

> **Note on SELinux (`:Z` flags):** The `:Z` flag ensures correct SELinux relabeling when running with Podman on Fedora, RHEL, or Linux Mint/Ubuntu systems with SELinux enabled.

---

