package golang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/base"
	"github.com/ermos/polyrule/internal/pkg/types"
	"github.com/ermos/strlang"
)

func ifBuilder(b *strlang.Builder, name, condition string) {
	b.Block(fmt.Sprintf("if %s {", condition), func() {
		b.WriteStringln(fmt.Sprintf(`errors = append(errors, "%s")`, name))
	}, "}", 2)
}

func validatorBuilder(b *strlang.Golang, vType types.Type, rules map[string]interface{}, generators map[string]lang.Rule) {
	for name, value := range rules {
		if name == "regex" {
			b.AddImports("regexp")
		}

		if err := lang.GetGenerator(name, generators)(b.Builder, name, vType, value); err != nil {
			panic(err)
		}
	}

	b.WriteStringln("return len(errors) == 0, errors")

	return
}

func messageBuilder(b *strlang.Golang, v interface{}) {
	base.MessageBuilder(b.Builder, nil, v, true, map[string]string{
		"key":        `"%v": `,
		"arrayStart": "[]interface{} {\n",
		"arrayEnd":   "}",
		"mapStart":   "map[string]interface{} {\n",
		"mapEnd":     "}",
		"string":     `"%s"`,
		"number":     "%v",
		"separator":  ",\n",
		"close":      "\n",
		"quote":      "\"",
	})
}
