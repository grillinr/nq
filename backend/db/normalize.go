package db

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// NormalizeName returns a normalized, lowercase, diacritics-stripped
// version of the input suitable for deduplication and searching.
func NormalizeName(s string) string {
	if s == "" {
		return ""
	}

	// Trim and lowercase
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Decompose Unicode characters and remove diacritic marks
	t := norm.NFD.String(s)
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			// skip combining marks
			continue
		}
		b.WriteRune(r)
	}
	s = b.String()

	// Replace any non-alphanumeric characters with a single space
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, " ")

	// Collapse multiple spaces and trim
	s = strings.Join(strings.Fields(s), " ")
	return s
}
