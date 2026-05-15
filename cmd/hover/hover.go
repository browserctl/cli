package hover

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "hover <sessionId> <selector>",
	Short: "hover operation",
	Args:  cobra.ExactArgs(2),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	selector := args[1]
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.Hover(context.Background(), sessionId, selector)
	if err != nil {
		return fmt.Errorf("hover: %w", err)
	}

	cmd.Println(output.Format(map[string]string{"ok": "hover"}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
