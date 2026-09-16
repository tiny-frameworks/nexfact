// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package invoice

import (
	"time"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/input"
)

type ZUGFeRDmaster struct {
	Mandant *Mandant
	Options *Options
	Invoice *Invoice
	Paths   *Paths
	Headers *Headers
	Footers *Footers
}

type Mandant struct {
	MandantID string
	Seller    string
}

type Options struct {
	Queue          string
	PDFSource      string
	TemplateSource string
	FacturXSource  string
	Attachments    []string
}

type Invoice struct {
	// Identität & Kontext
	InvoiceID string
	Profile   Profile    // EN16931 (enum)
	Currency  *CacheItem // ISO 4217
	IssueDate time.Time

	// Parteien
	Seller Party
	Buyer  Party

	// Referenzen & Perioden
	BuyerReference string
	ServicePeriod  *ServicePeriod
	ServiceDate    time.Time

	// Inhalte
	Lines  []LineItem
	Vats   []VatItem
	Totals Totals

	// Zahlungsbedingungen
	PaymentTerms *PaymentTerms
	PaymentMean  *CacheItem
}

type Party struct {
	ID      string // interne Referenz
	Name    string
	Name2   string
	Name3   string
	Address Address
	Contact *Contact
	IBAN    string
	TaxID   string // VAT-ID
}

type Address struct {
	Street  string
	Zip     string
	City    string
	Country *CacheItem // ISO 3166
}

type LineItem struct {
	Pos         int
	ID          string
	Description string

	Quantity  float64
	Unit      *CacheItem // normiert (UNECE)
	UnitPrice float64

	LineTotal float64
	TaxRate   float64
}

type VatItem struct {
	TaxRate     float64
	Category    string // S, Z, AE, ...
	Description string
	BasisAmount float64
	VatAmount   float64
}

type Totals struct {
	LineTotalAmount      float64
	TaxExclusiveAmount   float64
	ChargeTotalAmount    float64
	AllowanceTotalAmount float64
	DuePayableAmount     float64
	TaxTotalAmount       float64
	GrandTotalAmount     float64
}

type ServicePeriod struct {
	Start time.Time
	End   time.Time
}

type PaymentTerms struct {
	DueDate     time.Time
	Description string
}

type Contact struct {
	Name  string
	Phone string
	Email string
	Fax   string
}

type Paths struct {
	OTTfile     string
	ODTfile     string
	PDFfile     string
	PARfile     string
	XMLfile     string
	PARtemplate string
	XMLtemplate string
}

type Headers struct {
	SenderLine  string
	PostalNote1 string
	PostalNote2 string
	PostalNote3 string
	Department  string
	InfoDate    time.Time
}

type Footers struct {
	Addresses string
	Contacts  string
	Bank      string
	Legal     string
}

type Profile string

const (
	ProfileEN16931 Profile = "EN16931"
	SepaCT                 = "58"
)

func New() *ZUGFeRDmaster {

	m := &ZUGFeRDmaster{
		Mandant: &Mandant{},
		Options: &Options{},
		Invoice: &Invoice{},
		Paths:   &Paths{},
		Headers: &Headers{},
		Footers: &Footers{},
	}
	return m
}

func (m *ZUGFeRDmaster) NormalizeAndValidate(in *input.InputJob) error {

	switch in.Options.Queue {
	case job.QueueCombine, job.QueueExtract, job.QueueValidate: // only seller and option is necessary
		if err := m.normalizeOptions(in); err != nil { // Transform json inputJob to ZugferdMaster
			return err
		}

	default:
		// Transform json inputJob to ZugferdMaster
		if err := m.normalizeZF(in); err != nil { // Transform json inputJob to ZugferdMaster
			return err
		}

		if err := m.validateZF(); err != nil { // Validate ZugferdMaster
			return err
		}
	}
	return nil
}
