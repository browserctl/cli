package scroll

import (
	"context"
	"strconv"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "scroll <sessionId> <px>",
	Short: "scroll operation",
	Args:  cobra.ExactArgs(2),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	px, _ := strconv.Atoi(args[1])
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.Scroll(context.Background(), sessionId, px)
	if err != nil {
		return fmt.Errorf("scroll: %w", err)
	}

	cmd.Println(output.Format(map[string]string{"ok": "scroll"}, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
