// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/tiny-frameworks/nexutils/errors"
	yaml "github.com/goccy/go-yaml"
)

const (
	ProviderNative    = "native"
	ProviderContainer = "container"
)

type SystemConfig struct {
	Version    int               `yaml:"version"`
	EnvRoot    string            `yaml:"-"`
	Engine     EngineConfig      `yaml:"engine"`
	Paths      SystemPathsConfig `yaml:"paths"`
	Caches     CacheConfig       `yaml:"caches"`
	SystemRoot string
}

type EngineConfig struct {
	HeartbeatFrequency    time.Duration `yaml:"heartbeat_frequency"`
	WatchSleep            time.Duration `yaml:"watch_sleep"`
	CmdTimeout            time.Duration `yaml:"cmd_timeout"`
	ConfigReloadFrequency time.Duration `yaml:"config_reload"`
	DefaultQueue          string        `yaml:"default_queue"`
	ContainerEngine       string        `yaml:"container_engine"`
	PdfEngine             string        `yaml:"pdf_engine"`
	ZugferdEngine         string        `yaml:"zugferd_engine"`
	PdfProvider           string        `yaml:"pdf_provider"`
	ZugferdProvider       string        `yaml:"zugferd_provider"`
}

type SystemPathsConfig struct {
	CountriesCSV  string `yaml:"countries_csv"`
	UnitsCSV      string `yaml:"units_csv"`
	CurrenciesCSV string `yaml:"currencies_csv"`
	PaymentsCSV   string `yaml:"payments_csv"`
}

type CacheConfig struct {
	Seller   SingleCache `yaml:"seller_cache"`
	Country  SingleCache `yaml:"country_cache"`
	Currency SingleCache `yaml:"currency_cache"`
	Unit     SingleCache `yaml:"unit_cache"`
	Payment  SingleCache `yaml:"payment_cache"`
}

type SingleCache struct {
	TTL      time.Duration `yaml:"ttl"`
	Cleanup  time.Duration `yaml:"cleanup"`
	Capacity int           `yaml:"capacity"`
}

var SystemParams *SystemConfig

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
