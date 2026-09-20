// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

// //go:build !native && !webservice

//go:build suite

package orchestrator_test

import (
	"os"
	"path/filepath"
	"testing"

	api "codeberg.org/tiny-frameworks/nexfact/orchestrator"

	// These lines activate the engines in the background:
	_ "codeberg.org/tiny-frameworks/nexfact/engines/base"
	_ "codeberg.org/tiny-frameworks/nexfact/engines/pdf"
	_ "codeberg.org/tiny-frameworks/nexfact/engines/zugferd"
)

/*
 * use of 'const DEBUG = true'
 * table-driven-test can freeze in debug modus. Therefore
 * we use the TestOrchestratorSuite as 'normal' integration tests
 * to do this, set DEBUG = false
 *
 * if you must debug a test, set DEBUG = false and us the testcase in the
 * second section of this test-file
 */

// Begin section: table-driven-test: the central entry-point
func TestOrchestratorSuite(t *testing.T) {

	// Protection against "forgetting" (CI/Pipeline)
	if DEBUG && os.Getenv("CI") == "" {
		t.Log("DEBUG Mode active or in CI Pipeline 'const DEBUG = true' :=> no table driven tests!")
		t.Skip("in DEBUG Mode or in CI Pipeline :=> skip TestOrchestratorSuite")
		return
	}

	// === providers we test ==
	// At the moment we test both: native and container
	// If you wont test only container replace: providers := []string{NATIVE, CONTAINER} with: providers := []string{CONTAINER}
	// see 'orchestrator_suite_test.go - corresponding func TestOrchestratorSuite(t *testing.T)'
	// === providers we test ==
	providers := []string{NATIVE, CONTAINER}

	for _, provider := range providers {
		t.Run("Provider_"+provider, func(t *testing.T) {
			// setup for this provider
			envRoot, resourcePath := setupTestEnv(t, provider)

			// now run all scenaries for this provider
			for _, scenario := range scenarios {
				t.Run(scenario.name, func(t *testing.T) {
					// Now we run direct the generic Test-Logic
					runGenericOrchestratorTest(t, envRoot, resourcePath, scenario.jsonFile, scenario.expFile)
				})
			}
		})
	}
}

func runGenericOrchestratorTest(t *testing.T, envRoot, resourcePath, jsonFile, expectedOutputName string) {
	orch, err := api.New(envRoot)
	if err != nil {
		t.Fatalf("failed to create api: %v", err)
	}

	file := filepath.Join(resourcePath, jsonFile)
	jsonData, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read json: %v", err)
	}

	expFile := filepath.Join(envRoot, "sellers", testSeller, "outputs", expectedOutputName)
	result, err := orch.RunJsonJob(jsonData)
	if err != nil {
		t.Errorf("job execution failed: %v", err)
		return
	}

	if expFile != result.OutputPath {
		t.Errorf("expected: %s\nfound: %s", expFile, result.OutputPath)
	}
}

// End section: table-driven-test

/**********************************************
* Begin section: single go test for debugging *
**********************************************/

const (
	// ** control of debug-test [ CONTAINER | NATIVE ]
	DEBUG_PROVIDER = NATIVE
	//DEBUG_PROVIDER = CONTAINER

	// ** available scenarions
	DEBUG_SCENARIO = PAR_SCENE
	//DEBUG_SCENARIO = FACTURX_SCENE
	//DEBUG_SCENARIO = PDF_SCENE
	//DEBUG_SCENARIO = RUN_MIN
	//DEBUG_SCENARIO = COMBINE_OK
	//DEBUG_SCENARIO = VALIDATE_OK
	//DEBUG_SCENARIO = VALIDATE_ERROR
	//DEBUG_SCENARIO = EXTRACT_OK
)

func TestOrchestrator_Run_in_Debug(t *testing.T) {

	// Protection against "forgetting" (CI/Pipeline)
	if !DEBUG || os.Getenv("CI") != "" {
		t.Log("not in DEBUG Mode 'const DEBUG = false' :=> no normal table driven tests!")
		t.Skip("not in DEBUG-Mode or in CI Pipeline :=> skip single-test")
		return
	}

	for _, scenario := range scenarios {
		if scenario.name == DEBUG_SCENARIO {
			envRoot, resourcePath := setupTestEnv(t, DEBUG_PROVIDER)
			runGenericOrchestratorTest(t, envRoot, resourcePath, scenario.jsonFile, scenario.expFile)
			return // test successfully startet and ended
		}
	}
	t.Fatalf("Scenario %s was not found in table!", DEBUG_SCENARIO)

}

// End section: single go test for debugging
