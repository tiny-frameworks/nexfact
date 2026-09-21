// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

/*
A simple `main.go` snippet serves as the entry point. Default values ​​are used here:
  - The `env` structure is located in the working directory alongside `main.go`.
  - The `default` seller is used (found in `env/sellers/default`).
  - The `native` provider is used (specified in `env/system/systm.yaml` via `zugferd_provider=native` and `pdf_provider=native`).

You can find a detailed code example showing how to integrate nexfact as a library in `orchestrator.orchestrator_test`.
*/
package main

import (
	stdErrors "errors"
	"os"
	"os/exec"
	"runtime"

	"codeberg.org/tiny-frameworks/nexfact/orchestrator"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
	"codeberg.org/tiny-frameworks/nexutils/errors"
	"codeberg.org/tiny-frameworks/nexutils/logger"
)

func main() {

	// Executes the logic and ensures that `defer` calls are guaranteed to run!
	os.Exit(run())

}
func run() int {

	defer config.StopCaches()

	// Required only for the local monorepo example:
	os.Setenv(orchestrator.EnvRootKey, "../env")

	// In production, the orchestrator searches env-directory with:
	// 1: orchestrator.NewDefault():
	//   1.1: environment-variable: NEXFACT_ENVROOT, if empty:
	//   1.2: default directory position (env-dir in working-directory)
	// 2: orchestrator.New(env-path)
	orch, err := orchestrator.NewDefault()

	if err := setupLogger(); err != nil {
		return 1
	}

	if err != nil {
		errors.LogError(logger.Logger, err)
		os.Exit(errors.GetExitCode(err))
	}

	// 2. Execute job (expects JSON payload for the toolchain)
	jsonPayload := jsonJob()
	result, err := orch.RunJsonJob(jsonPayload)
	if err != nil {
		errors.LogError(logger.Logger, err)
		os.Exit(errors.GetExitCode(err))
	}

	logger.Logger.Info("Job succeeded", "JobID:", result.JobID, "OutputPath", result.OutputPath)
	if err := openPDF(result.OutputPath); err != nil {
		errors.LogError(logger.Logger, err)
		os.Exit(errors.GetExitCode(err))
	}
	return 0 // SUCCESS Exit-Code
}

func setupLogger() error {
	cfg := &logger.LoggerConfig{
		Filename:   "example_main.log",
		UseLocking: false,
		Level:      logger.LevelInfo,
	}

	return logger.SetupLogging(cfg)
}

func openPDF(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return stdErrors.New("unsupported operating system: " + runtime.GOOS)
	}

	return cmd.Start()
}

func jsonJob() []byte {

	return []byte(`{
  "invoice": {
    "currency": "EUR",
    "invoice_id": "RE-0814-min",
    "issue_date": "20260105",
    "service_date": "20260104",
    "buyer": {
      "id": "471115",
      "name1": "Franz Fröhlich AG",
      "name2": "Einkauf Herr Ansprech",
      "name3": "Carnica",
      "street": "Musterstrasse 64",
      "zip": "01221",
      "city": "Musterdorf",
      "country": "Deutschland"
    },
    "info": {
	  "info_date": "20260102"
    },
    "items": [
      {
        "pos": 1,
        "id": "A1",
        "description": "Go Entwicklungsservice",
        "quantity": 10.5,
        "unit": "Stunde",
        "net_price": 95.00,
        "line_total": 997.5,
        "tax_rate": 19
      },
      {
        "pos": 2,
        "id": "B5",
        "description": "Dokumentations-Pauschale",
        "quantity": 1,
        "unit": "Pauschale",
        "net_price": 45.50,
        "line_total": 45.50,
        "tax_rate": 19
      },
      {
        "pos": 3,
        "id": "4711",
        "description": "Go-Grundkurs für Fortgeschrittene",
        "quantity": 1,
        "unit": "Stück",
        "net_price": 150.00,
        "line_total": 150.00,
        "tax_rate": 19
      },
      {
        "pos": 4,
        "id": "0815",
        "description": "LibreOffice Basic Beruhigungstee",
        "quantity": 5,
        "unit": "Liter",
        "net_price": 3.50,
        "line_total": 17.50,
        "tax_rate": 7
      }
    ],
    "vats": [
      {
        "pos": 1,
        "id": "S",
        "description": "USt",
        "vat_base": 1193.00,
        "vat_price": 226.67,
        "tax_rate": 19
      },
      {
        "pos": 2,
        "id": "S",
        "description": "USt",
        "vat_base": 17.50,
        "vat_price": 1.23,
        "tax_rate": 7
      }
    ],
    "totals": {
      "line_total_amount": 1210.50,
      "tax_total_amount": 227.91,
      "grand_total_amount": 1438.41
    }
  }
}`)

}
