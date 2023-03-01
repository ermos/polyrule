package utils

import (
	"unicode"
)

func ToSnake(s string) (result string) {
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i > 0 && s[i-1] != '_' {
				result += "_"
			}
			result += string(unicode.ToLower(c))
			continue
		}
		result += string(c)
	}
	return
}

func ToPascal(s string) (result string) {
	return handlePascalAndCamel(s, true)
}

func ToCamel(s string) (result string) {
	return handlePascalAndCamel(s, false)
}

func handlePascalAndCamel(s string, isPascal bool) (result string) {
	capitalize := isPascal
	for _, c := range s {
		if c == '_' {
			capitalize = true
		} else if capitalize {
			result += string(unicode.ToUpper(c))
			capitalize = false
		} else {
			result += string(c)
		}
	}
	return
}
