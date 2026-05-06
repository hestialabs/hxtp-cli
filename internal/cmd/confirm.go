package cmd

import (
	"fmt"

	"github.com/hestialabs/hxtp-cli/internal/auth"
	"github.com/hestialabs/hxtp-cli/internal/ui"
	"github.com/hestialabs/hxtp-go/client"
	"github.com/spf13/cobra"
)

var confirmCmd = &cobra.Command{
	Use:   "confirm [device_id] [token]",
	Short: "Confirm a critical gateway command",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		theme := ui.GetTheme()
		deviceId := args[0]
		token := args[1]

		authToken, cfg, err := auth.RequireAuth()
		if err != nil {
			return err
		}

		hxtpClient := client.NewClient(client.ClientConfig{
			BaseURL: cfg.ApiUrl,
			Token:   authToken,
		})

		fmt.Print("🔐 Validating safety token... ")

		// Using SDK's baked-in ConfirmCommand
		_, err = hxtpClient.ConfirmCommand(deviceId, token)
		if err != nil {
			fmt.Println("❌")
			return err
		}

		fmt.Println("✅")
		fmt.Printf("%s\n", theme.SuccessMsg.Render("Safety gate passed. Action released to hardware."))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(confirmCmd)
}
