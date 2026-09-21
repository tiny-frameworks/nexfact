
## nexfact/orchestrator
<sup>the *orchestrator module*, member of **GSF-nexfact ZUGFeRD** family.<sup>

---

The `orchestrator` module is the central control unit and main entry point for Go applications and microservices using **GSF-NexFact**.
It manages the high-level invoice generation lifecycle, transforms incoming JSON payloads or RPC network requests into
internal domain models, resolves multi-tenant configurations, stages binary assets, and delegates processing to
specialized `writer` engines.

---

## Core Responsibilities

* **Configuration Management:** Initializes system parameters (`system.yaml`) and tenant settings (`seller.yaml`).
* **Payload Validation & Preprocessing:** Transforms incoming JSON payloads into the strongly typed `MasterZugFerd` domain structure. For WebService requests, an automated **Preprocessor** transparently stages Base64-encoded files or remote HTTP URLs into isolated temporary workspaces.
* **Pipeline Execution:** Coordinates rendering steps across engines (e.g., XML generation $\rightarrow$ PDF layout $\rightarrow$ ZUGFeRD PDF injection).
* **Engine Delegation:** Uses `api/writer` factory logic to dynamically instantiate and invoke required engines (`native` vs. `container`).
* **P2P / WebService Exposure:** Exposes high-performance JSON-RPC 2.0 endpoints over WebSockets for distributed or cross-language worker environments.
---

## Module Directory Structure

```text
orchestrator/
├── config/                  # Configuration loaders (system.yaml, seller.yaml, environment paths)
├── orchestrator.go          # Core Orchestrator struct, constructors, and RunJsonJob entry point
├── generateBase.go          # Pipeline step: Pure XML / PAR generation
├── generatePdf.go           # Pipeline step: Layout PDF generation via LibreOffice
├── generateZugferd.go       # Pipeline step: Hybrid ZUGFeRD / Factur-X assembly via Mustang
└── orchestrator_test.go     # Comprehensive integration test suite (local & service mode)

```

---

## Architecture & Data Flow

The orchestrator abstracts the complex multi-step rendering process behind a clean, single-method API (`RunJsonJob`).

```text
               +-----------------------------------+  
               | In-Process Call: RunJsonJob()     |  
               +-----------------+-----------------+  
                                 |
                                 v
             +---------------------------------------+
             | 1. Parse JSON -> MasterZugFerd Struct |
             +-------------------+-------------------+
                                 |
                                 v
             +---------------------------------------+
             | 2. Resolve Tenant & System Configs    |
             +-------------------+-------------------+
                                 |
           +---------------------+---------------------+
           |                     |                     |
           v                     v                     v
 +--------------------+ +--------------------+ +--------------------+
| generateBase()     | | generatePdf()      | | generateZugferd()  |
| (XML Generation)   | | (LibreOffice PDF)  | | (Mustang Assembly) |
+----------+---------+ +----------+---------+ +----------+---------+
           |                      |                     |
           +----------------------+---------------------+
                                  |
                                  v
              +---------------------------------------+
              | 3. Finalize Output & Return OutputPath |
              +---------------------------------------+

```

---

## Job Payload Specification (`invoice.json`)

The `orchestrator` consumes JSON payloads. Over direct `RunJsonJob` calls, file properties expect local file paths. 

### Payload Structure Overview

| Section | Mandatory? | Description |
| --- | --- | --- |
| `mandant.seller` | No | Overrides target tenant folder (e.g., `"testSeller"` $\rightarrow$ `env/sellers/testSeller`). |
| `options.queue` | No | Target pipeline queue (`"pdf"`, `"facturx"`, `"base"`). Defaults to `"zugferd"`. |
| `options.template` | No | Path, URL, or Base64 object for a custom LibreOffice template (`.ott`). |
| `options.attachments` | No | Array of file paths, Base64 objects, or URLs to embed into the ZUGFeRD PDF. |
| `invoice.currency` | Yes | ISO 4217 currency code (e.g., `"EUR"`). |
| `invoice.buyer` | Yes | Buyer address and metadata. |
| `invoice.items` | Yes | Array of line items (quantity, net price, tax rates). |
| `invoice.vats` | Yes | Grouped VAT totals required for tax compliance. |
| `invoice.totals` | Yes | Calculated totals (`line_total_amount`, `tax_total_amount`, `grand_total_amount`). |

---

### Payload Examples

#### In-Process / File-Path Payload (`min-invoice.json`)

Used for direct `orch.RunJsonJob()` execution on local filesystems:

```json
{
  "invoice": {
    "currency": "EUR",
    "invoice_id": "RE-0814-min",
    "issue_date": "20260105",
    "service_date": "20260104",
    "buyer": {
      "id": "471115",
      "name1": "Franz Fröhlich AG",
      "street": "Musterstrasse 64",
      "zip": "01221",
      "city": "Musterdorf",
      "country": "Deutschland"
    },
    "items": [
      {
        "pos": 1,
        "description": "Go Entwicklungsservice",
        "quantity": 10.5,
        "unit": "Stunde",
        "net_price": 95.00,
        "line_total": 997.5,
        "tax_rate": 19
      }
    ],
    "vats": [
      {
        "pos": 1,
        "id": "S",
        "description": "USt",
        "vat_base": 997.50,
        "vat_amount": 189.53,
        "tax_rate": 19
      }
    ],
    "totals": {
      "line_total_amount": 997.50,
      "tax_total_amount": 189.53,
      "grand_total_amount": 1187.03
    }
  }
}

```

---

## API Entry Points

```go
// 1. Initialize Orchestrator with default environment resolution
orch, err := orchestrator.NewDefault()

// 2. Execute a complete invoice rendering pipeline in-process
result, err := orch.RunJsonJob(jsonBytes)

```

---

## Integration & WebService Testing

The `orchestrator` contains a comprehensive integration test suite verifying both in-process engine pipelines and networked WebService execution.

| Test File | Engine / Target | Description |
| --- | --- | --- |
| **`orchestrator_test-suite_native_test.go`** | `NATIVE` | Direct execution without containers or network overhead. |
| **`orchestrator_suite_test.go`** | `NATIVE` & `CONTAINER` | Full integration suite with Docker/Podman container isolation. |

--- 
### 1. Native Testing (Fast & Lightweight, Recommended for Getting Started)

If you do not have Docker or Podman installed—or simply want a lightweight, fast local test run—use the native test suite.

```bash
go test -v -tags=native ./...
```

### 2. Full Integration Testing

To run the full suite including containerized provider tests (requires Docker or Podman):

```bash
go test -v ./...
```

---
