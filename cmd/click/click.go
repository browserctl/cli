package click

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "click <sessionId> <selector>",
	Short: "Click an element by CSS selector",
	Args:  cobra.ExactArgs(2),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId, selector := args[0], args[1]
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.Click(context.Background(), sessionId, selector)
	if err != nil {
		return fmt.Errorf("click: %w", err)
	}

	cmd.Println(output.Format(map[string]string{"ok": "clicked"}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}