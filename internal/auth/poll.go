package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/hestialabs/hxtp-cli/internal/ui"
)

// HandshakeResult defines the token returned from the browser-based auth bridge.
type HandshakeResult struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}

// LoginWithBrowser executes the Vercel-style "Device Handshake" flow.
func LoginWithBrowser(apiUrl string) (string, error) {
	// 1. Generate a unique, high-entropy handshake session (32 bytes = 256 bits)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("AUTH_ERROR: Failed to generate premium handshake ID: %v", err)
	}
	handshakeId := fmt.Sprintf("HXT-%s", hex.EncodeToString(bytes))
	loginUrl := fmt.Sprintf("%s/auth/cli?code=%s", apiUrl, handshakeId)

	// Create OSC 8 clickable link for modern terminals
	clickableUrl := fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", loginUrl, loginUrl)
	maskedId := fmt.Sprintf("HXT-%s", "****************************************************************")

	fmt.Println(ui.GetTheme().BoldHeader.Render("HxTP Secure Handshake"))
	fmt.Printf("Action Required: %s\n", ui.GetTheme().SubHeader.Render("Authorize this device in your browser"))
	fmt.Printf("Security URL: %s\n", ui.GetTheme().AccentMsg.Render(clickableUrl))
	fmt.Printf("Verification Code: %s\n\n", ui.GetTheme().AccentMsg.Render(maskedId))

	// 2. Launch Browser (Background)
	_ = browser.OpenURL(loginUrl)

	// 3. Poll for the token using our premium Pulse UX
	var token string
	err := ui.Pulse("Waiting for browser authorization...", func() error {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		timeout := time.After(5 * time.Minute)

		for {
			select {
			case <-timeout:
				return fmt.Errorf("AUTH_TIMEOUT: Handshake session expired")
			case <-ticker.C:
				// Polling Endpoint (Assumed backend route: GET /auth/handshake/:id)
				pollUrl := fmt.Sprintf("%s/auth/handshake/%s", apiUrl, handshakeId)
				resp, err := http.Get(pollUrl)
				if err != nil {
					continue // Network jitter
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					body, _ := io.ReadAll(resp.Body)
					var res HandshakeResult
					json.Unmarshal(body, &res)
					if res.Token != "" {
						token = res.Token
						return nil
					}
				}
			}
		}
	})

	return token, err
}
