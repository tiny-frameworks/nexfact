// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexutils/errors"
	yaml "github.com/goccy/go-yaml"
)

type SellerConfig struct {
	Version    int               `yaml:"version"`
	Name       string            `yaml:"name"`
	LeitwegID  string            `yaml:"leitweg_id"`
	Address    Address           `yaml:"address"`
	Tax        Tax               `yaml:"tax"`
	Contact    Contact           `yaml:"contact"`
	Bank       Bank              `yaml:"bank"`
	Legal      Legal             `yaml:"legal"`
	DIN5008    DIN5008           `yaml:"din_5008"`
	Paths      SellerPathsConfig `yaml:"paths"`
	Templates  TemplatesConfig   `yaml:"templates"`
	SellerRoot string
}

type Address struct {
	Street  string `yaml:"street"`
	Zip     string `yaml:"zip"`
	City    string `yaml:"city"`
	Country string `yaml:"country"`
}

type Tax struct {
	VatID     string `yaml:"vat_id"`
	TaxNumber string `yaml:"tax_number"`
}

type Contact struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	Phone string `yaml:"phone"`
}

type Bank struct {
	Name string `yaml:"name"`
	IBAN string `yaml:"iban"`
	BIC  string `yaml:"bic"`
}

type Legal struct {
	Form             string `yaml:"form"`
	HRB              string `yaml:"hrb"`
	ManagingDirector string `yaml:"managing_director"`
}

type DIN5008 struct {
	AbsenderZeile string    `yaml:"absender_zeile"`
	Impressum     Impressum `yaml:"impressum"`
}

type Impressum struct {
	AdresseZeile string `yaml:"adresse_zeile"`
	KontaktZeile string `yaml:"kontakt_zeile"`
	BankZeile    string `yaml:"bank_zeile"`
	LegalZeile   string `yaml:"legal_zeile"`
}

type SellerPathsConfig struct {
	LogFile string `yaml:"log_file"`
}

type TemplatesConfig struct {
	ParTemplate string `yaml:"par"`
	XmlTemplate string `yaml:"xml"`
	OttTemplate string `yaml:"ott"`
}

const DefaultSeller = "default"

var SellerParams *SellerConfig

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
