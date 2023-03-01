package golang

import (
	"github.com/ermos/polyrule/internal/pkg/types"
	"path/filepath"
	"reflect"
)

func (l Lang) getPackageName() string {
	return filepath.Base(filepath.Dir(l.OutputDirPath()))
}

func (l Lang) interfaceToType(v interface{}) string {
	t := reflect.TypeOf(v).Kind()

	if t == reflect.Map {
		return "map[string]interface{}"
	} else if t == reflect.String {
		return "string"
	} else if t == reflect.Bool {
		return "bool"
	} else if t == reflect.Float64 {
		return "float64"
	} else if t == reflect.Array {
		return "[]interface{}"
	}

	return "interface{}"
}

func localToType(t types.Type) string {
	if t == types.String {
		return "string"
	} else if t == types.Int {
		return "int"
	} else if t == types.Float {
		return "float64"
	} else if t == types.Bool {
		return "bool"
	}
	return "interface{}"
}

func localToEmpty(t types.Type) string {
	if t == types.String {
		return `""`
	} else if t == types.Int {
		return "0"
	} else if t == types.Float {
		return "0"
	} else if t == types.Bool {
		return "false"
	}
	return "nil"
}
