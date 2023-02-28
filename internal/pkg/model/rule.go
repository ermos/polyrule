package model

import "github.com/ermos/polyrule/internal/pkg/types"

type Rule struct {
	Path    string                 `json:"-"`
	Type    types.Type             `json:"type"`
	Message interface{}            `json:"message"`
	Rules   map[string]interface{} `json:"rules"`
}
