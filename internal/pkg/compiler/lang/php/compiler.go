package php

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

type Lang struct{}

var ruleGenerator = map[string]lang.Rule{
	"required": ruleRequired,
	"regex":    ruleRegex,
	"min":      ruleMin,
	"max":      ruleMax,
}

func (Lang) GetExtension() string {
	return "php"
}

func (Lang) Compile(cmd *cobra.Command, path, name string, rules map[string]model.Rule) (content string, err error) {
	b := &strings.Builder{}

	b.WriteString("<?php\n")

	namespace := cmd.Flag("namespace").Value.String()
	if namespace != "" {
		b.WriteString(fmt.Sprintf(
			"namespace %s;\n",
			strings.ReplaceAll(filepath.Join(namespace, path), "/", "\\"),
		))
	}

	utils.Block(b, 0, fmt.Sprintf("\nclass %sRules {\n", utils.Capitalize(name)), func(i int) {
		for n, rule := range rules {
			if err = writeRules(b, n, rule); err != nil {
				return
			}
		}
	}, "};")

	return b.String(), err
}

func writeRules(b *strings.Builder, name string, rule model.Rule) (err error) {
	utils.Indent(b, 1, fmt.Sprintf("protected static mixed $%s =", utils.LowerFirst(name)))
	messageBuilder(b, 1, nil, rule.Message)
	utils.Jump(b, 1)

	m := fmt.Sprintf("public static function %sMessage(): mixed {", utils.LowerFirst(name))
	utils.Block(b, 1, m, func(i int) {
		utils.Indent(b, i, fmt.Sprintf("return self::$%s;\n", utils.LowerFirst(name)))
	}, "}\n")

	m = fmt.Sprintf(
		"public static function %sValidate(%s $value, bool $with_errors = false): %s|array {",
		utils.LowerFirst(name),
		mapType(rule.Type),
		mapType(rule.Type),
	)
	utils.Block(b, 1, m, func(i int) {
		utils.Indent(b, i, "$errors = [];\n\n")

		validatorBuilder(b, rule.Type, i, rule.Rules)

		utils.Block(b, i, "if ($with_errors) {", func(i int) {
			utils.Block(b, i, "return [", func(i int) {
				utils.Indent(b, i, "'errors' => $errors,\n")
				utils.Indent(b, i, "'valid' => empty($errors)\n")
			}, "];")
		}, "}\n")

		utils.Indent(b, 2, "return empty($errors);\n")
	}, "}\n")

	return
}
