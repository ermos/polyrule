package js

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
	b := &strings.Builder{}

	utils.Block(b, 0, fmt.Sprintf("export const %sRules = {", utils.Capitalize(l.FileName)), func(i int) {
		for n, rule := range l.Rules {
			utils.Block(b, 1, fmt.Sprintf("%s: {", n), func(i int) {
				messageBuilder(b, i, "message", rule.Message)
				validatorBuilder(b, rule.Type, i, rule.Rules, l.RuleGenerators)
			}, "},")
		}
	}, "};")

	return b.String(), err
}
