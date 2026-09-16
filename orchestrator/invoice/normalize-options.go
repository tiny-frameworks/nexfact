// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package invoice

import (
	"fmt"
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexfact/orchestrator/input"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

func (m *ZUGFeRDmaster) normalizeOptions(in *input.InputJob) error {

	// Seller anreichern
	if err := m.enrichSeller(in); err != nil {
		return err
	}

	m.Invoice.InvoiceID = in.Invoice.InvoiceID
	m.Invoice.Profile = ProfileEN16931

	m.Options.Queue = in.Options.Queue
	m.Options.Attachments = in.Options.Attachments
	if in.Options.PDFSource != "" {
		if pdfSource, err := filepath.Abs(in.Options.PDFSource); err != nil {
			return errors.Wrap(
				errors.ReadError,
				fmt.Sprintf("Pdf Path %s is no absolute path", in.Options.PDFSource),
				"orchestrator.invoice.normalizeOptions",
				err,
			)
		} else {
			m.Options.PDFSource = pdfSource
		}
	}

	if in.Options.FacturXSource != "" {
		if fxSource, err := filepath.Abs(in.Options.FacturXSource); err != nil {
			return errors.Wrap(
				errors.ReadError,
				fmt.Sprintf("Factur-X Path %s is no absolute path", in.Options.FacturXSource),
				"orchestrator.invoice.normalizeOptions",
				err,
			)
		} else {
			m.Options.FacturXSource = fxSource
		}
	}
	if in.Options.TemplateSource != "" {
		if template, err := filepath.Abs(in.Options.TemplateSource); err != nil {
			return errors.Wrap(
				errors.ReadError,
				fmt.Sprintf("Template Path %s is no absolute path", in.Options.TemplateSource),
				"orchestrator.invoice.normalizeOptions",
				err,
			)
		} else {
			m.Options.TemplateSource = template
		}
	}

	return nil
}
