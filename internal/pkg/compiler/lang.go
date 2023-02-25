package compiler

import "github.com/ermos/polyrule/internal/pkg/model"

type Lang interface {
	Compile(list map[string]map[string]model.Rule) (content string, err error)
}
