// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package config

import "time"

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
