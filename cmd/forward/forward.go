package forward

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "forward <sessionId>",
	Short: "forward operation",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close() //nolint: errcheck

	err = cli.Forward(context.Background(), sessionId)
	if err != nil {
		return fmt.Errorf("forward: %w", err)
	}

	cmd.Println(output.Format(map[string]string{"ok": "forward"}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
