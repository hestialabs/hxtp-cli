package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var transportType string

var rootCmd = &cobra.Command{
	Use:     "hxtp-cli",
	Version: "v1.0.1",
	Short:   "HxTP is a secure developer-first CLI",
	Long:    `The official Hestia Labs Cross-Platform Trust Protocol CLI. Built for developers to add, control, and manage your devices instantly.`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&transportType, "transport", "rest", "Transport layer to use (rest, mqtt, ws)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("\n❌ %v\n", err)
		os.Exit(1)
	}
}
