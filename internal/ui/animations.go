package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// PixelSpinner defines the "Cybernetic Pulse" motion frames for 2026.
// Optimized for high-bitrate TUI rendering with a "Liquid Braille" effect.
var PixelSpinnerSet = spinner.Spinner{
	Frames: []string{
		"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
		"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷",
	},
	FPS: time.Second / 30, // Ultra-high refresh for "Hestia-class" performance
}

// GetPixelSpinner returns a pre-styled spinner with HxTP Cyber-Cyan accents.
func GetPixelSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = PixelSpinnerSet
	s.Style = lipgloss.NewStyle().Foreground(GetTheme().Accent)
	return s
}

// Pulse executes a function while showing a high-end pixel animation.
// It handles the animation loop in a background goroutine.
func Pulse(message string, fn func() error) error {
	s := GetPixelSpinner()
	done := make(chan bool)
	errChan := make(chan error)

	fmt.Print("\033[?25l") // Hide cursor
	defer fmt.Print("\033[?25h") // Show cursor

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// Clear line and redraw
				fmt.Printf("\r%s %s", s.View(), GetTheme().SubHeader.Render(message))
				s, _ = s.Update(spinner.TickMsg{})
				time.Sleep(s.Spinner.FPS)
			}
		}
	}()

	go func() {
		errChan <- fn()
		done <- true
	}()

	err := <-errChan
	fmt.Print("\r\033[K") // Clear the animation line
	return err
}
