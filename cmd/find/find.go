package find

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "find <sessionId>",
	Short: "find operation",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.find(context.Background(), sessionId)
	if err != nil {
		return fmt.Errorf("find: %w", err)
	}

	cmd.Println(output.Format(map[string]string{"ok": "find"}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
