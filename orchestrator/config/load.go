// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/writer"
	"codeberg.org/tiny-frameworks/nexutils/errors"
	yaml "github.com/goccy/go-yaml"
)

var (
	SystemParams *SystemConfig
	SellerParams *SellerConfig
)

func LoadSystem(envRoot string) error {
	sRoot := filepath.Join(envRoot, "system")
	yamlPath := filepath.Join(sRoot, "system.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return errors.Wrap(
			errors.ReadError,
			fmt.Sprintf("Can not read %s", yamlPath),
			"orchestrator.config.LoadSystem",
			err,
		)
	}

	var cfg SystemConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return errors.Wrap(
			errors.InternalError,
			fmt.Sprintf("Can not unmarshal data %v", data),
			"orchestrator.config.LoadSystem",
			err,
		)
	}

	cfg.EnvRoot = envRoot
	cfg.SystemRoot = sRoot
	SystemParams = &cfg

	if err = normalizeSystem(); err != nil {
		return err
	}

	loadCaches()

	return nil
}

func normalizeSystem() error {
	// Small helper function to avoid code duplication
	absPath := func(base string, elem ...string) (string, error) {
		p, err := filepath.Abs(filepath.Join(base, filepath.Join(elem...)))
		if err != nil {
			return "", errors.Wrap(
				errors.ReadError,
				fmt.Sprintf("can't build absolute path for components %v", elem),
				"orchestrator.config.load.normalizeSystem",
				err,
			)
		}
		return p, nil
	}

	var err error
	p := SystemParams.SystemRoot

	if SystemParams.Paths.CountriesCSV, err = absPath(p, "csv", SystemParams.Paths.CountriesCSV); err != nil {
		return err
	}
	if SystemParams.Paths.UnitsCSV, err = absPath(p, "csv", SystemParams.Paths.UnitsCSV); err != nil {
		return err
	}
	if SystemParams.Paths.CurrenciesCSV, err = absPath(p, "csv", SystemParams.Paths.CurrenciesCSV); err != nil {
		return err
	}
	if SystemParams.Paths.PaymentsCSV, err = absPath(p, "csv", SystemParams.Paths.PaymentsCSV); err != nil {
		return err
	}
	if SystemParams.Engine.PdfEngine, err = absPath(p, SystemParams.Engine.PdfEngine); err != nil {
		return err
	}
	if SystemParams.Engine.ZugferdEngine, err = absPath(p, SystemParams.Engine.ZugferdEngine); err != nil {
		return err
	}

	return nil
}

func checkContainer() error {
	heartbeatFrequency := SystemParams.Engine.HeartbeatFrequency * 2
	providerRoot := filepath.Join(SellerParams.SellerRoot, "data")

	if SystemParams.Engine.PdfProvider == ProviderContainer {
		hbFile := filepath.Join(providerRoot, "pdfs", "logs", "nexgate.heartbeat")
		if _, err := writer.CheckHeartbeat(hbFile, heartbeatFrequency); err != nil {
			return err
		}
	}
	if SystemParams.Engine.ZugferdProvider == ProviderContainer {
		hbFile := filepath.Join(providerRoot, "zugferds", "logs", "nexgate.heartbeat")
		if _, err := writer.CheckHeartbeat(hbFile, heartbeatFrequency); err != nil {
			return err
		}
	}
	return nil
}

func checkTimeout(timeout, lower, upper time.Duration) error {
	if timeout < lower {
		return errors.New(
			errors.InvalidValue,
			fmt.Sprintf("timeout too small: %v", timeout),
			"orchestrator.config.load.timeout",
		)
	}
	if timeout > upper {
		return errors.New(
			errors.InvalidValue,
			fmt.Sprintf("timeout too big: %v", timeout),
			"orchestrator.config.load.timeout",
		)
	}
	return nil
}

func LoadSeller(sellerID string) error {
	sRoot := filepath.Join(SystemParams.EnvRoot, "sellers", sellerID)
	seller := &SellerConfig{
		SellerRoot: sRoot,
	}

	cfg, err := Caches["seller"].GetOrLoad(sellerID, seller.loader)
	if err != nil {
		return err
	}

	SellerParams = cfg.(*SellerConfig)

	if err := normalizeSeller(); err != nil {
		return err
	}
	if wErr := writeEnv(); wErr != nil {
		return wErr
	}

	if err := checkContainer(); err != nil {
		return err
	}

	return nil
}

func (seller *SellerConfig) loader() (interface{}, error) {
	yamlPath := filepath.Join(seller.SellerRoot, "seller.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	var cfg SellerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.SellerRoot = seller.SellerRoot
	return &cfg, nil
}

func normalizeSeller() error {
	absPath := func(base string, elem ...string) (string, error) {
		p, err := filepath.Abs(filepath.Join(base, filepath.Join(elem...)))
		if err != nil {
			return "", errors.Wrap(
				errors.ReadError,
				fmt.Sprintf("can't build absolute path for components %v", elem),
				"orchestrator.config.load.normalizeSeller",
				err,
			)
		}
		return p, nil
	}

	var err error
	p := SellerParams.SellerRoot

	if SellerParams.Paths.LogFile, err = absPath(p, "logs", SellerParams.Paths.LogFile); err != nil {
		return err
	}
	if SellerParams.Templates.OttTemplate, err = absPath(p, "templates", SellerParams.Templates.OttTemplate); err != nil {
		return err
	}
	if SellerParams.Templates.ParTemplate, err = absPath(p, "templates", SellerParams.Templates.ParTemplate); err != nil {
		return err
	}
	if SellerParams.Templates.XmlTemplate, err = absPath(p, "templates", SellerParams.Templates.XmlTemplate); err != nil {
		return err
	}

	return nil
}

func writeEnv() error {
	sellerEnv := map[string]time.Duration{
		"CMD_TIMEOUT":         SystemParams.Engine.CmdTimeout,
		"HEARTBEAT_FREQUENCY": SystemParams.Engine.HeartbeatFrequency,
		"WATCH_SLEEP":         SystemParams.Engine.WatchSleep,
		"CONFIG_RELAOAD":      SystemParams.Engine.ConfigReloadFrequency,
	}

	sellerDataPath := filepath.Join(SellerParams.SellerRoot, "data")
	sellerPdfEnvPath := filepath.Join(sellerDataPath, "pdfs", "logs", "nexgate.env")
	sellerZugferdEnvPath := filepath.Join(sellerDataPath, "zugferds", "logs", "nexgate.env")

	pFile, pErr := os.Create(sellerPdfEnvPath)
	if pErr != nil {
		return errors.Wrap(errors.WriteError, fmt.Sprintf("Error creating file: %s", sellerPdfEnvPath), "orchestrator.config.load.writeEnv", pErr)
	}
	defer pFile.Close()

	zFile, zErr := os.Create(sellerZugferdEnvPath)
	if zErr != nil {
		return errors.Wrap(errors.WriteError, fmt.Sprintf("Error creating file: %s", sellerZugferdEnvPath), "orchestrator.config.load.writeEnv", zErr)
	}
	defer zFile.Close()

	// An elegant Go approach: Wrap both files in a MultiWriter. A single write operation writes to both.
	multiWriter := io.MultiWriter(pFile, zFile)

	for key, val := range sellerEnv {
		secondString := strconv.FormatFloat(val.Seconds(), 'f', -1, 64)
		envString := fmt.Sprintf("%s=%s\n", key, secondString)
		if _, err := multiWriter.Write([]byte(envString)); err != nil {
			return errors.Wrap(errors.WriteError, "failed to write env configuration", "orchestrator.config.load.writeEnv", err)
		}
	}

	return nil
}
