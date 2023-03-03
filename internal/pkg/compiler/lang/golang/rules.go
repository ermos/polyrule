package golang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"github.com/ermos/strlang"
	"reflect"
	"strings"
)

func ruleRequired(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	// Strict lang, check only string, int and float64
	if vType != types.String && vType != types.Int && vType != types.Float {
		return nil
	}

	if utils.ForceBool(value) {
		ifBuilder(b, name, fmt.Sprintf("input == %s", localToEmpty(vType)))
	}

	return nil
}

func ruleRegex(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	ref := reflect.TypeOf(value)

	if vType != types.String {
		return fmt.Errorf("%s type not allowed for regex", vType)
	}

	if ref.Kind() == reflect.Map {
		m, ok := value.(map[string]interface{})
		if ok {
			for n, v := range m {
				ifBuilder(b, fmt.Sprintf("regex.%s", n), fmt.Sprintf(
					`!regexp.MustCompile("%s").MatchString(input)`,
					strings.ReplaceAll(utils.ForceString(v), "\\", "\\\\"),
				))
			}
		}
		return nil
	}

	ifBuilder(b, name, fmt.Sprintf(
		`!regexp.MustCompile("%s").MatchString(input)`,
		strings.ReplaceAll(utils.ForceString(value), "\\", "\\\\"),
	))

	return nil
}

func ruleMin(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	return ruleMinMax(b, name, vType, value, "<")
}

func ruleMax(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	return ruleMinMax(b, name, vType, value, ">")
}

func ruleMinMax(b *strlang.Builder, name string, vType types.Type, value interface{}, sign string) error {
	f := "len(input)"
	if vType == types.Float || vType == types.Int {
		f = "input"
	}
	ifBuilder(b, name, fmt.Sprintf("%s %s %v", f, sign, value))
	return nil
}
