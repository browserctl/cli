package eval

import (
	"context"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"
	"browserctl/cli/output"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "eval <sessionId> <expr>",
	Short: "Evaluate JavaScript",
	Args:  cobra.ExactArgs(2),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	expr := args[1]

	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close() //nolint: errcheck

	result, err := cli.Eval(context.Background(), sessionId, expr)
	if err != nil {
		return err
	}

	cmd.Println(output.Format(result, root.Output))
	return nil
}

func init() {
	root.RootCmd.AddCommand(cmd)
}
