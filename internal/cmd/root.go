package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

var transportType string

var rootCmd = &cobra.Command{
	Use:     "hxtp-cli",
	Version: "v1.0.1",
	Short:   "HxTP is a secure developer-first CLI",
	Long:    `The official Hestia Labs Cross-Platform Trust Protocol CLI. Built for developers to add, control, and manage your devices instantly.`,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the CLI to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Checking for updates...")
		// Use the official install script logic
		updateCmd := "curl -fsSL https://hestialabs.in/install.sh | bash"
		c := exec.Command("bash", "-c", updateCmd)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	rootCmd.PersistentFlags().StringVar(&transportType, "transport", "rest", "Transport layer to use (rest, mqtt, ws)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("\n❌ %v\n", err)
		os.Exit(1)
	}
}
