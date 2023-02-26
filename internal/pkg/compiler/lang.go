package compiler

import (
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/spf13/cobra"
)

type Lang interface {
	GetExtension() string
	Compile(cmd *cobra.Command, path, name string, rules map[string]model.Rule) (content string, err error)
}
