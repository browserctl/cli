package fill

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "fill <sessionId> <selector> <value>",
	Short: "Fill an input by CSS selector",
	Args:  cobra.ExactArgs(3),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId, selector, value := args[0], args[1], args[2]
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.Fill(context.Background(), sessionId, selector, value)
	if err != nil {
		return fmt.Errorf("fill: %w", err)
	}

	cmd.Println(output.Format(map[string]string{"ok": "filled"}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}