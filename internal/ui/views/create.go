package views

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/hestialabs/hxtp-cli/internal/auth"
	"github.com/hestialabs/hxtp-cli/internal/ui"
	"github.com/hestialabs/hxtp-go/client"
)

// DeviceCreateFlow executes the interactive device registration wizard.
func DeviceCreateFlow(initialHomeId, initialDeviceType string) error {
	theme := ui.GetTheme()
	token, cfg, err := auth.RequireAuth()
	if err != nil {
		return err
	}

	hxtpClient := client.NewClient(client.ClientConfig{
		BaseURL: cfg.ApiUrl,
		Token:   token,
	})

	var (
		homeId     = initialHomeId
		deviceType = initialDeviceType
	)

	fmt.Println(theme.BoldHeader.Render("Add New Device"))
	fmt.Println(theme.SubHeader.Render("Connecting a new device to your home."))
	fmt.Println()

	// ── 1. Discover Homes with Pulse Animation ────
	var homeResult map[string]interface{}
	err = ui.Pulse("Discovering available homes...", func() error {
		var pulseErr error
		homeResult, pulseErr = hxtpClient.ListHomes()
		return pulseErr
	})

	if err != nil {
		fmt.Println("❌")
		return fmt.Errorf("Discovery failed: %v", err)
	}
	fmt.Printf("✅ %s\n", theme.SuccessMsg.Render("Home discovery complete."))

	homesList, ok := homeResult["homes"].([]interface{})
	if !ok || len(homesList) == 0 {
		return fmt.Errorf("No homes found. Please create a home in the dashboard first.")
	}

	options := make([]huh.Option[string], len(homesList))
	for i, h := range homesList {
		hMap := h.(map[string]interface{})
		// Use home_name to match backend
		name, _ := hMap["home_name"].(string)
		if name == "" {
			name, _ = hMap["name"].(string) // Fallback
		}
		id := hMap["id"].(string)
		options[i] = huh.NewOption(name, id)
	}

	// Skip wizard if flags are provided
	if homeId != "" && deviceType != "" {
		fmt.Printf("⏩ %s\n", theme.InfoMsg.Render("Flags provided. Bypassing interactive wizard."))
	} else {
		// ── 2. The Wizard ──────────────────────────────
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select Target Home").
					Options(options...).
					Value(&homeId),

				huh.NewInput().
					Title("Device Type").
					Value(&deviceType).
					Placeholder("e.g. gateway, light, sensor").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("Device type is required")
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeCharm())

		err = form.Run()
		if err != nil {
			return err
		}
	}

	// ── 3. Register with Pulse Animation ──────────
	var regResult map[string]interface{}
	err = ui.Pulse("Registering device...", func() error {
		var pulseErr error
		regResult, pulseErr = hxtpClient.RegisterDevice(deviceType, homeId, nil)
		return pulseErr
	})

	if err != nil {
		fmt.Println("❌")
		return fmt.Errorf("Registration failed: %v", err)
	}
	fmt.Printf("✅ %s\n", theme.SuccessMsg.Render("Device registered successfully."))

	// ── 4. Premium Result ─────────────────────────
	deviceId := regResult["device_id"].(string)
	secret := regResult["device_secret"].(string)

	fmt.Println()
	fmt.Println(theme.SuccessMsg.Render("Device Successfully Added!"))
	fmt.Printf("ID:     %s\n", lipgloss.NewStyle().Foreground(theme.Accent).Render(deviceId))
	fmt.Printf("Secret: %s\n", lipgloss.NewStyle().Foreground(theme.Accent).Render(secret))
	fmt.Println()
	fmt.Println(theme.InfoMsg.Render("IMPORTANT: Store this secret securely. It will never be shown again."))
	
	if f, ok := regResult["api_base_url"].(string); ok {
		fmt.Printf("Endpoint: %s\n", f)
	}

	return nil
}
