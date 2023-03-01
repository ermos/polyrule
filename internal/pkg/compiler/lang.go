package compiler

import (
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/spf13/cobra"
)

type Lang interface {
	New(cmd *cobra.Command, outputPath, subPath, fileName string, rules map[string]model.Rule) interface{}
	OutputDirPath() string
	Compile() (content string, err error)
}
