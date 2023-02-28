package js

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"reflect"
	"strings"
)

func ifBuilder(b *strings.Builder, name, condition string, indent int) {
	utils.Block(b, indent, fmt.Sprintf("if (%s) {", condition), func(i int) {
		utils.Indent(b, i, fmt.Sprintf("errors.push('%s');\n", name))
	}, "}\n")
}

func validatorBuilder(b *strings.Builder, vType types.Type, indent int, rules map[string]interface{}) {
	utils.Block(b, indent, "validate(input, withErrors = false) {", func(i int) {
		utils.Indent(b, i, "const errors = [];\n\n")

		for name, value := range rules {
			name = strings.ToLower(name)

			generator := ruleGenerator[name]
			if generator == nil {
				panic(fmt.Errorf(
					"%s's rule isn't currently supported by choosen programing language compiler",
					name,
				))
			}

			if err := generator(b, name, vType, value, 3); err != nil {
				panic(err)
			}
		}

		errorBuilder(b, 3)
	}, "}")
}

func errorBuilder(b *strings.Builder, indent int) {
	utils.Block(b, indent, "if (withErrors) {", func(i int) {
		utils.Block(b, i, "return {", func(i int) {
			utils.Indent(b, i, "errors: errors,\n")
			utils.Indent(b, i, "valid: errors.length === 0,\n")
		}, "}")
	}, "}\n")

	utils.Indent(b, indent, "return errors.length === 0\n")
}

func messageBuilder(b *strings.Builder, indent int, key interface{}, v interface{}) {
	if key != nil {
		utils.Indent(b, indent, fmt.Sprintf("%v: ", key))
	} else {
		utils.Indent(b, indent, "")
	}

	ref := reflect.TypeOf(v)
	if ref.Kind() == reflect.Array || ref.Kind() == reflect.Slice {
		b.WriteString("[\n")

		m, ok := v.([]interface{})
		if ok {
			for _, value := range m {
				messageBuilder(b, indent+1, nil, value)
			}
		}

		utils.Indent(b, indent, "],\n")
	} else if ref.Kind() == reflect.Map {
		b.WriteString("{\n")

		m, ok := v.(map[string]interface{})
		if ok {
			for name, value := range m {
				messageBuilder(b, indent+1, name, value)
			}
		}

		utils.Indent(b, indent, "},\n")
	} else if ref.Kind() == reflect.String {
		b.WriteString(fmt.Sprintf("'%s',\n", strings.ReplaceAll(v.(string), "'", "\\'")))
	} else {
		// number or boolean ?
		b.WriteString(fmt.Sprintf("%v,\n", v))
	}
}
