package js

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"github.com/ermos/strlang"
	"reflect"
)

func ruleRequired(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	if utils.ForceBool(value) {
		ifBuilder(b, name, "!input")
	}
	return nil
}

func ruleRegex(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	ref := reflect.TypeOf(value)

	if ref.Kind() == reflect.Map {
		m, ok := value.(map[string]interface{})
		if ok {
			for n, v := range m {
				ifBuilder(
					b,
					fmt.Sprintf("regex.%s", n),
					fmt.Sprintf("!/%s/.test(input)", utils.ForceString(v)),
				)
			}
		}
		return nil
	}

	ifBuilder(
		b,
		name,
		fmt.Sprintf("!/%s/.test(input)", utils.ForceString(value)),
	)

	return nil
}

func ruleMin(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	return ruleMinMax(b, name, vType, value, "<")
}

func ruleMax(b *strlang.Builder, name string, vType types.Type, value interface{}) error {
	return ruleMinMax(b, name, vType, value, ">")
}

func ruleMinMax(b *strlang.Builder, name string, vType types.Type, value interface{}, sign string) error {
	f := "input.length"
	if vType == types.Float || vType == types.Int {
		f = "input"
	}
	ifBuilder(b, name, fmt.Sprintf("%s %s %v", f, sign, value))
	return nil
}
