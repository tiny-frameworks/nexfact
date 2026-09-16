// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"fmt"
	"html"
	"path/filepath"
	"text/template"
	"time"

	"codeberg.org/tiny-frameworks/nexfact/orchestrator/invoice"
	"codeberg.org/tiny-frameworks/nexutils/errors"
)

type FxGenerator struct {
	tpl *template.Template
}

func (g *FxGenerator) Profile() invoice.Profile {
	return "EN16931"
}

func NewFX(templatePath string) (Generator, error) {

	baseName := filepath.Base(templatePath)
	tpl, err := template.New(baseName).
		Funcs(*xmlFuncMap()).
		Option("missingkey=error").
		ParseFiles(templatePath)
	if err != nil {
		return nil, errors.Wrap(
			errors.InvalidValue,
			fmt.Sprintf("basename [%s] or templatepath [%s] not found", baseName, templatePath),
			"generator.en16931.New()",
			err,
		)
	}

	return &FxGenerator{tpl: tpl}, nil
}

func (g *FxGenerator) Generate(inv *invoice.ZUGFeRDmaster) ([]byte, error) {
	if inv == nil {
		return nil, errors.New(
			errors.MissingField,
			"Invoice must not be nil",
			"generator.en16931.Generate()",
		)
	}

	var buf bytes.Buffer

	if err := g.tpl.Execute(&buf, inv.Invoice); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func xmlFuncMap() *template.FuncMap {
	funcMap := template.FuncMap{

		"toCIIcode": func(data interface{}) string {
			// 1. Check if Interface is nil
			if data == nil {
				return ""
			}
			// 2. We only check: Does the object being passed in have a "GetCode()" method?
			type Coder interface {
				GetCode() string
			}

			if v, ok := data.(Coder); ok && v != nil {
				return v.GetCode()
			}
			return ""
		},

		"xmlEscape": func(s string) string {
			return html.EscapeString(s)
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

				return ""
			}

			// Consistent formatting: 2 decimal places, dot as separator
			return fmt.Sprintf("%.2f", amount)
		},
		"toXMLdate": func(datum time.Time) string {
			return datum.Format("20060102")
		},
	}
	return &funcMap
}
