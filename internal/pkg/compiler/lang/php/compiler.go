package php

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/ermos/strlang"
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
	namespace := l.Command.Flag("namespace").Value.String()
	if namespace != "" {
		namespace = strings.ReplaceAll(filepath.Join(namespace, l.SubPath), "/", "\\")
	}

	b := strlang.NewPHP(namespace)

	b.Class(fmt.Sprintf("%sRules", utils.Capitalize(l.FileName)), func() {
		for n, rule := range l.Rules {
			if err = l.writeRules(b, n, rule); err != nil {
				return
			}
		}
	})

	return b.String(), err
}

func (l Lang) writeRules(b *strlang.PHP, name string, rule model.Rule) (err error) {
	name = utils.LowerFirst(name)

	b.WriteString(fmt.Sprintf("protected static mixed $%s =", name))
	messageBuilder(b, rule.Message)
	b.WriteStringln("")

	b.ClassFunc("public static", fmt.Sprintf("%sMessage", name), "", "mixed", func() {
		b.WriteStringln(fmt.Sprintf("return self::$%s;", name))
	})

	b.ClassFunc(
		"public static",
		fmt.Sprintf("%sValidate", name),
		fmt.Sprintf("%s $value, bool $with_errors = false", mapType(rule.Type)),
		fmt.Sprintf("%s|array", mapType(rule.Type)),
		func() {
			b.WriteStringln("$errors = [];", 2)

			validatorBuilder(b, rule.Type, rule.Rules, l.RuleGenerators)

			b.If("$with_errors", func() {
				b.Block("return [", func() {
					b.WriteStringln("'errors' => $errors,")
					b.WriteStringln("'valid' => empty($errors)")
				}, "];")
			}, 2)

			b.WriteStringln("return empty($errors);")
		},
	)

	return
}
