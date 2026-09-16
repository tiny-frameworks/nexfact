// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package orchestrator_test

import (
	"os"
	"path/filepath"
	"testing"

	writer "codeberg.org/tiny-frameworks/nexfact/api/writer"

	// These lines activate the engines in the background:

	_ "codeberg.org/tiny-frameworks/nexfact/engines/base"
	_ "codeberg.org/tiny-frameworks/nexfact/engines/pdf"
	_ "codeberg.org/tiny-frameworks/nexfact/engines/zugferd"
)

const DEBUG = false

const testSeller = "testSeller"
const ( // provider consts
	CONTAINER = "container"
	NATIVE    = "native"
)

const ( // scenario consts
	FACTURX_SCENE  = "Run_FacturX"
	PAR_SCENE      = "Run_Par"
	PDF_SCENE      = "Run_PDF"
	RUN_MIN        = "Run_Min"
	COMBINE_OK     = "combine_ok"
	VALIDATE_OK    = "validate_ok"
	VALIDATE_ERROR = "validate_error"
	EXTRACT_OK     = "extract_ok"
)

// test-szenarien, which should run for each provider
var scenarios = []struct {
	name     string
	jsonFile string
	expFile  string
}{
	{PDF_SCENE, "test_generate_pdf.json", "RE-0815-max-ohne.pdf"},
	{FACTURX_SCENE, "test_generate_fx.json", "RE-0815-max-fx.xml"},
	{PAR_SCENE, "test_generate_par.json", "RE-0815-max-par.par"},
	{RUN_MIN, "test_generate_min_invoice.json", "RE-0814-min.zf.pdf"},
	{COMBINE_OK, "test_combine_invoice.json", "RE-0815-combine-mit.zf.pdf"},
	{VALIDATE_OK, "test_validate_ok.json", "test_valid_zf.result.pdf"},
	{VALIDATE_ERROR, "test_validate_error.json", "test_error_fx.result.pdf"},
	{EXTRACT_OK, "test_extract_zf.json", "factur-x.xml"},
}

/*
 * use of 'const DEBUG = true'
 * table-driven-test can freeze in debug modus. Therefore
 * we use the TestOrchestratorSuite as 'normal' integration tests
 * to do this, set DEBUG = false
 *
 * if you must debug a test, set DEBUG = false and us the testcase in the
 * second section of this test-file
 */

// Helper Func: setupTestEnv return paths
func setupTestEnv(t *testing.T, provider string) (string, string) {
	envRoot, _ := filepath.Abs(filepath.Join("testdata", "env"))
	resourcePath, _ := filepath.Abs(filepath.Join("testdata", "resources"))

	systemNativeResource := filepath.Join(resourcePath, "system.native")
	systemContainerResource := filepath.Join(resourcePath, "system.container")
	systemYamlTarget := filepath.Join(envRoot, "system", "system.yaml")

	// clean data/artefacts, data/base/in, data/pdfs/in, data/zugferds/in and outputs
	dataRoot := filepath.Join(envRoot, "sellers", testSeller, "data")

	artefactsRoot := filepath.Join(dataRoot, "artefacts")
	removeFiles(t, artefactsRoot)

	baseIn := filepath.Join(dataRoot, "base/in")
	removeFiles(t, baseIn)

	pdfsIn := filepath.Join(dataRoot, "pdfs/in")
	removeFiles(t, pdfsIn)

	zfsIn := filepath.Join(dataRoot, "zugferds/in")
	removeFiles(t, zfsIn)

	outputs := filepath.Join(envRoot, "sellers", testSeller, "outputs")
	removeFiles(t, outputs)

	// copy config based to provider
	switch provider {
	case NATIVE:
		writer.CopyFile(systemNativeResource, systemYamlTarget)
	case CONTAINER:
		writer.CopyFile(systemContainerResource, systemYamlTarget)
	}

	return envRoot, resourcePath
}

func removeFiles(t *testing.T, dir string) error {

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if err = os.RemoveAll(filePath); err != nil {
			t.Log(err)
		}
	}

	return nil
}
