package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/hestialabs/hxtp-cli/internal/auth"
	"github.com/hestialabs/hxtp-cli/internal/ui"
	"github.com/hestialabs/hxtp-go/client"
)

// LoginFlow executes the interactive 2026 TUI login wizard.
func LoginFlow() error {
	theme := ui.GetTheme()

	var (
		loginMethod string
		apiUrl      string = "https://api.hestialabs.in/api/v1" // Default Production
		token       string
	)

	fmt.Println(theme.BoldHeader.Render("Welcome to HxTP 3.0"))
	fmt.Println(theme.SubHeader.Render("Log in to your Hestia Cloud account."))
	fmt.Println()

	// ── 1. Select Login Method ────────────────────
	methodForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How do you want to authenticate?").
				Options(
					huh.NewOption("Browser Login", "browser"),
					huh.NewOption("Manual API Token", "manual"),
				).
				Value(&loginMethod),
		),
	).WithTheme(huh.ThemeCharm())

	err := methodForm.Run()
	if err != nil {
		return err
	}

	if loginMethod == "browser" {
		// ── 2. Seamless Browser Flow ──────────────────
		token, err = auth.LoginWithBrowser(apiUrl)
		if err != nil {
			return err
		}
	} else {
		// ── 3. Manual Token Flow ──────────────────────
		fmt.Println(theme.SubHeader.Render("Enter your secure Gateway URL and Token from the dashboard."))
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("API Gateway URL").
					Value(&apiUrl).
					Placeholder("https://api.hestialabs.in/api/v1"),
				huh.NewInput().
					Title("Secure API Token").
					Value(&token).
					Password(true).
					Placeholder("hxtp_..."),
			),
		).WithTheme(huh.ThemeCharm())

		err = form.Run()
		if err != nil {
			return err
		}
	}

	// ── 2. Unified Verification with Pulse ──────
	hxtpClient := client.NewClient(client.ClientConfig{
		BaseURL: apiUrl,
		Token:   token,
	})

	err = ui.Pulse("Verifying credentials...", func() error {
		// We'll use a simple Check (listing devices) to verify the token
		_, pulseErr := hxtpClient.SendCommand("ping", "status", nil, false)
		if pulseErr != nil && !strings.Contains(pulseErr.Error(), "404") {
			return pulseErr
		}
		return nil
	})

	if err != nil {
		fmt.Println("❌")
		return fmt.Errorf("Authentication failed: %v", err)
	}

	// Save
	err = auth.SaveToken(token)
	if err != nil {
		return fmt.Errorf("Failed to save to keychain: %v", err)
	}

	err = auth.SaveConfig(&auth.Config{
		ApiUrl:    apiUrl,
		LastLogin: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	fmt.Printf("✅ %s\n", theme.SuccessMsg.Render("Credentials verified."))
	fmt.Println()
	fmt.Printf("%s\n", theme.SuccessMsg.Render("Successfully authenticated!"))
	fmt.Printf("%s credentials secured in system keychain.\n", theme.InfoMsg.Render("Encrypted"))

	return nil
}
