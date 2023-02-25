package cmd

import (
	"github.com/ermos/polyrule/internal/pkg/command"
	"github.com/spf13/cobra"
)

var langCmd = &cobra.Command{
	Use:   "lang",
	Short: "Get supported programming languages",
	Long:  `Get supported programming languages`,
	Run:   command.RunLang,
}

func init() {
	rootCmd.AddCommand(langCmd)
}
