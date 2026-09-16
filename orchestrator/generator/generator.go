// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"codeberg.org/tiny-frameworks/nexfact/orchestrator/invoice"
)

type Generator interface {
	Generate(inv *invoice.ZUGFeRDmaster) ([]byte, error)
	Profile() invoice.Profile
}
