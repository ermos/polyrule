package command

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler"
	"github.com/spf13/cobra"
)

func RunLang(cmd *cobra.Command, args []string) {
	fmt.Println("Available programming language :")
	for _, v := range compiler.GetAvailableLang() {
		fmt.Printf("- %s\n", v)
	}
}
