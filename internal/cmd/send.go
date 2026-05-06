package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hestialabs/hxtp-cli/internal/auth"
	"github.com/hestialabs/hxtp-cli/internal/ui"
	"github.com/hestialabs/hxtp-go/client"
	"github.com/hestialabs/hxtp-go/transport"
	"github.com/spf13/cobra"
)

var params []string
var dryRun bool

var sendCmd = &cobra.Command{
	Use:   "send [device_id] [action]",
	Short: "Send a command to a hardware device",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		theme := ui.GetTheme()
		deviceId := args[0]
		action := args[1]

		token, cfg, err := auth.RequireAuth()
		if err != nil {
			return err
		}

		paramMap := make(map[string]interface{})
		for _, p := range params {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) == 2 {
				paramMap[parts[0]] = parts[1]
			}
		}

		hxtpClient := client.NewClient(client.ClientConfig{
			BaseURL:  cfg.ApiUrl,
			Token:    token,
			DeviceId: cfg.DeviceId,
			ClientId: cfg.ClientId,
			Secret:   cfg.Secret,
		})

		// Setup Transport if native requested
		if transportType == "mqtt" {
			broker := "tcp://mqtt.hestialabs.in:1883" 
			if envBroker := os.Getenv("HXTP_BROKER"); envBroker != "" {
				broker = envBroker
			}
			
			mqttTransport := transport.NewMQTTTransport(broker, "hxtpctl-"+cfg.ClientId, "", "")
			err = mqttTransport.Connect(cmd.Context())
			if err != nil {
				return fmt.Errorf("MQTT_CONNECT_ERROR: %w", err)
			}
			defer mqttTransport.Disconnect(cmd.Context())
			
			hxtpClient.SetTransport(mqttTransport)
		}

		fmt.Printf("📡 Sending command %s to device %s via %s... ", 
			lipgloss.NewStyle().Foreground(ui.GetTheme().Accent).Render(action),
			deviceId,
			lipgloss.NewStyle().Foreground(ui.GetTheme().Secondary).Render(transportType),
		)

		resp, err := hxtpClient.SendCommand(deviceId, action, paramMap, dryRun)
		if err != nil {
			fmt.Println("❌")
			return err
		}

		fmt.Println("✅")

		if status, ok := resp["status"].(string); ok && status == "dry_run_required" {
			dryRunToken := resp["dry_run_token"].(string)
			fmt.Println()
			fmt.Println(theme.WarningMsg.Render("⚠️  SAFETY GATE TRIGGERED"))
			fmt.Printf("This action requires a second confirmation. Dry-run results look stable.\n")
			fmt.Printf("To finalize this action, run:\n\n")
			fmt.Printf("   %s\n\n", fmt.Sprintf("hxtp confirm %s --device-id %s", dryRunToken, deviceId))
		} else {
			fmt.Printf("%s\n", theme.SuccessMsg.Render("Command successfully sent to device!"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringSliceVarP(&params, "param", "p", []string{}, "Key-value parameters (key=value)")
	sendCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview action without execution")
}
