package switchcmd

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "switch <sessionId> <tabId>",
	Short: "Switch to a different tab",
	Args:  cobra.ExactArgs(2),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	tabId := args[1]

	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	newSessionId, err := cli.Attach(context.Background(), tabId)
	if err != nil {
		return fmt.Errorf("switch: %w", err)
	}

	out := output.Format(map[string]string{"sessionId": newSessionId}, root.Output)
	cmd.Println(out)
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}