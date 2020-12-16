// Copyright © Trevor N. Suarez (Rican7)

// Package normalize provides mechanisms to normalize PlayStation software names
// and references.
package normalize

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regexSerialCodeLoose  = regexp.MustCompile(`([A-Za-z]+).*?(\d+)`)
	regexSerialCodeStrict = regexp.MustCompile(`^([A-Z]{4})-(\d{5})$`)
)

// SerialCode takes a serial code string and returns a normalized variant.
func SerialCode(serialCode string) string {
	normalized := serialCode

	normalized = strings.TrimSpace(normalized)
	normalized = strings.ToUpper(normalized)

	serialCodeMatches := regexSerialCodeLoose.FindStringSubmatch(normalized)
	if len(serialCodeMatches) == 3 {
		normalized = fmt.Sprintf("%s-%s", serialCodeMatches[1], serialCodeMatches[2])
	}

	if !regexSerialCodeStrict.MatchString(normalized) {
		return ""
	}

	return normalized
}

// Title takes a title string and returns a normalized variant and any common
// variations of that title.
func Title(title string) (string, []string) {
	normalized := title

	normalized = strings.TrimSpace(normalized)

	return normalized, nil
}
