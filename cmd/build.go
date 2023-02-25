package cmd

import (
	"github.com/ermos/polyrule/internal/pkg/command"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your rules in given programming languages",
	Long:  `Build your rules in given programming languages`,
	Run:   command.RunBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP(
		"input",
		"i",
		"",
		"specifies the input file or directory for the rule compilation process",
	)
	buildCmd.MarkFlagRequired("input")

	buildCmd.Flags().StringP(
		"output",
		"o",
		"",
		"specifies the output file or directory for the compiled rules",
	)
	buildCmd.MarkFlagRequired("output")

	buildCmd.Flags().StringP(
		"lang",
		"l",
		"",
		"specifies the programming language to use for compiling the rules",
	)
	buildCmd.MarkFlagRequired("lang")
}
