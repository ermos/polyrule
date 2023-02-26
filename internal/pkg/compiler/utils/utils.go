package utils

import (
	"fmt"
	"strings"
)

func Capitalize(s string) string {
	if len(s) != 0 {
		firstChar := strings.ToUpper(string(s[0]))
		rest := strings.ToLower(s[1:])
		s = firstChar + rest
	}
	return s
}

func LowerFirst(s string) string {
	if len(s) != 0 {
		firstChar := strings.ToLower(string(s[0]))
		s = firstChar + s[1:]
	}
	return s
}

func Tab(b *strings.Builder, nb int, s string) {
	tab := strings.Repeat("\t", nb)
	b.WriteString(fmt.Sprintf("%s%s", tab, s))
}
