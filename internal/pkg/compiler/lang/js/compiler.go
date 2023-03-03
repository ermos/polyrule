package js

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
}

func (l Lang) OutputDirPath() string {
	return filepath.Join(l.OutputPath, l.SubPath, fmt.Sprintf("%s.js", l.FileName))
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
	b := strlang.NewJavascript()

	b.Export().Object("const", fmt.Sprintf("%sRules", utils.Capitalize(l.FileName)), func() {
		for n, rule := range l.Rules {
			b.Block(fmt.Sprintf("%s: {", n), func() {
				messageBuilder(b, "message", rule.Message)
				validatorBuilder(b, rule.Type, rule.Rules, l.RuleGenerators)
			}, "},")
		}
	})

	return b.String(), err
}
