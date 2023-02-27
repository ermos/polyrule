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

	// Required flags

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

	// Optional flags

	buildCmd.Flags().StringP(
		"namespace",
		"n",
		"",
		"specifies the current namespace (php)",
	)

	buildCmd.Flags().BoolP(
		"clean",
		"c",
		false,
		"remove file or directory output if already exists",
	)

	buildCmd.Flags().BoolP(
		"verbose",
		"v",
		false,
		"output more detailed information about the build process",
	)
}
