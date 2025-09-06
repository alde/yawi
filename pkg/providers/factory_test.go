package providers

import (
	"strings"
	"testing"

	"github.com/alde/yawi/pkg/compositor"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name          string
		compositorType compositor.Type
		expectError   bool
		expectedType  string
	}{
		{
			name:          "Hyprland provider",
			compositorType: compositor.Hyprland,
			expectError:   false,
			expectedType:  "*providers.HyprlandProvider",
		},
		{
			name:          "Sway provider",
			compositorType: compositor.Sway,
			expectError:   false,
			expectedType:  "*providers.SwayProvider",
		},
		{
			name:          "GNOME provider",
			compositorType: compositor.GNOME,
			expectError:   false,
			expectedType:  "*providers.GNOMEProvider",
		},
		{
			name:          "macOS provider",
			compositorType: compositor.MacOS,
			expectError:   false,
			expectedType:  "*providers.MacOSProvider",
		},
		{
			name:          "Unknown compositor",
			compositorType: compositor.Unknown,
			expectError:   true,
			expectedType:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(tt.compositorType)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for compositor type %v, but got none", tt.compositorType)
				}
				if provider != nil {
					t.Errorf("Expected nil provider for unknown compositor, got %T", provider)
				}
				// Check error message mentions supported platforms
				if !strings.Contains(err.Error(), "Supported") {
					t.Errorf("Error message should mention supported platforms, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for compositor type %v: %v", tt.compositorType, err)
				}
				if provider == nil {
					t.Errorf("Expected provider for compositor type %v, got nil", tt.compositorType)
				}
				// Just verify provider exists and has a name
				if provider != nil && provider.Name() == "" {
					t.Errorf("Provider should have a non-empty name")
				}
			}
		})
	}
}

func TestProviderNames(t *testing.T) {
	// Test that provider names are consistent
	tests := []struct {
		compositorType compositor.Type
		expectedName   string
	}{
		{compositor.Hyprland, "Hyprland"},
		{compositor.Sway, "Sway"},
		{compositor.GNOME, "GNOME Shell"},
		{compositor.MacOS, "macOS"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedName, func(t *testing.T) {
			provider, err := NewProvider(tt.compositorType)
			if err != nil {
				t.Fatalf("Unexpected error creating provider: %v", err)
			}

			name := provider.Name()
			if name != tt.expectedName {
				t.Errorf("Provider.Name() = %q, want %q", name, tt.expectedName)
			}
		})
	}
}