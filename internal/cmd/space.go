package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hestialabs/hxtp-cli/internal/auth"
	"github.com/hestialabs/hxtp-go/client"
	"github.com/spf13/cobra"
)

var spaceCmd = &cobra.Command{
	Use:     "space",
	Aliases: []string{"s", "smartspace"},
	Short:   "Manage your Smart Spaces",
}

var spaceListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all your registered Smart Spaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, cfg, err := auth.RequireAuth()
		if err != nil {
			return err
		}

		hxtpClient := client.NewClient(client.ClientConfig{
			BaseURL: cfg.ApiUrl,
			Token:   token,
		})
		res, err := hxtpClient.ListHomes()
		if err != nil {
			return err
		}

		homes, ok := res["homes"].([]interface{})
		if !ok || len(homes) == 0 {
			fmt.Println("No Smart Spaces found. Create one in the dashboard or via 'space add'.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "CURRENT\tID\tNAME\tTIMEZONE")
		for _, h := range homes {
			hMap := h.(map[string]interface{})
			id := hMap["id"].(string)
			
			activeMarker := ""
			if id == cfg.ActiveSpaceID {
				activeMarker = "*"
			}

			name, _ := hMap["home_name"].(string)
			if name == "" {
				name, _ = hMap["name"].(string)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", activeMarker, id, name, hMap["timezone"])
		}
		w.Flush()
		return nil
	},
}

var spaceUseCmd = &cobra.Command{
	Use:   "use [id]",
	Short: "Set the default Smart Space for subsequent commands",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, cfg, err := auth.RequireAuth()
		if err != nil {
			return err
		}

		id := args[0]
		cfg.ActiveSpaceID = id

		if err := auth.SaveConfig(cfg); err != nil {
			return err
		}

		fmt.Printf("✅ Active Smart Space set to: %s\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(spaceCmd)
	spaceCmd.AddCommand(spaceListCmd)
	spaceCmd.AddCommand(spaceUseCmd)
}
