// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package invoice

import (
	"fmt"
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/input"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

const noBuyerRef string = "NICHT ANGEGEBEN"

func (m *ZUGFeRDmaster) normalizeZF(in *input.InputJob) error {

	// 1. triviales Mapping
	m.mapBasicFields(in)
	if err := m.enrichInvoice(in); err != nil {
		return err
	}

	// 2. Parties auflösen
	if err := m.enrichSeller(in); err != nil {
		return err
	}
	if err := m.enrichBuyer(in); err != nil {
		return err
	}

	// 3. Lines & VAT
	if err := m.mapLines(in); err != nil {
		return err
	}

	if err := m.mapVATs(in); err != nil {
		return err
	}

	// 4. Totals
	if err := m.mapTotals(in); err != nil {
		return err
	}

	// 5. Zahlungsbedingungen / Perioden
	if err := m.mapPayment(in); err != nil {
		return err
	}

	if err := m.mapServicePeriod(in); err != nil {
		return err
	}

	if err := m.mapPaths(in); err != nil {
		return err
	}

	if err := m.mapHeaderFooter(in); err != nil {
		return err
	}
	return nil
}

func (m *ZUGFeRDmaster) mapBasicFields(in *input.InputJob) {
	m.Invoice.InvoiceID = in.Invoice.InvoiceID
	m.Invoice.Profile = ProfileEN16931
	m.Invoice.Currency = m.getCurrency(in.Invoice.Currency)
	m.Invoice.IssueDate = date(in.Invoice.IssueDate)
	m.Options = &Options{
		Queue:          in.Options.Queue,
		PDFSource:      in.Options.PDFSource,
		FacturXSource:  in.Options.FacturXSource,
		TemplateSource: in.Options.TemplateSource,
		Attachments:    in.Options.Attachments,
	}
}

func (m *ZUGFeRDmaster) enrichInvoice(in *input.InputJob) error {

	servicedate := in.Invoice.ServiceDate
	if servicedate == "" {
		servicedate = in.Invoice.ServiceEnd
	}
	m.Invoice.ServiceDate = date(servicedate)

	buyerRef := in.Invoice.BuyerRef
	if buyerRef == "" {
		buyerRef = noBuyerRef
	}
	m.Invoice.BuyerReference = buyerRef

	return nil
}

func (m *ZUGFeRDmaster) enrichBuyer(in *input.InputJob) error {
	m.Invoice.Buyer = Party{
		ID:    in.Invoice.Buyer.ID,
		Name:  in.Invoice.Buyer.Name1,
		Name2: in.Invoice.Buyer.Name2,
		Name3: in.Invoice.Buyer.Name3,
		Address: Address{
			Street:  in.Invoice.Buyer.Street,
			Zip:     in.Invoice.Buyer.Zip,
			City:    in.Invoice.Buyer.City,
			Country: m.getCountry(in.Invoice.Buyer.Country),
		},
		Contact: &Contact{
			Phone: in.Invoice.Info.Phone,
			Fax:   in.Invoice.Info.Fax,
			Email: in.Invoice.Info.Email,
			Name:  in.Invoice.Info.Contact,
		},
	}

	return nil
}

func (m *ZUGFeRDmaster) mapLines(in *input.InputJob) error {

	for _, l := range in.Invoice.Items {

		item := LineItem{
			Pos:         l.Pos,
			ID:          l.ID,
			Description: l.Description,
			Quantity:    l.Quantity,
			Unit:        m.getUnit(l.Unit),
			UnitPrice:   l.UnitPrice,
			LineTotal:   l.LineTotal,
			TaxRate:     l.TaxRate,
		}
		m.Invoice.Lines = append(m.Invoice.Lines, item)
	}

	return nil
}

func (m *ZUGFeRDmaster) mapTotals(in *input.InputJob) error {
	m.Invoice.Totals = Totals{
		LineTotalAmount:  in.Invoice.Totals.LineTotalAmount,
		TaxTotalAmount:   in.Invoice.Totals.TaxTotalAmount,
		GrandTotalAmount: in.Invoice.Totals.GrandTotalAmount,
		DuePayableAmount: in.Invoice.Totals.GrandTotalAmount,
	}
	return nil
}

func (m *ZUGFeRDmaster) mapVATs(in *input.InputJob) error {

	for _, l := range in.Invoice.VATs {

		item := VatItem{
			Category:    l.ID,          // S | Z | E |
			Description: l.Description, // USt
			VatAmount:   l.VATAmount,
			BasisAmount: l.VATBase,
			TaxRate:     l.TaxRate,
		}
		m.Invoice.Vats = append(m.Invoice.Vats, item)
	}
	return nil
}

func (m *ZUGFeRDmaster) mapPayment(in *input.InputJob) error {

	m.Invoice.PaymentTerms = &PaymentTerms{
		DueDate:     date(in.Invoice.PaymentTerms.DueDate),
		Description: in.Invoice.PaymentTerms.DueDescription,
	}
	m.Invoice.PaymentMean = m.getPayment(SepaCT)
	return nil
}

func (m *ZUGFeRDmaster) mapServicePeriod(in *input.InputJob) error {
	m.Invoice.ServicePeriod = &ServicePeriod{
		Start: date(in.Invoice.ServiceStart),
		End:   date(in.Invoice.ServiceEnd),
	}
	return nil
}

func (m *ZUGFeRDmaster) mapPaths(in *input.InputJob) (err error) {
	var (
		outputFile  string
		parFile     string
		xmlFile     string
		ottTemplate string
		//parTemplate string
		//xmlTemplate string
	)

	outputsDir := filepath.Join(config.SellerParams.SellerRoot, "outputs")
	baseName := fmt.Sprintf("%s", m.Invoice.InvoiceID)

	// .odt, .pdf, .par files are output documents. Attributes in LoParameter
	// are either empty or an absolute path or ~ (absolute from env/source/ e.g. env/example)
	if m.Paths.ODTfile == "" { // use default env Odt Dir/name.odt
		odtName := fmt.Sprintf("%s.odt", baseName)
		if outputFile, err = filepath.Abs(filepath.Join(outputsDir, odtName)); err != nil {
			return errors.New(
				errors.NormalizationFailed,
				"toAbsolute()",
				"odtfile",
			)
		}
		m.Paths.ODTfile = outputFile
	} else {
		if !filepath.IsAbs(m.Paths.ODTfile) { // simply pass through
			// nothing to do
		}
	}

	if m.Paths.PDFfile == "" { // use default env Odt Dir/name.pdf
		pdfName := fmt.Sprintf("%s.pdf", baseName)
		if outputFile, err = filepath.Abs(filepath.Join(outputsDir, pdfName)); err != nil {
			return errors.New(
				errors.NormalizationFailed,
				"toAbolute()",
				"pdffile",
			)

		}
		m.Paths.PDFfile = outputFile
	} else {
		if !filepath.IsAbs(m.Paths.PDFfile) { // default output dir
			// nothing to do
		}
	}

	parName := fmt.Sprintf("%s.par", baseName)
	if parFile, err = filepath.Abs(filepath.Join(outputsDir, parName)); err != nil {
		return errors.New(
			errors.NormalizationFailed,
			"toAbolute()",
			"parFile",
		)
	}
	m.Paths.PARfile = parFile

	xmlName := fmt.Sprintf("%s.xml", baseName)
	if xmlFile, err = filepath.Abs(filepath.Join(outputsDir, xmlName)); err != nil {
		return errors.New(
			errors.NormalizationFailed,
			"toAbolute()",
			"xmlFile",
		)
	}
	m.Paths.XMLfile = xmlFile

	// .ott- , .par.tpl or .xml.tpl- file are templates NOT output documents!  therefore
	// they MUST be in env/templates folder; attributes are NOT absolute
	templatesDir := filepath.Join(config.SellerParams.SellerRoot, "templates")
	templateSource := m.Options.TemplateSource
	if m.Paths.OTTfile == "" {
		if templateSource == "" {
			ottTemplate = config.SellerParams.Templates.OttTemplate
		} else {
			ottTemplate = templateSource
			if !filepath.IsAbs(ottTemplate) {
				ottTemplate = filepath.Join(templatesDir, ottTemplate)
			}
		}
		m.Paths.OTTfile = ottTemplate
	}

	// now get parameter .par template
	/*
		if parTemplate, err = filepath.Abs(filepath.Join(templatesDir, config.SellerParams.Templates.ParTemplate)); err != nil {
			return errors.New(
				errors.NormalizationFailed,
				"toAbolute()",
				"parTemplate",
			)
		}
		m.Paths.PARtemplate = parTemplate
	*/
	m.Paths.PARtemplate = config.SellerParams.Templates.ParTemplate
	/*
		if xmlTemplate, err = filepath.Abs(filepath.Join(templatesDir, config.SellerParams.Templates.XmlTemplate)); err != nil {
			return errors.New(
				errors.NormalizationFailed,
				"toAbolute()",
				"xmlTemplate",
			)
		}
		m.Paths.XMLtemplate = xmlTemplate
	*/
	m.Paths.XMLtemplate = config.SellerParams.Templates.XmlTemplate

	// Now normalize attachments
	var attach []string
	if m.Options.Attachments == nil { // slice options is not yet initializes
		attach = make([]string, 0)
	}
	for _, item := range in.Options.Attachments {
		if filepath.IsAbs(item) {
			attach = append(attach, item)
		} else {
			itemPath, _ := filepath.Abs(item)
			attach = append(attach, itemPath)
		}
	}
	m.Options.Attachments = attach
	return nil
}

func (m *ZUGFeRDmaster) mapHeaderFooter(in *input.InputJob) error {
	senderline := config.SellerParams.DIN5008.AbsenderZeile
	if in.Invoice.SenderHeader != "" {
		senderline = in.Invoice.SenderHeader
	}
	m.Headers.SenderLine = senderline

	if len(in.Invoice.PostalNotes) != 0 {
		for i := range len(in.Invoice.PostalNotes) {
			switch i {
			case 0:
				m.Headers.PostalNote1 = clean(in.Invoice.PostalNotes[i])
			case 1:
				m.Headers.PostalNote2 = clean(in.Invoice.PostalNotes[i])
			case 2:
				m.Headers.PostalNote3 = clean(in.Invoice.PostalNotes[i])
			}
		}
	}
	m.Headers.Department = in.Invoice.Info.Department
	m.Headers.InfoDate = date(in.Invoice.Info.InfoDate)

	// Footer vorbelegen
	footAddresses := config.SellerParams.DIN5008.Impressum.AdresseZeile
	if in.Invoice.Impressum.Addresses != "" {
		footAddresses = in.Invoice.Impressum.Addresses
	}
	m.Footers.Addresses = footAddresses

	footContacts := config.SellerParams.DIN5008.Impressum.KontaktZeile
	if in.Invoice.Impressum.Contacts != "" {
		footContacts = in.Invoice.Impressum.Contacts
	}
	m.Footers.Contacts = footContacts

	footBank := config.SellerParams.DIN5008.Impressum.BankZeile
	if in.Invoice.Impressum.Bank != "" {
		footBank = in.Invoice.Impressum.Bank
	}
	m.Footers.Bank = footBank

	footLegal := config.SellerParams.DIN5008.Impressum.LegalZeile
	if in.Invoice.Impressum.Legal != "" {
		footLegal = in.Invoice.Impressum.Legal
	}
	m.Footers.Legal = footLegal

	return nil
}
