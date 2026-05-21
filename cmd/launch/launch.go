package launch

import (
	"browserctl/cli/client"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	addr    string
	secret  string
	extPath string
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch Chrome browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New(addr, secret)
		if err != nil {
			return fmt.Errorf("connect: %w", err)
		}
		defer c.Close() //nolint: errcheck

		err = c.Launch(context.Background(), extPath)
		if err != nil {
			return fmt.Errorf("launch: %w", err)
		}
		return nil
	},
}

func init() {
	launchCmd.Flags().StringVar(&addr, "addr", "ws://localhost:9222", "WebSocket address")
	launchCmd.Flags().StringVar(&secret, "secret", "", "Secret for authentication")
	launchCmd.Flags().StringVar(&extPath, "ext", "", "Path to Chrome extension to load")
}