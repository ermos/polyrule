package model

type Rule struct {
	Path    string                 `json:"-"`
	Message interface{}            `json:"message"`
	Rules   map[string]interface{} `json:"rules"`
}
