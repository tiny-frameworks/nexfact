// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package invoice

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
)

type CacheItem struct {
	ID    string // stunden
	Code  string // XRechnung (HUR)
	Label string // PDF (Std)
}

const (
	CURRENCIES string = "currencies"
	COUNTRIES  string = "countries"
	UNITS      string = "units"
	PAYMENTS   string = "payments"
)

/*
	func newCache(capacity int, ttl time.Duration, cleanup time.Duration) *cache.LRUCache {
		// Cache mit Kapazität 20, TTL 8h, Cleanup alle 2h
		c := cache.New(capacity, ttl, cleanup)
		return c
	}
*/
func (c *CacheItem) GetLabel() string {
	return c.Label
}

func (c *CacheItem) GetCode() string {
	return c.Code
}

// Loader für cache
func (m *ZUGFeRDmaster) CsvLoader(searchKey string, source string) (interface{}, error) {

	var csvFile string

	if filepath.IsAbs(source) { // to simplifiy tests
		csvFile = source
	} else {
		switch source {
		case CURRENCIES:
			csvFile = config.SystemParams.Paths.CurrenciesCSV
		case COUNTRIES:
			csvFile = config.SystemParams.Paths.CountriesCSV
		case UNITS:
			csvFile = config.SystemParams.Paths.UnitsCSV
		case PAYMENTS:
			csvFile = config.SystemParams.Paths.PaymentsCSV
		}
	}

	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	key := strings.ToLower(searchKey)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if record[0] == key {
			return &CacheItem{
				ID:    record[0],
				Code:  record[1],
				Label: record[2],
			}, nil
		}
	}
	return nil, nil
}

func (m *ZUGFeRDmaster) getCountry(key string) *CacheItem {
	if key == "" {
		return nil
	}
	found, err := config.Caches["country"].GetOrLoad(key, func() (interface{}, error) {
		return m.CsvLoader(key, COUNTRIES)
	})
	if found == nil || err != nil {
		return nil
	}

	return found.(*CacheItem)
}

func (m *ZUGFeRDmaster) getCurrency(key string) *CacheItem {
	if key == "" {
		return nil
	}
	found, err := config.Caches["currency"].GetOrLoad(key, func() (interface{}, error) {
		return m.CsvLoader(key, CURRENCIES)
	})
	if found == nil || err != nil {
		return nil
	}
	return found.(*CacheItem)
}

func (m *ZUGFeRDmaster) getUnit(key string) *CacheItem {
	if key == "" {
		return nil
	}
	found, err := config.Caches["unit"].GetOrLoad(key, func() (interface{}, error) {
		return m.CsvLoader(key, UNITS)
	})
	if found == nil || err != nil {
		return nil
	}
	return found.(*CacheItem)
}

func (m *ZUGFeRDmaster) getPayment(key string) *CacheItem {
	if key == "" {
		return nil
	}
	found, err := config.Caches["payment"].GetOrLoad(key, func() (interface{}, error) {
		return m.CsvLoader(key, PAYMENTS)
	})
	if found == nil || err != nil {
		return nil
	}
	return found.(*CacheItem)
}
