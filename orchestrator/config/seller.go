// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package config

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
