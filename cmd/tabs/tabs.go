package tabs

import (
	"context"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "tabs",
	Short: "List open tabs",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close() //nolint: errcheck

	tabs, err := cli.GetTabs(context.Background())
	if err != nil {
		return err
	}

	cmd.Println(output.Format(tabs, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(listCmd)
}
