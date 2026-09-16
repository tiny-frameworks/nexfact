// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"encoding/json"
	"os"

	"codeberg.org/tiny-frameworks/nexfact/api/job"
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/config"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

const NEXGATE_SELLER = "NEXGATE_MANDANT"

func Parse(jsonData []byte) (job *InputJob, err error) {
	if err := json.Unmarshal(jsonData, &job); err != nil {
		return nil, errors.Wrap(
			errors.InvalidFormat,
			"error unmarshalling json input",
			"orchestrator.input.parse",
			err,
		)
	}
	return job, nil
}

func EnrichSeller(inputjob *InputJob) error {

	// seller and jobQueue are mandatory. For convenience we use Defaults if missing
	if inputjob.Options.Queue == "" {
		inputjob.Options.Queue = job.QueueZugferd
	}

	if inputjob.Mandant.Seller == "" {
		inputjob.Mandant.Seller = config.DefaultSeller
		if envSeller := os.Getenv(NEXGATE_SELLER); envSeller != "" {
			inputjob.Mandant.Seller = envSeller
		}
	}

	return config.LoadSeller(inputjob.Mandant.Seller)
}

func validateJob(iJob *InputJob) error {
	// for testing func validate()
	return ValidateInputInvoice(iJob)
}

func ValidateInputInvoice(inv *InputJob) error {

	if inv.Invoice.InvoiceID == "" {
		return errors.New(
			errors.MissingField,
			"Invoice number is required",
			"orchestrator.input.invoice.number",
		)
	}

	if inv.Invoice.IssueDate == "" {
		return errors.New(
			errors.MissingField,
			"Invoice issueDate is required",
			"orchestrator.input.invoice.issueDate",
		)
	}
	if inv.Invoice.Currency == "" {
		return errors.New(
			errors.MissingField,
			"Invoice currency is required",
			"orchestrator.input.invoice.currency",
		)
	}

	// --- Parties ---
	if inv.Mandant.Seller == "" {
		return errors.New(
			errors.MissingField,
			"Seller name is required",
			"input.invoice.seller.name",
		)
	}
	if inv.Invoice.Buyer.Name1 == "" {
		return errors.New(
			errors.MissingField,
			"Invoice buyer name1 is required",
			"orchestrator.input.invoice.buyer.name1",
		)
	}

	// --- Lines ---
	if len(inv.Invoice.Items) == 0 {
		return errors.New(
			errors.MissingField,
			"at least one invoice line items is required",
			"orchestrator.input.invoice.items",
		)
	}

	for idx, l := range inv.Invoice.Items {
		path := linePath(idx)

		if l.ID == "" {
			return errors.New(
				errors.MissingField,
				"Invoice "+path+".id is required",
				"orchestrator.input.invoice.item.id",
			)
		}
		if l.Description == "" {
			return errors.New(
				errors.MissingField,
				"Invoice lineitem description is required",
				"orchestrator.input.invoice.item.description",
			)
		}
		if l.Quantity <= 0 {
			return errors.New(
				errors.InvalidValue,
				"Invoice "+path+".quantity is required",
				"orchestrator.input.invoice.item.quantity",
			)
		}
		if l.UnitPrice < 0 {
			return errors.New(
				errors.InvalidValue,
				"Invoice "+path+".net_price is required",
				"orchestrator.input.invoice.item.unitprice",
			)
		}
		if l.TaxRate < 0 {
			return errors.New(
				errors.InvalidValue,
				"Invoice "+path+".tax_rateis required",
				"orchestrator.input.invoice.item.taxrate",
			)
		}

		// einfache Konsistenz
		if l.Quantity*l.UnitPrice != l.LineTotal {
			return errors.New(
				errors.Inconsistent,
				"line_total != quantity * net_price",
				"input.line",
			)
		}
	}

	return nil
}

func linePath(i int) string {
	return "lineitem[" + itoa(i) + "]"
}

// minimal, ohne strconv import
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := []byte{}
	for i > 0 {
		buf = append([]byte{byte('0' + i%10)}, buf...)
		i /= 10
	}
	return string(buf)
}
