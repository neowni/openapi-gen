package common

import (
	"strings"
	"unicode"
)

func NameCamelCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_'
	})

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}

	return strings.Join(parts, "")
}

func NameSnakeCase(s string) string {
	var result strings.Builder

	for i, char := range s {
		if unicode.IsUpper(char) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(char))
	}

	return result.String()
}
