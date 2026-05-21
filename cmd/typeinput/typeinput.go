package typeinput

import (
	"context"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "type <sessionId> <selector> <text>",
	Short: "type text into an element",
	Args:  cobra.ExactArgs(3),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	selector := args[1]
	text := args[2]

	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close() //nolint: errcheck

	if err := cli.Type(context.Background(), sessionId, selector, text); err != nil {
		return err
	}

	cmd.Println(output.Format(map[string]bool{"ok": true}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}