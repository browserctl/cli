package screenshot

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"browserctl/cli/client"
	"browserctl/cli/cmd/root"

	"github.com/spf13/cobra"
)

var outputPath string

var cmd = &cobra.Command{
	Use:   "screenshot <sessionId> [outputPath]",
	Short: "Take a screenshot",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	sessionId := args[0]
	path := ""
	if len(args) > 1 {
		path = args[1]
	}

	cli, err := client.New(root.SvcAddr, root.Secret)
	if err != nil {
		return err
	}
	defer cli.Close()

	data, err := cli.Screenshot(context.Background(), sessionId)
	if err != nil {
		return fmt.Errorf("screenshot: %w", err)
	}

	// Decode base64 to raw bytes for PNG
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		// data might already be raw bytes
		raw = data
	}

	if path == "" {
		os.Stdout.Write(raw)
	} else {
		if err := os.WriteFile(path, raw, 0644); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		cmd.Printf("saved to %s\n", path)
	}
	return nil
}

func init() {
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path")
	root.RootCmd.AddCommand(cmd)
}