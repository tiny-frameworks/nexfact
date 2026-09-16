// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

package invoice

import (
	"math"
	"strings"
	"time"
)

// Round to 2 decimal places (standard rounding)
func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

func clean(input string) string {
	// Replace semicolons with a space or a comma
	// so that macro parsing does not break
	safe := strings.ReplaceAll(input, ";", ",")

	// Remove line breaks, as they would corrupt the .par file.
	safe = strings.ReplaceAll(safe, "\n", " ")

	return strings.TrimSpace(safe)
}

func date(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
