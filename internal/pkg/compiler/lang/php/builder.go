package php

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/base"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"strings"
)

func ifBuilder(b *strings.Builder, name, condition string, indent int) {
	utils.Block(b, indent, fmt.Sprintf("if (%s) {", condition), func(i int) {
		utils.Indent(b, i, fmt.Sprintf("errors.push('%s');\n", name))
	}, "}\n")
}

func validatorBuilder(b *strings.Builder, vType types.Type, indent int, rules map[string]interface{}) {
	for name, value := range rules {
		name = strings.ToLower(name)

		generator := ruleGenerator[name]
		if generator == nil {
			panic(fmt.Errorf(
				"%s's rule isn't currently supported by choosen programing language compiler",
				name,
			))
		}

		if err := generator(b, name, vType, value, indent); err != nil {
			panic(err)
		}
	}
}

func messageBuilder(b *strings.Builder, indent int, key interface{}, v interface{}) {
	base.MessageBuilder(b, indent, key, v, true, map[string]string{
		"key":        "'%v' => ",
		"arrayStart": "[\n",
		"arrayEnd":   "]",
		"mapStart":   "[\n",
		"mapEnd":     "]",
		"string":     "'%s'",
		"number":     "%v",
		"separator":  ",\n",
		"close":      ";\n",
	})
}
