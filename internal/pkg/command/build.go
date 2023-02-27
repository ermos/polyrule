package command

import (
	"github.com/ermos/polyrule/internal/pkg/compiler"
	"github.com/ermos/polyrule/internal/pkg/log"
	"github.com/spf13/cobra"
)

func RunBuild(cmd *cobra.Command, args []string) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if verbose && err == nil {
		log.SetLogLevel(log.LevelVerbose)
	}

	err = compiler.Compile(
		cmd,
		cmd.Flag("lang").Value.String(),
		cmd.Flag("input").Value.String(),
		cmd.Flag("output").Value.String(),
	)
	if err != nil {
		panic(err)
	}
}
