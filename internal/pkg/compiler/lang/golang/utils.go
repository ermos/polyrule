package golang

import (
	"github.com/ermos/polyrule/internal/pkg/types"
	"path/filepath"
)

func (l Lang) getPackageName() string {
	return filepath.Base(filepath.Dir(l.OutputDirPath()))
}

func (Lang) interfaceToType(v interface{}) (t string) {
	switch v.(type) {
	case map[string]interface{}:
		t = "map[string]interface{}"
		break
	case string:
		t = "string"
		break
	case bool:
		t = "bool"
		break
	case float64:
		t = "float64"
		break
	case []interface{}:
		t = "[]interface{}"
		break
	default:
		t = "interface{}"
	}
	return
}

func localToType(lt types.Type) (t string) {
	switch lt {
	case types.String:
		t = "string"
		break
	case types.Int:
		t = "int"
		break
	case types.Float:
		t = "float64"
		break
	case types.Bool:
		t = "bool"
		break
	default:
		t = "interface{}"
	}
	return
}

func localToEmpty(t types.Type) (e string) {
	switch t {
	case types.String:
		e = `""`
		break
	case types.Int, types.Float:
		e = "0"
		break
	case types.Bool:
		e = "false"
		break
	default:
		e = "nil"
	}
	return
}
