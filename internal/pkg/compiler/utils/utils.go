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

func EscapeSimple(v interface{}) string {
	return strings.ReplaceAll(ForceString(v), "'", "\\'")
}

func EscapeDouble(v interface{}) string {
	return strings.ReplaceAll(ForceString(v), "\"", "\\\"")
}

func Indent(b *strings.Builder, nb int, s string) {
	tab := strings.Repeat("\t", nb)
	b.WriteString(fmt.Sprintf("%s%s", tab, s))
}

func Block(b *strings.Builder, indent int, start string, f func(indent int), end string) {
	Indent(b, indent, start+"\n")
	f(indent + 1)
	Indent(b, indent, end+"\n")
}

func Jump(b *strings.Builder, nb int) {
	jump := strings.Repeat("\n", nb)
	b.WriteString(jump)
}
