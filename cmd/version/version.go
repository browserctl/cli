package version

import (
	"fmt"

	"browserctl/cli/cmd/root"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("browserctl v0.1.0")
	},
}

func init() {
	root.RootCmd.AddCommand(cmd)
}