package compositor

import (
	"os"
	"runtime"
	"strings"
)

// Type represents different window managers/compositors
type Type int

const (
	Unknown Type = iota
	Hyprland
	Sway
	GNOME
	MacOS
)

func (c Type) String() string {
	switch c {
	case Hyprland:
		return "Hyprland"
	case Sway:
		return "Sway"
	case GNOME:
		return "GNOME"
	case MacOS:
		return "macOS"
	default:
		return "Unknown"
	}
}

// Detect attempts to determine which window manager/compositor is currently running
// It checks environment variables and other indicators to make an educated guess
func Detect() Type {
	// Check if we're on macOS first
	if runtime.GOOS == "darwin" {
		return MacOS
	}

	// Hyprland sets a unique instance signature
	if hyprInstance := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"); hyprInstance != "" {
		return Hyprland
	}

	// Sway provides a socket path for IPC
	if swaySocket := os.Getenv("SWAYSOCK"); swaySocket != "" {
		return Sway
	}

	// GNOME can be detected through desktop environment variables
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	session := strings.ToLower(os.Getenv("XDG_SESSION_DESKTOP"))

	if strings.Contains(desktop, "gnome") || strings.Contains(session, "gnome") {
		return GNOME
	}

	return Unknown
}