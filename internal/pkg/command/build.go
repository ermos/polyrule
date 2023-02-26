package command

import (
	"github.com/ermos/polyrule/internal/pkg/compiler"
	"github.com/spf13/cobra"
)

func RunBuild(cmd *cobra.Command, args []string) {
	err := compiler.Compile(
		cmd,
		cmd.Flag("lang").Value.String(),
		cmd.Flag("input").Value.String(),
		cmd.Flag("output").Value.String(),
	)
	if err != nil {
		panic(err)
	}
}
