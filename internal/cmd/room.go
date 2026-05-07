package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hestialabs/hxtp-cli/internal/auth"
	"github.com/hestialabs/hxtp-go/client"
	"github.com/spf13/cobra"
)

var roomHomeId string

var roomCmd = &cobra.Command{
	Use:     "room",
	Aliases: []string{"r"},
	Short:   "Manage rooms within a Smart Space",
}

var roomListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List rooms in a specific Smart Space",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, cfg, err := auth.RequireAuth()
		if err != nil {
			return err
		}

		if roomHomeId == "" && cfg.ActiveSpaceID != "" {
			roomHomeId = cfg.ActiveSpaceID
		}

		if roomHomeId == "" {
			return fmt.Errorf("Smart Space ID is required. Use --space [id] or 'space use [id]'")
		}

		hxtpClient := client.NewClient(client.ClientConfig{
			BaseURL: cfg.ApiUrl,
			Token:   token,
		})
		res, err := hxtpClient.ListRooms(roomHomeId)
		if err != nil {
			return err
		}

		rooms, ok := res["rooms"].([]interface{})
		if !ok || len(rooms) == 0 {
			fmt.Println("No rooms found in this space.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME")
		for _, r := range rooms {
			rMap := r.(map[string]interface{})
			fmt.Fprintf(w, "%s\t%s\n", rMap["id"], rMap["name"])
		}
		w.Flush()
		
		fmt.Printf("\nTotal: %d rooms\n", len(rooms))
		return nil
	},
}

var roomAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new room to a Smart Space",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, cfg, err := auth.RequireAuth()
		if err != nil {
			return err
		}

		if roomHomeId == "" && cfg.ActiveSpaceID != "" {
			roomHomeId = cfg.ActiveSpaceID
		}

		if roomHomeId == "" {
			return fmt.Errorf("Smart Space ID is required. Use --space [id] or 'space use [id]'")
		}
		name := args[0]

		hxtpClient := client.NewClient(client.ClientConfig{
			BaseURL: cfg.ApiUrl,
			Token:   token,
		})
		res, err := hxtpClient.CreateRoom(roomHomeId, name)
		if err != nil {
			return err
		}

		roomId := "unknown"
		if id, ok := res["id"].(string); ok {
			roomId = id
		}

		fmt.Printf("✅ Room '%s' successfully added to space.\n", name)
		fmt.Printf("ID: %s\n", roomId)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(roomCmd)
	roomCmd.AddCommand(roomListCmd)
	roomCmd.AddCommand(roomAddCmd)

	roomCmd.PersistentFlags().StringVarP(&roomHomeId, "space", "s", "", "ID of the target Smart Space (defaults to 'space use' setting)")
}
