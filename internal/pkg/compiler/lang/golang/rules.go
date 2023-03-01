package golang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"reflect"
	"strings"
)

func ruleRequired(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	// Strict lang, check only string, int and float64
	if vType != types.String && vType != types.Int && vType != types.Float {
		return nil
	}

	if utils.ForceBool(value) {
		ifBuilder(b, name, fmt.Sprintf("input == %s", localToEmpty(vType)), indent)
	}

	return nil
}

func ruleRegex(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	ref := reflect.TypeOf(value)

	if vType != types.String {
		return fmt.Errorf("%s type not allowed for regex", vType)
	}

	if ref.Kind() == reflect.Map {
		m, ok := value.(map[string]interface{})
		if ok {
			for n, v := range m {
				ifBuilder(
					b,
					fmt.Sprintf("regex.%s", n),
					fmt.Sprintf(
						`!regexp.MustCompile("%s").MatchString(input)`,
						strings.ReplaceAll(utils.ForceString(v), "\\", "\\\\"),
					),
					indent,
				)
			}
		}
		return nil
	}

	ifBuilder(
		b,
		name,
		fmt.Sprintf(
			`!regexp.MustCompile("%s").MatchString(input)`,
			strings.ReplaceAll(utils.ForceString(value), "\\", "\\\\"),
		),
		indent,
	)

	return nil
}

func ruleMin(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	f := "len(input)"
	if vType == types.Float || vType == types.Int {
		f = "input"
	}
	ifBuilder(b, name, fmt.Sprintf("%s < %v", f, value), indent)
	return nil
}

func ruleMax(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	f := "len(input)"
	if vType == types.Float || vType == types.Int {
		f = "input"
	}
	ifBuilder(b, name, fmt.Sprintf("%s > %v", f, value), indent)
	return nil
}
