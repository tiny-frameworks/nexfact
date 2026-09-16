// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/orchestrator/invoice"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

type ParInvoice struct {
	Invoice   *invoice.ZUGFeRDmaster
	LineItems []string
}

type ParGenerator struct {
	tpl *template.Template
}

func (g *ParGenerator) Profile() invoice.Profile {
	return "EN16931"
}

func NewPar(templatePath string) (Generator, error) {
	baseName := filepath.Base(templatePath)
	tpl, err := template.New(baseName).
		Funcs(*parFuncMap()).
		Option("missingkey=error").
		ParseFiles(templatePath)
	if err != nil {
		return nil, errors.Wrap(
			errors.InvalidValue,
			fmt.Sprintf("basename [%s] or templatepath [%s] not found", baseName, templatePath),
			"par.simple.generator.New",
			err,
		)
	}
	return &ParGenerator{tpl: tpl}, nil

}

func (g *ParGenerator) Generate(inv *invoice.ZUGFeRDmaster) ([]byte, error) {
	var buf bytes.Buffer

	parInv := g.setupLineItems(inv)

	if err := g.tpl.Execute(&buf, parInv); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *ParGenerator) setupLineItems(inv *invoice.ZUGFeRDmaster) *ParInvoice {
	var line string

	sort.SliceStable(inv.Invoice.Lines, func(i, j int) bool {
		return inv.Invoice.Lines[i].Pos < inv.Invoice.Lines[j].Pos
	})

	// adding lineItems
	items := []string{}
	for idx, li := range inv.Invoice.Lines {
		unitLabel := li.Unit.Label
		line = fmt.Sprintf("t-leistung[li%02d]=%02d;%s;%s;%.2f;%.2f;%.2f;%.2f",
			idx, li.Pos, li.Description, unitLabel, li.Quantity,
			li.UnitPrice, li.LineTotal, li.TaxRate)
		items = append(items, line)
	}

	// adding lineNetSum
	//5;Summe; netto;;; 1210.50;
	line = fmt.Sprintf("t-leistung[ls01]= ;Summe;netto; ; ;%.2f;", inv.Invoice.Totals.LineTotalAmount)
	items = append(items, line)

	// adding vatItems
	// t-leistung[vi01] = 6;;USt;;17.50;1.23;7
	for idx, vi := range inv.Invoice.Vats {
		line = fmt.Sprintf("t-leistung[vi%02d]= ; ;%s; ;%.2f;%.2f;%.2f",
			idx, vi.Description, vi.BasisAmount, vi.VatAmount, vi.TaxRate)
		items = append(items, line)
	}

	// adding vatSum
	//line = fmt.Sprintf("t-leistung[vs01]=%02d;Summe;USt;;;%.2f;", tIdx, inv.TotalVat)
	line = fmt.Sprintf("t-leistung[vs01]= ;Summe;USt; ; ;%.2f;", inv.Invoice.Totals.TaxTotalAmount)
	items = append(items, line)

	line = fmt.Sprintf("t-leistung[zb01]= ;Gesamtbetrag;EUR; ; ;%.2f;", inv.Invoice.Totals.GrandTotalAmount)
	items = append(items, line)

	return &ParInvoice{
		Invoice:   inv,
		LineItems: items,
	}
}

func parFuncMap() *template.FuncMap {
	funcMap := template.FuncMap{
		"toPARdate": func(datum time.Time) string {
			return datum.Format("02.01.2006")
		},
		"formatCity": func(zip, city string) string {
			if zip == "" {
				return city
			}
			return zip + " " + city
		},
		"toCIILabel": func(data interface{}) string {
			// 1. Prüfen, ob das Interface nil ist
			if data == nil {
				return ""
			}
			// 2. Wir prüfen nur: Hat das, was da reinkommt, eine Methode "GetCode()"?
			type Labeler interface {
				GetLabel() string
			}

			if v, ok := data.(Labeler); ok && v != nil {
				return v.GetLabel()
			}
			return ""
		},
		"formatEuro": func(val interface{}) string {
			var amount float64

			switch v := val.(type) {
			case float64:
				amount = v
			case float32:
				amount = float64(v)
			case int:
				amount = float64(v)
			default:
				return "0.00"
			}

			// Einheitliche Formatierung: 2 Nachkommastellen, Punkt als Trenner
			return fmt.Sprintf("%.2f", amount)
		},
	}
	return &funcMap
}
