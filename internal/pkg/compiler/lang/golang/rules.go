package golang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"reflect"
	"strings"
)

func ruleRequired(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	if utils.ForceBool(value) {
		ifBuilder(b, name, "!input", indent)
	}
	return nil
}

func ruleRegex(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	ref := reflect.TypeOf(value)

	if ref.Kind() == reflect.Map {
		m, ok := value.(map[string]interface{})
		if ok {
			for n, v := range m {
				ifBuilder(
					b,
					fmt.Sprintf("regex.%s", n),
					fmt.Sprintf("!/%s/.test(input)", utils.ForceString(v)),
					indent,
				)
			}
		}
		return nil
	}

	ifBuilder(
		b,
		name,
		fmt.Sprintf("!/%s/.test(input)", utils.ForceString(value)),
		indent,
	)

	return nil
}

func ruleMin(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	ifBuilder(b, name, fmt.Sprintf("input.length < %v", value), indent)
	return nil
}

func ruleMax(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	ifBuilder(b, name, fmt.Sprintf("input.length > %v", value), indent)
	return nil
}
