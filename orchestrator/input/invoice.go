// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package input

type InputJob struct {
	Mandant MandantType `json:"mandant"`
	Options OptionsType `json:"options,omitempty"`
	Invoice InvoiceType `json:"invoice"`
}

type MandantType struct {
	MandantID string `json:"mandant_id"`
	Seller    string `json:"seller"`
}

type OptionsType struct {
	Queue          string   `json:"queue"`
	PDFSource      string   `json:"pdf"`
	FacturXSource  string   `json:"factur-x"`
	TemplateSource string   `json:"template"`
	Attachments    []string `json:"attachments"`
}

type InvoiceType struct {
	Currency     string `json:"currency"`
	InvoiceID    string `json:"invoice_id"`
	IssueDate    string `json:"issue_date"`
	ServiceDate  string `json:"service_date"`
	ServiceStart string `json:"service_start"`
	ServiceEnd   string `json:"service_end"`
	BuyerRef     string `json:"buyer_reference"`

	Buyer Buyer `json:"buyer"`

	SenderHeader string   `json:"sender_header"`
	PostalNotes  []string `json:"postal_notes"`

	Info         Info         `json:"info"`
	Impressum    Impressum    `json:"impressum"`
	PaymentTerms PaymentTerms `json:"payment_terms"`

	Items  []Item `json:"items"`
	VATs   []VAT  `json:"vats"`
	Totals Totals `json:"totals"`
}

type Buyer struct {
	ID      string `json:"id"`
	Name1   string `json:"name1"`
	Name2   string `json:"name2"`
	Name3   string `json:"name3"`
	Street  string `json:"street"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Country string `json:"country"`
	Email   string `json:"email"`
}

type Info struct {
	Department string `json:"department"`
	Contact    string `json:"contact"`
	Phone      string `json:"phone"`
	Fax        string `json:"fax"`
	Email      string `json:"email"`
	InfoDate   string `json:"info_date"`
}

type Impressum struct {
	Addresses string `json:"addresses"`
	Contacts  string `json:"contacts"`
	Bank      string `json:"bank"`
	Legal     string `json:"legal"`
}

type PaymentTerms struct {
	DueDescription string `json:"due_description"`
	DueDate        string `json:"due_date"`
}

type Item struct {
	Pos         int     `json:"pos"`
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`

	UnitPrice float64 `json:"net_price"` // Price pro Unit

	LineTotal float64 `json:"line_total"`
	TaxRate   float64 `json:"tax_rate"`
}

type VAT struct {
	Pos         int     `json:"pos"`
	ID          string  `json:"id"`
	Description string  `json:"description"`
	VATBase     float64 `json:"vat_base"`
	VATAmount   float64 `json:"vat_price"`
	TaxRate     float64 `json:"tax_rate"`
}

type Totals struct {
	LineTotalAmount  float64 `json:"line_total_amount"`
	TaxTotalAmount   float64 `json:"tax_total_amount"`
	GrandTotalAmount float64 `json:"grand_total_amount"`
}
