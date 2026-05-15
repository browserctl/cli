package navigate

import (
	"context"

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

	if err := cli.Navigate(context.Background(), sessionId, url); err != nil {
		return err
	}

	cmd.Println(output.Format(map[string]bool{"ok": true}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
