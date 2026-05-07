package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hestialabs/hxtp-cli/internal/auth"
	"github.com/hestialabs/hxtp-cli/internal/ui/views"
	"github.com/hestialabs/hxtp-go/client"
	"github.com/spf13/cobra"
)

var deviceType string
var homeId string

var deviceCmd = &cobra.Command{
	Use:     "device",
	Aliases: []string{"d", "gear"},
	Short:   "Manage your smart hardware (Gear)",
}

var deviceListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all your connected devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, cfg, err := auth.RequireAuth()
		if err != nil {
			return err
		}

		hxtpClient := client.NewClient(client.ClientConfig{
			BaseURL: cfg.ApiUrl,
			Token:   token,
		})
		res, err := hxtpClient.ListDevices()
		if err != nil {
			return err
		}

		devices, ok := res["devices"].([]interface{})
		if !ok || len(devices) == 0 {
			fmt.Println("No devices found. Add one using 'device add'.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tSTATUS\tROOM")
		for _, d := range devices {
			dMap := d.(map[string]interface{})
			status := "online"
			if active, ok := dMap["active"].(bool); ok && !active {
				status = "offline"
			}
			room, _ := dMap["room_id"].(string)
			if room == "" {
				room = "unassigned"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", dMap["id"], dMap["device_type"], status, room)
		}
		w.Flush()
		return nil
	},
}

var deviceCreateCmd = &cobra.Command{
	Use:   "add",
	Short: "Connect new hardware to your space",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, cfg, _ := auth.RequireAuth()
		if homeId == "" && cfg != nil && cfg.ActiveSpaceID != "" {
			homeId = cfg.ActiveSpaceID
		}
		return views.DeviceCreateFlow(homeId, deviceType)
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceCreateCmd)

	deviceCreateCmd.Flags().StringVarP(&deviceType, "type", "t", "", "Type of hardware (e.g., smart_switch)")
	deviceCreateCmd.Flags().StringVarP(&homeId, "space", "s", "", "ID of the target Smart Space")
}
