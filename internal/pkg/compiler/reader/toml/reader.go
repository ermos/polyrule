package toml

import (
	"github.com/BurntSushi/toml"
	"github.com/ermos/polyrule/internal/pkg/model"
)

type Reader struct{}

func (Reader) Read(input string) (rules map[string]model.Rule, err error) {
	_, err = toml.DecodeFile(input, &rules)
	return
}
