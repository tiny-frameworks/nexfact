
## nexfact-env
<sup>the *environment module*, member of **GSF-nexfact ZUGFeRD** family.</sup>

---

The `nexgate-env` module defines the mandatory runtime environment and directory structure for the **NexFact** ecosystem. It enforces a strict **"Convention over Configuration"** architecture for both native processes and containerized workers.

All components (Go Orchestrator, LibreOffice Basic Macros, Mustang Java Engine, and Container Workers) interact asynchronously by monitoring and manipulating files within this predictable directory layout.

---

## Directory Structure

### 1. Application Directory Structure (Host / System View)

```text
env/
├── sellers/                   # Multi-tenant directory
│   └── [SellerName]/          # Client-specific namespace (e.g., "default")
│       ├── data/
│       │   ├── artefacts/     # Temporary build directory (intermediate PDFs, XMLs)
│       │   ├── base/          # Job directory for pure XML / PAR generation
│       │   │   ├── in/        # Job queue (job.json -> job.process -> job.done)
│       │   │   └── logs/      # Log files, heartbeats, and environment state
│       │   ├── pdfs/          # Job directory for PDF rendering
│       │   │   ├── in/        # Job queue for PDF tasks
│       │   │   └── logs/
│       │   └── zugferds/      # Job directory for hybrid ZUGFeRD / Factur-X invoices
│       │       ├── in/        # Job queue for final assembly
│       │       └── logs/
│       ├── outputs/           # Final, validated end products (e.g., RE-2026-001.pdf)
│       ├── templates/         # Invoice templates (*.ott, *.tpl)
│       └── seller.yaml        # Client-specific configuration
└── system/                    # Global binaries, lookup tables, and system configs
    ├── csv/                   # System-wide reference data (e.g., units.csv, vats.csv)
    ├── loEngine.odt           # LibreOffice Macro container (for native PDF rendering)
    ├── mustang-cli.jar        # Mustang ZUGFeRD Java library (for native XML injection)
    └── system.yaml            # Global system parameters (provider flags, timeouts)

```

*Das Verzeichnis templates/ enthält die master .ott-Vorlagen, die von den PDF-Engines
(sowohl native als auch im Container) benötigt werden, um das Dokument-Layout zu erzeugen.*

### 2. Container Mapped Directory Structure

When running in **Container Mode**, workers (e.g., Podman/Docker containers for PDF or ZUGFeRD processing) do not need access to the entire `env` tree. Instead, a specific task directory is mounted directly into the container's `/data` volume:

**Host Path:** `env/sellers/[SellerName]/data/pdfs/` $\rightarrow$ **Container Path:** `/data/`

```text
/data/                         # Mounted inside the container
├── artefacts/                 # Shared intermediate processing folder
├── in/                        # Input job queue monitored by the container worker
└── logs/                      # Container execution logs and heartbeat output

```

*Example Mount:*

```bash
podman run -v ./env/sellers/default/data/pdfs:/data:Z nexfact/pdfmanager

```

---
## How It Works (Job Lifecycle)

1. **Job Creation:**
The Go `orchestrator` parses the input JSON payload, builds the internal structure, and writes a `job.json` file into the target `in/` directory (e.g., `sellers/default/data/pdfs/in/job.json`).
2. **Job Processing & Execution:**  
   Depending on the configured provider (`native` vs. `container`), the execution strategy differs:
   * **Native Mode:** The Go `writer` directly invokes the local runtime (e.g., executing Java/LibreOffice as a child process) and blocks until completion or timeout.
   * **Container Mode:** The container worker continuously polls its mounted `in/` directory, locks `job.json` by renaming it to **`job.process`**, emits heartbeats into `logs/`, and renames it to **`job.done`** upon completion.
3. **Finalization:**
The Go `orchestrator` validates the output inside `artefacts/`, cleans up temporary state files, and delivers the finalized invoice to the tenant's `outputs/` folder.

```

---
