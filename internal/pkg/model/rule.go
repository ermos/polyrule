package model

type Rule struct {
	Message interface{}            `json:"message"`
	Rules   map[string]interface{} `json:"rules"`
}
