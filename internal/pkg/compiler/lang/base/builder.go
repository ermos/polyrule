package base

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/utils"
	"reflect"
	"strings"
)

func MessageBuilder(b *strings.Builder, indent int, key interface{}, v interface{}, first bool, fragment map[string]string) {
	if key != nil {
		utils.Indent(b, indent, fmt.Sprintf(fragment["key"], key))
	} else {
		utils.Indent(b, indent, "")
	}

	ref := reflect.TypeOf(v)
	if ref.Kind() == reflect.Array || ref.Kind() == reflect.Slice {
		b.WriteString(fragment["arrayStart"])

		m, ok := v.([]interface{})
		if ok {
			for _, value := range m {
				MessageBuilder(b, indent+1, nil, value, false, fragment)
			}
		}

		utils.Indent(b, indent, fragment["arrayEnd"])
	} else if ref.Kind() == reflect.Map {
		b.WriteString(fragment["mapStart"])

		m, ok := v.(map[string]interface{})
		if ok {
			for name, value := range m {
				MessageBuilder(b, indent+1, name, value, false, fragment)
			}
		}

		utils.Indent(b, indent, fragment["mapEnd"])
	} else if ref.Kind() == reflect.String {
		b.WriteString(fmt.Sprintf(fragment["string"], strings.ReplaceAll(
			v.(string),
			fragment["quote"],
			"\\"+fragment["quote"],
		)))
	} else {
		// number or boolean ?
		b.WriteString(fmt.Sprintf(fragment["number"], v))
	}

	if first {
		b.WriteString(fragment["close"])
	} else {
		b.WriteString(fragment["separator"])
	}
}
