package js

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/model"
	"reflect"
	"strings"
)

type Lang struct{}

func (Lang) Compile(list map[string]map[string]model.Rule) (content string, err error) {
	b := &strings.Builder{}

	for namespace, rules := range list {
		b.WriteString(fmt.Sprintf("\nexport const %sRules = {\n", utils.Capitalize(namespace)))

		for name, rule := range rules {
			if err = writeRules(b, name, rule); err != nil {
				return
			}
		}

		b.WriteString("};\n")
	}

	return b.String(), err
}

func writeRules(b *strings.Builder, name string, rule model.Rule) (err error) {
	utils.Tab(b, 1, fmt.Sprintf("%s: {\n", name))

	fromInterface(b, "message", rule.Message, 2)

	utils.Tab(b, 2, "validate(input, withErrors = false) {\n")

	utils.Tab(b, 3, "const errors = [];\n\n")

	for ruleName, r := range rule.Rules {
		ruleName = strings.ToLower(ruleName)
		push := true

		switch ruleName {
		case "required":
			if utils.ForceBool(r) {
				utils.Tab(b, 3, "if (!input) {\n")
			} else {
				push = false
			}
			break
		case "regex":
			ref := reflect.TypeOf(r)
			if ref.Kind() == reflect.Map {
				push = false

				m, ok := r.(map[string]interface{})
				if ok {
					for n, v := range m {
						utils.Tab(b, 3, fmt.Sprintf("if (!/%s/.test(input)) {\n", utils.ForceString(v)))
						utils.Tab(b, 4, fmt.Sprintf("errors.push('regex.%s');\n", n))
						utils.Tab(b, 3, "}\n\n")
					}
				}
			} else {
				utils.Tab(b, 3, fmt.Sprintf("if (!/%s/.test(input)) {\n", utils.ForceString(r)))
			}
			break
		case "min":
			utils.Tab(b, 3, fmt.Sprintf("if (input.length < %v) {\n", r))
			break
		case "max":
			utils.Tab(b, 3, fmt.Sprintf("if (input.length > %v) {\n", r))
			break
		default:
			return fmt.Errorf(
				"%s's rule isn't currently supported by choosen programing language compiler",
				ruleName,
			)
		}

		if push {
			utils.Tab(b, 4, fmt.Sprintf("errors.push('%s');\n", ruleName))
			utils.Tab(b, 3, "}\n\n")
		}
	}

	utils.Tab(b, 3, "if (withErrors) {\n")
	utils.Tab(b, 4, "return {\n")
	utils.Tab(b, 5, "errors: errors,\n")
	utils.Tab(b, 5, "valid: errors.length === 0,\n")
	utils.Tab(b, 4, "}\n")
	utils.Tab(b, 3, "}\n\n")

	utils.Tab(b, 3, "return errors.length === 0\n")

	utils.Tab(b, 2, "},\n")

	utils.Tab(b, 1, "},\n")

	return
}

func fromInterface(b *strings.Builder, key interface{}, v interface{}, currTab int) {
	if key != nil {
		utils.Tab(b, currTab, fmt.Sprintf("%v: ", key))
	} else {
		utils.Tab(b, currTab, "")
	}

	ref := reflect.TypeOf(v)
	if ref.Kind() == reflect.Array || ref.Kind() == reflect.Slice {
		b.WriteString("[\n")

		m, ok := v.([]interface{})
		if ok {
			for _, value := range m {
				fromInterface(b, nil, value, currTab+1)
			}
		}

		utils.Tab(b, currTab, "],\n")
	} else if ref.Kind() == reflect.Map {
		b.WriteString("{\n")

		m, ok := v.(map[string]interface{})
		if ok {
			for name, value := range m {
				fromInterface(b, name, value, currTab+1)
			}
		}

		utils.Tab(b, currTab, "},\n")
	} else if ref.Kind() == reflect.String {
		b.WriteString(fmt.Sprintf("'%s',\n", strings.ReplaceAll(v.(string), "'", "\\'")))
	} else {
		// number or boolean ?
		b.WriteString(fmt.Sprintf("%v,\n", v))
	}
}
