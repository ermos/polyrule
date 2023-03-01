package golang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/spf13/cobra"
	"path/filepath"
	"reflect"
	"strings"
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
	imports := make(map[string]bool)
	top := &strings.Builder{}
	b := &strings.Builder{}

	utils.Indent(top, 0, fmt.Sprintf("package %s\n\n", l.getPackageName()))

	l.StructName = fmt.Sprintf("%sRules", utils.ToPascal(l.FileName))
	utils.Indent(b, 0, fmt.Sprintf("type %s struct{}\n\n", l.StructName))

	for n, rule := range l.Rules {
		utils.Block(b, 0, fmt.Sprintf(
			"func (%s) Get%sMessage() %s {",
			l.StructName,
			utils.ToPascal(n),
			l.interfaceToType(rule.Message),
		), func(i int) {
			utils.Indent(b, 1, "return ")
			messageBuilder(b, i, "message", rule.Message)
		}, "}\n")

		utils.Block(b, 0, fmt.Sprintf(
			"func (%s) Validate%sWithErrors(input %s) (isValid bool, errors []string) {",
			l.StructName,
			utils.ToPascal(n),
			localToType(rule.Type),
		), func(i int) {
			localImports := validatorBuilder(b, rule.Type, i, rule.Rules, l.RuleGenerators)
			for _, localImport := range localImports {
				imports[localImport] = true
			}
		}, "}\n")

		utils.Block(b, 0, fmt.Sprintf(
			"func (r %s) Validate%s(input %s) (isValid bool) {",
			l.StructName,
			utils.ToPascal(n),
			localToType(rule.Type),
		), func(i int) {
			utils.Indent(b, i, fmt.Sprintf("isValid, _ = r.Validate%sWithErrors(input)\n", utils.ToPascal(n)))
			utils.Indent(b, i, "return\n")
		}, "}\n")
	}

	// Must get imports before write this part
	list := reflect.ValueOf(imports).MapKeys()
	if len(list) != 0 {
		utils.Block(top, 0, "import (", func(i int) {
			for _, importItem := range reflect.ValueOf(imports).MapKeys() {
				utils.Indent(top, 1, fmt.Sprintf("\"%s\"\n", importItem.String()))
			}
		}, ")\n")
	}

	return top.String() + b.String(), err
}
