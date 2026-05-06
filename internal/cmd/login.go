package cmd

import (
	"github.com/hestialabs/hxtp-cli/internal/ui/views"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with HxTP Cloud",
	Long:  `Starts the interactive TUI login wizard to secure your credentials and link the CLI to Hestia Cloud.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return views.LoginFlow()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
