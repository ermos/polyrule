package js

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/base"
	"github.com/ermos/polyrule/internal/pkg/types"
	"github.com/ermos/strlang"
)

func ifBuilder(b *strlang.Builder, name, condition string) {
	b.Block(fmt.Sprintf("if (%s) {", condition), func() {
		b.WriteStringln(fmt.Sprintf("errors.push('%s');", name))
	}, "}", 2)
}

func validatorBuilder(b *strlang.Javascript, vType types.Type, rules map[string]interface{}, generators map[string]lang.Rule) {
	b.Block("validate(input, withErrors = false) {", func() {
		b.WriteStringln("const errors = [];", 2)

		for name, value := range rules {
			if err := lang.GetGenerator(name, generators)(b.Builder, name, vType, value); err != nil {
				panic(err)
			}
		}

		b.If("withErrors", func() {
			b.Block("return {", func() {
				b.WriteStringln("errors: errors,")
				b.WriteStringln("valid: errors.length === 0,")
			}, "}")
		})

		b.WriteStringln("return errors.length === 0")
	}, "}")
}

func messageBuilder(b *strlang.Javascript, key interface{}, v interface{}) {
	base.MessageBuilder(b.Builder, key, v, true, map[string]string{
		"key":        "%v: ",
		"arrayStart": "[\n",
		"arrayEnd":   "]",
		"mapStart":   "{\n",
		"mapEnd":     "}",
		"string":     "'%s'",
		"number":     "%v",
		"separator":  ",\n",
		"close":      ",\n",
		"quote":      "'",
	})
}
