package fsm

import (
	"strings"
	"unicode"
)

// ToCamelCase transform string to upper or lower camel case
func ToCamelCase(s string, upper bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var (
		nextUpper = upper
		prev      rune
		builder   strings.Builder
	)

	for _, r := range s {
		if !unicode.IsLetter(r) {
			nextUpper = true
			continue
		}

		if nextUpper {
			builder.WriteRune(unicode.ToUpper(r))
			nextUpper = false
			continue
		}

		if unicode.IsLower(prev) {
			builder.WriteRune(r)
			continue
		}

		builder.WriteRune(unicode.ToLower(r))
		prev = r
	}

	return builder.String()
}
