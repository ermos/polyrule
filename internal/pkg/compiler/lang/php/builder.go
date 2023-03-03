package php

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/base"
	"github.com/ermos/polyrule/internal/pkg/types"
	"github.com/ermos/strlang"
)

func ifBuilder(b *strlang.Builder, name, condition string) {
	b.Block(fmt.Sprintf("if (%s) {", condition), func() {
		b.WriteStringln(fmt.Sprintf(`$errors[] = "%s";`, name))
	}, "}", 2)
}

func validatorBuilder(b *strlang.PHP, vType types.Type, rules map[string]interface{}, generators map[string]lang.Rule) {
	for name, value := range rules {
		if err := lang.GetGenerator(name, generators)(b.Builder, name, vType, value); err != nil {
			panic(err)
		}
	}
}

func messageBuilder(b *strlang.PHP, v interface{}) {
	base.MessageBuilder(b.Builder, nil, v, true, map[string]string{
		"key":        "'%v' => ",
		"arrayStart": "[\n",
		"arrayEnd":   "]",
		"mapStart":   "[\n",
		"mapEnd":     "]",
		"string":     "'%s'",
		"number":     "%v",
		"separator":  ",\n",
		"close":      ";\n",
		"quote":      "'",
	})
}
