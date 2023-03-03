package lang

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/types"
	"github.com/ermos/strlang"
)

type Rule func(b *strlang.Builder, name string, vType types.Type, value interface{}) error

func GetGenerator(name string, list map[string]Rule) Rule {
	generator := list[name]
	if generator == nil {
		panic(fmt.Errorf(
			"%s's rule isn't currently supported by choosen programing language compiler",
			name,
		))
	}
	return generator
}
