package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var SvcAddr string
var Secret   string
var Output   string
var Timeout  int

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&SvcAddr, "svc", "s", "ws://localhost:9222", "browserctl svc address")
	RootCmd.PersistentFlags().StringVar(&Secret, "secret", "", "auth secret")
	RootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "json", "output format: json, text, pretty")
	RootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 30000, "timeout in milliseconds")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./browserctl.yaml)")
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

var RootCmd = &cobra.Command{
	Use:   "browserctl",
	Short: "browserctl is a CLI for AI agents to control Chrome via CDP",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}