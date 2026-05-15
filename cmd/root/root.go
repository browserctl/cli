package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var svcAddr string
var secret   string
var output   string
var timeout  int

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&svcAddr, "svc", "s", "ws://localhost:9222", "browserctl svc address")
	rootCmd.PersistentFlags().StringVar(&secret, "secret", "", "auth secret")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "json", "output format: json, text, pretty")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 30000, "timeout in milliseconds")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./browserctl.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("browserctl")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/browserctl")
	}
	viper.AutomaticEnv()
}

var rootCmd = &cobra.Command{
	Use:   "browserctl",
	Short: "browserctl is a CLI for AI agents to control Chrome via CDP",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}