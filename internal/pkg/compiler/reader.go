package compiler

import "github.com/ermos/polyrule/internal/pkg/model"

type Reader interface {
	Read(input string) (rules map[string]model.Rule, err error)
}
