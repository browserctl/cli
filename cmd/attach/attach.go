package attach

import (
	"context"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "attach <tabId>",
	Short: "Attach to a tab and get sessionId",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	tabId := args[0]
	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close() //nolint: errcheck

	sessionId, err := cli.Attach(context.Background(), tabId)
	if err != nil {
		return err
	}

	out := output.Format(map[string]string{"sessionId": sessionId}, root.Output)
	cmd.Println(out)
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}