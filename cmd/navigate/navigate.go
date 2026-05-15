package navigate

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "navigate <sessionId> <url>",
	Short: "Navigate to a URL",
	Args:  cobra.ExactArgs(2),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	url := args[1]

	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	result, err := cli.Navigate(context.Background(), sessionId, url)
	if err != nil {
		return err
	}

	out := output.Format(result, root.Output)
	cmd.Println(out)
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}