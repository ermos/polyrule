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

type Lang struct {
	Command        *cobra.Command
	FileName       string
	OutputPath     string
	SubPath        string
	Rules          map[string]model.Rule
	RuleGenerators map[string]lang.Rule
}

func (l Lang) OutputDirPath() string {
	return filepath.Join(l.OutputPath, l.SubPath, fmt.Sprintf("%s.php", l.FileName))
}

func (Lang) New(cmd *cobra.Command, outputPath, subPath, fileName string, rules map[string]model.Rule) interface{} {
	return Lang{
		Command:    cmd,
		OutputPath: outputPath,
		SubPath:    subPath,
		FileName:   fileName,
		Rules:      rules,
		RuleGenerators: map[string]lang.Rule{
			"required": ruleRequired,
			"regex":    ruleRegex,
			"min":      ruleMin,
			"max":      ruleMax,
		},
	}
}

func (l Lang) Compile() (content string, err error) {
	b := &strings.Builder{}

	b.WriteString("<?php\n")

	namespace := l.Command.Flag("namespace").Value.String()
	if namespace != "" {
		b.WriteString(fmt.Sprintf(
			"namespace %s;\n",
			strings.ReplaceAll(filepath.Join(namespace, l.SubPath), "/", "\\"),
		))
	}

	utils.Block(b, 0, fmt.Sprintf("\nclass %sRules {\n", utils.Capitalize(l.FileName)), func(i int) {
		for n, rule := range l.Rules {
			if err = l.writeRules(b, n, rule); err != nil {
				return
			}
		}
	}, "};")

	return b.String(), err
}

func (l Lang) writeRules(b *strings.Builder, name string, rule model.Rule) (err error) {
	utils.Indent(b, 1, fmt.Sprintf("protected static mixed $%s =", utils.LowerFirst(name)))
	messageBuilder(b, 1, rule.Message)
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

		validatorBuilder(b, rule.Type, i, rule.Rules, l.RuleGenerators)

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
