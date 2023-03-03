package base

import (
	"fmt"
	"github.com/ermos/strlang"
	"reflect"
	"strings"
)

func MessageBuilder(b *strlang.Builder, key interface{}, v interface{}, first bool, fragment map[string]string) {
	keyStr := ""
	if key != nil {
		keyStr = fmt.Sprintf(fragment["key"], key)
	}

	b.WriteString(keyStr)

	ref := reflect.TypeOf(v)
	if ref.Kind() == reflect.Array || ref.Kind() == reflect.Slice {
		b.WriteNoIdentString(fragment["arrayStart"])

		m, ok := v.([]interface{})
		if ok {
			for _, value := range m {
				b.Indent()
				MessageBuilder(b, nil, value, false, fragment)
				b.StripIndent()
			}
		}

		b.WriteString(fragment["arrayEnd"])
	} else if ref.Kind() == reflect.Map {
		b.WriteNoIdentString(fragment["mapStart"])

		m, ok := v.(map[string]interface{})
		if ok {
			for name, value := range m {
				b.Indent()
				MessageBuilder(b, name, value, false, fragment)
				b.StripIndent()
			}
		}

		b.WriteString(fragment["mapEnd"])
	} else if ref.Kind() == reflect.String {
		b.WriteNoIdentString(fmt.Sprintf(fragment["string"], strings.ReplaceAll(
			v.(string),
			fragment["quote"],
			"\\"+fragment["quote"],
		)))
	} else {
		// number or boolean ?
		b.WriteNoIdentString(fmt.Sprintf(fragment["number"], v))
	}

	if first {
		b.WriteNoIdentString(fragment["close"])
	} else {
		b.WriteNoIdentString(fragment["separator"])
	}
}
