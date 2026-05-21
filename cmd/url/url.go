package url

import (
	"context"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "url <sessionId>",
	Short: "url operation",
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

	urlStr, err := cli.GetUrl(context.Background(), sessionId)
	if err != nil {
		return err
	}

	cmd.Println(output.Format(urlStr, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
