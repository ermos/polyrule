package golang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/ermos/strlang"
	"github.com/spf13/cobra"
	"path/filepath"
)

type Lang struct {
	Command        *cobra.Command
	FileName       string
	OutputPath     string
	SubPath        string
	Rules          map[string]model.Rule
	RuleGenerators map[string]lang.Rule
	StructName     string
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

func (l Lang) OutputDirPath() string {
	return filepath.Join(l.OutputPath, l.SubPath, fmt.Sprintf("%s.go", l.FileName))
}

func (l Lang) Compile() (content string, err error) {
	l.StructName = fmt.Sprintf("%sRules", utils.ToPascal(l.FileName))

	b := strlang.NewGolang(l.getPackageName())

	b.Struct(l.StructName, func() {})

	for n, rule := range l.Rules {
		// GetNameMessage() function
		b.Func(
			l.StructName,
			fmt.Sprintf("Get%sMessage", utils.ToPascal(n)),
			"",
			l.interfaceToType(rule.Message),
			func() {
				b.WriteString("return ")
				messageBuilder(b, rule.Message)
			},
		)

		// ValidateNameWithErrors() function
		b.Func(
			l.StructName,
			fmt.Sprintf("Validate%sWithErrors", utils.ToPascal(n)),
			"input "+localToType(rule.Type),
			"isValid bool, errors []string",
			func() {
				validatorBuilder(b, rule.Type, rule.Rules, l.RuleGenerators)
			},
		)

		// Validate() function
		b.Func(
			"r "+l.StructName,
			fmt.Sprintf("Validate%s", utils.ToPascal(n)),
			"input "+localToType(rule.Type),
			"isValid bool",
			func() {
				b.WriteStringln(fmt.Sprintf("isValid, _ = r.Validate%sWithErrors(input)", utils.ToPascal(n)))
				b.WriteStringln("return")
			},
		)
	}

	return b.String(), err
}
