package new

import (
	"context"
	"fmt"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var url string

var cmd = &cobra.Command{
	Use:   "new [url]",
	Short: "Open a new tab, optionally navigate to URL",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	var targetURL string
	if len(args) > 0 {
		targetURL = args[0]
	}

	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	tabId, err := cli.NewTab(context.Background(), targetURL)
	if err != nil {
		return fmt.Errorf("new tab: %w", err)
	}

	out := output.Format(map[string]string{"tabId": tabId}, root.Output)
	cmd.Println(out)
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}