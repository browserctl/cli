package cookies

import (
	"context"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "cookies <sessionId>",
	Short: "cookies operation",
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

	cookieStr, err := cli.Cookies(context.Background(), sessionId)
	if err != nil {
		return err
	}

	cmd.Println(output.Format(cookieStr, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
