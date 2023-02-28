package types

type Type string

const (
	String Type = "string"
	Bool   Type = "boolean"
	Float  Type = "float"
	Int    Type = "int"
)

func IsValidType(t Type) bool {
	return t == String ||
		t == Bool ||
		t == Float ||
		t == Int
}
