package providers

import (
	"fmt"

	"yawi/pkg/compositor"
	"yawi/pkg/window"
)

// NewProvider creates a window provider for the given compositor/window manager type
func NewProvider(comp compositor.Type) (window.Provider, error) {
	switch comp {
	case compositor.Hyprland:
		return &HyprlandProvider{}, nil
	case compositor.Sway:
		return &SwayProvider{}, nil
	case compositor.GNOME:
		return &GNOMEProvider{}, nil
	case compositor.MacOS:
		return &MacOSProvider{}, nil
	default:
		return nil, fmt.Errorf("unsupported compositor: %s\nSupported: Hyprland, Sway, GNOME Shell, macOS", comp)
	}
}