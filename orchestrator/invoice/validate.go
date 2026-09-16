// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package invoice

import (
	"math"

	"codeberg.org/tiny-frameworks/nexutils/errors"
)

func (m *ZUGFeRDmaster) validateZF() error {

	if m.Invoice.InvoiceID != "" {
		if err := validateInvoice(m.Invoice); err != nil {
			return err
		}
	}
	if err := validateOptions(m.Options); err != nil {
		return err
	}
	return nil
}

func validateOptions(opts *Options) error {
	return nil
}

func validateInvoice(inv *Invoice) error {

	if inv.Seller.Name == "" {
		return errors.New(
			errors.IncompleteParty, "seller name required", "core.invoice.seller.name",
		)
	}

	if inv.InvoiceID == "" {
		return errors.New(
			errors.MissingField, "Invoice number is required", "core.invoice.id",
		)
	}
	if inv.ServiceDate.IsZero() && inv.ServicePeriod.End.IsZero() {
		return errors.New(
			errors.MissingField, "Servicedate or ServicePeriod is required", "core.invoice.servicedate",
		)
	}
	if len(inv.Lines) == 0 {
		return errors.New(
			errors.NormalizationFailed, "No line items", "core.invoice.Lines",
		)
	}

	if inv.Totals.GrandTotalAmount <= 0 {
		return errors.New(
			errors.TotalMismatch, "Invalid Totals", "core.invoice.Totals",
		)
	}

	// Einfacher Check: Ist Netto + Steuer = Brutto?
	expectedGross := inv.Totals.LineTotalAmount + inv.Totals.TaxTotalAmount
	if math.Abs(expectedGross-inv.Totals.GrandTotalAmount) > 0.01 {
		return errors.New(
			errors.TotalMismatch,
			"calculated Totals - inv Totals > 0.01",
			"core.invoice.Totals",
		)
	}

	// ValidateTaxBreakdown prüft die Regel BR-CO-17
	for _, tax := range inv.Vats {
		// Die Formel des Validators: Basis * (Satz / 100)
		expectedAmount := round2(tax.BasisAmount * (tax.TaxRate / 100))

		if math.Abs(tax.VatAmount-expectedAmount) > 0.001 {
			return errors.New(
				errors.RuleViolation,
				"Steuer Rechenfehler",
				"tax.TaxRate",
			)
		}
	}

	return nil
}
