// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package invoice

import (
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/input"
)

func (m *ZUGFeRDmaster) enrichSeller(in *input.InputJob) error {

	m.Invoice.Seller = Party{
		ID:   in.Mandant.Seller,
		Name: config.SellerParams.Name,
		Address: Address{
			Street:  config.SellerParams.Address.Street,
			Zip:     config.SellerParams.Address.Zip,
			City:    config.SellerParams.Address.City,
			Country: m.getCountry(config.SellerParams.Address.Country),
		},
		TaxID: config.SellerParams.Tax.VatID,
		IBAN:  config.SellerParams.Bank.IBAN,
		Contact: &Contact{
			Name:  config.SellerParams.Contact.Name,
			Email: config.SellerParams.Contact.Email,
			Phone: config.SellerParams.Contact.Phone,
		},
	}
	m.Mandant.MandantID = in.Mandant.MandantID
	m.Mandant.Seller = in.Mandant.Seller

	return nil
}
