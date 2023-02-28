package php

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/errors"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"github.com/ermos/polyrule/internal/pkg/types"
	"reflect"
	"strings"
)

func ruleRequired(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	if utils.ForceBool(value) {
		ifBuilder(b, name, "empty($value)", indent)
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
					fmt.Sprintf("empty(preg_match_all('/%s/', $value, $matches, PREG_SET_ORDER))", utils.EscapeSimple(v)),
					indent,
				)
			}
		}
		return nil
	}

	ifBuilder(
		b,
		name,
		fmt.Sprintf("empty(preg_match_all('/%s/', $value, $matches, PREG_SET_ORDER))", utils.EscapeSimple(value)),
		indent,
	)

	return nil
}

func ruleMin(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	return ruleMinMax(b, name, vType, value, indent, "<")
}

func ruleMax(b *strings.Builder, name string, vType types.Type, value interface{}, indent int) error {
	return ruleMinMax(b, name, vType, value, indent, ">")
}

func ruleMinMax(b *strings.Builder, name string, vType types.Type, value interface{}, indent int, operator string) error {
	var c string
	if vType == types.Int || vType == types.Float {
		c = "$value"
	} else if vType == types.String {
		c = "strlen($value)"
	} else {
		return errors.UnsupportedTypeForRule(vType, "min/max")
	}

	ifBuilder(b, name, fmt.Sprintf("%s %s %v", c, operator, value), indent)
	return nil
}
