package php

import (
	"errors"
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/spf13/cobra"
	"path/filepath"
	"reflect"
	"strings"
)

type Lang struct{}

func (Lang) GetExtension() string {
	return "php"
}

func (Lang) Compile(cmd *cobra.Command, path, name string, rules map[string]model.Rule) (content string, err error) {
	b := &strings.Builder{}

	namespace := cmd.Flag("namespace").Value.String()
	if namespace == "" {
		err = errors.New("you must provide a namespace with -n or --namespace flag")
		return
	}

	b.WriteString("<?php\n")
	b.WriteString(fmt.Sprintf(
		"namespace %s;\n",
		strings.ReplaceAll(filepath.Join(namespace, path), "/", "\\"),
	),
	)
	b.WriteString(fmt.Sprintf("\nclass %sRules {\n", utils.Capitalize(name)))

	for n, rule := range rules {
		if err = writeRules(b, n, rule); err != nil {
			return
		}
	}

	b.WriteString("};\n")

	return b.String(), err
}

func writeRules(b *strings.Builder, name string, rule model.Rule) (err error) {
	utils.Tab(b, 1, fmt.Sprintf("protected static mixed $%s =", utils.LowerFirst(name)))

	fromInterface(b, nil, rule.Message, 1)

	utils.Tab(b, 1, fmt.Sprintf("public static function %sMessage(): mixed\n", utils.LowerFirst(name)))
	utils.Tab(b, 1, "{\n")
	utils.Tab(b, 2, fmt.Sprintf("return self::$%s;\n", utils.LowerFirst(name)))
	utils.Tab(b, 1, "}\n\n")

	//
	//        if (preg_match(`\d`, $value)) {
	//            $errors[] = "regex.number";
	//        }
	//

	utils.Tab(b, 1, fmt.Sprintf(
		"public static function %sValidate(mixed $value, bool $with_errors = false): bool|array\n",
		utils.LowerFirst(name),
	))
	utils.Tab(b, 1, "{\n")

	utils.Tab(b, 2, "$errors = [];\n\n")

	for ruleName, r := range rule.Rules {
		ruleName = strings.ToLower(ruleName)
		push := true

		switch ruleName {
		case "required":
			if utils.ForceBool(r) {
				utils.Tab(b, 2, "if (empty($value)) {\n")
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
						utils.Tab(b, 2, fmt.Sprintf(
							"if (empty(preg_match_all('/%s/', $value, $matches, PREG_SET_ORDER))) {\n",
							utils.ForceString(v),
						))
						utils.Tab(b, 3, fmt.Sprintf("$errors[] = 'regex.%s';\n", n))
						utils.Tab(b, 2, "}\n\n")
					}
				}
			} else {
				utils.Tab(b, 2, fmt.Sprintf(
					"if (empty(preg_match_all('/%s/', $value, $matches, PREG_SET_ORDER))) {\n",
					utils.ForceString(r),
				))
			}
			break
		case "min", "max":
			operator := "<"
			if ruleName == "max" {
				operator = ">"
			}

			utils.Tab(
				b,
				2,
				fmt.Sprintf(
					"if (((is_array($value) ? count($value) : 0) + "+
						"(is_string($value) ? strlen($value) : 0) + "+
						"(is_numeric($value) ? $value : 0)) %s %v) {\n",
					operator,
					r,
				),
			)

			break
		default:
			return fmt.Errorf(
				"%s's rule isn't currently supported by choosen programing language compiler",
				ruleName,
			)
		}

		if push {
			utils.Tab(b, 3, fmt.Sprintf("$errors[] = '%s';\n", ruleName))
			utils.Tab(b, 2, "}\n\n")
		}
	}

	utils.Tab(b, 2, "if ($with_errors) {\n")
	utils.Tab(b, 3, "return [\n")
	utils.Tab(b, 4, "'errors' => $errors,\n")
	utils.Tab(b, 4, "'valid' => empty($errors)\n")
	utils.Tab(b, 3, "];\n")
	utils.Tab(b, 2, "}\n\n")

	utils.Tab(b, 2, "return empty($errors);\n")

	utils.Tab(b, 1, "}\n")
	return
}

func fromInterface(b *strings.Builder, key interface{}, v interface{}, currTab int) {
	if key != nil {
		utils.Tab(b, currTab, fmt.Sprintf("'%v' => ", key))
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

		utils.Tab(b, currTab, "]")
	} else if ref.Kind() == reflect.Map {
		b.WriteString("[\n")

		m, ok := v.(map[string]interface{})
		if ok {
			for name, value := range m {
				fromInterface(b, name, value, currTab+1)
			}
		}

		utils.Tab(b, currTab, "]")
	} else if ref.Kind() == reflect.String {
		b.WriteString(fmt.Sprintf("'%s'", strings.ReplaceAll(v.(string), "'", "\\'")))
	} else {
		// number or boolean ?
		b.WriteString(fmt.Sprintf("%v", v))
	}

	if currTab == 1 {
		b.WriteString(";\n\n")
	} else {
		b.WriteString(",\n")
	}
}
