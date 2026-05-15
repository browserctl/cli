package tabs

import (
	"context"
	"encoding/json"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "tabs",
	Short: "List all http(s) tabs",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	tabs, err := cli.GetTabs(context.Background())
	if err != nil {
		return err
	}

	out := output.Format(tabs, root.Output)
	cmd.Println(out)
	return nil
}

func init() {
	root.RootCmd.AddCommand(listCmd)
}

var _ json.Marshaler = client.TargetInfo{}