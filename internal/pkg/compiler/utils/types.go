package utils

func ForceBool(v interface{}) bool {
	if _, ok := v.(bool); ok {
		return v.(bool)
	}
	return false
}

func ForceString(v interface{}) string {
	if _, ok := v.(string); ok {
		return v.(string)
	}
	return ""
}
