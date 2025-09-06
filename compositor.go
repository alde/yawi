package main

import (
	"os"
	"strings"
)

type Compositor int

const (
	Unknown Compositor = iota
	Hyprland
	Sway
	GNOME
)

func (c Compositor) String() string {
	switch c {
	case Hyprland:
		return "Hyprland"
	case Sway:
		return "Sway"
	case GNOME:
		return "GNOME"
	default:
		return "Unknown"
	}
}

func DetectCompositor() Compositor {

	if hyprInstance := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"); hyprInstance != "" {
		return Hyprland
	}

	if swaySocket := os.Getenv("SWAYSOCK"); swaySocket != "" {
		return Sway
	}

	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	session := strings.ToLower(os.Getenv("XDG_SESSION_DESKTOP"))

	if strings.Contains(desktop, "gnome") || strings.Contains(session, "gnome") {
		return GNOME
	}

	return Unknown
}
