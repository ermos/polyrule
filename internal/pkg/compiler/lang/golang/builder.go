package golang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/base"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"strings"
)

func ifBuilder(b *strings.Builder, name, condition string, indent int) {
	utils.Block(b, indent, fmt.Sprintf("if %s {", condition), func(i int) {
		utils.Indent(b, i, fmt.Sprintf("errors = append(errors, \"%s\")\n", name))
	}, "}\n")
}

func validatorBuilder(b *strings.Builder, vType types.Type, i int, rules map[string]interface{}, ruleGenerator map[string]lang.Rule) (imports []string) {
	for name, value := range rules {
		name = strings.ToLower(name)

		generator := ruleGenerator[name]
		if generator == nil {
			panic(fmt.Errorf(
				"%s's rule isn't currently supported by choosen programing language compiler",
				name,
			))
		}

		if name == "regex" {
			imports = append(imports, "regexp")
		}

		if err := generator(b, name, vType, value, i); err != nil {
			panic(err)
		}
	}

	utils.Indent(b, i, "return len(errors) == 0, errors\n")

	return
}

func messageBuilder(b *strings.Builder, indent int, key interface{}, v interface{}) {
	base.MessageBuilder(b, indent, nil, v, true, map[string]string{
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
