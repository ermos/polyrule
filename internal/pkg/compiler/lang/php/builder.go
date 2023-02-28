package php

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
	if key != nil {
		utils.Indent(b, indent, fmt.Sprintf("'%v' => ", key))
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

		utils.Indent(b, indent, "]")
	} else if ref.Kind() == reflect.Map {
		b.WriteString("[\n")

		m, ok := v.(map[string]interface{})
		if ok {
			for name, value := range m {
				messageBuilder(b, indent+1, name, value)
			}
		}

		utils.Indent(b, indent, "]")
	} else if ref.Kind() == reflect.String {
		b.WriteString(fmt.Sprintf("'%s'", strings.ReplaceAll(v.(string), "'", "\\'")))
	} else {
		// number or boolean ?
		b.WriteString(fmt.Sprintf("%v", v))
	}

	if indent == 1 {
		b.WriteString(";\n\n")
	} else {
		b.WriteString(",\n")
	}
}
