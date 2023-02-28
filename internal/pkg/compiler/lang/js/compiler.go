package js

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/spf13/cobra"
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
	return "js"
}

func (Lang) Compile(cmd *cobra.Command, path, name string, rules map[string]model.Rule) (content string, err error) {
	b := &strings.Builder{}

	utils.Block(b, 0, fmt.Sprintf("export const %sRules = {", utils.Capitalize(name)), func(i int) {
		for n, rule := range rules {
			utils.Block(b, 1, fmt.Sprintf("%s: {", n), func(i int) {
				messageBuilder(b, i, "message", rule.Message)
				validatorBuilder(b, rule.Type, i, rule.Rules)
			}, "}")
		}
	}, "};")

	return b.String(), err
}
