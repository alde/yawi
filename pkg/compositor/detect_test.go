package compositor

import (
	"os"
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	// Save original environment
	originalVars := map[string]string{
		"HYPRLAND_INSTANCE_SIGNATURE": os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"),
		"SWAYSOCK":                    os.Getenv("SWAYSOCK"),
		"XDG_CURRENT_DESKTOP":         os.Getenv("XDG_CURRENT_DESKTOP"),
		"XDG_SESSION_DESKTOP":         os.Getenv("XDG_SESSION_DESKTOP"),
	}

	// Clean up after test
	defer func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected Type
	}{
		{
			name: "Hyprland detection",
			envVars: map[string]string{
				"HYPRLAND_INSTANCE_SIGNATURE": "some-signature",
				"SWAYSOCK":                    "",
				"XDG_CURRENT_DESKTOP":         "",
				"XDG_SESSION_DESKTOP":         "",
			},
			expected: Hyprland,
		},
		{
			name: "Sway detection",
			envVars: map[string]string{
				"HYPRLAND_INSTANCE_SIGNATURE": "",
				"SWAYSOCK":                    "/run/user/1000/sway-ipc.sock",
				"XDG_CURRENT_DESKTOP":         "",
				"XDG_SESSION_DESKTOP":         "",
			},
			expected: Sway,
		},
		{
			name: "GNOME detection via XDG_CURRENT_DESKTOP",
			envVars: map[string]string{
				"HYPRLAND_INSTANCE_SIGNATURE": "",
				"SWAYSOCK":                    "",
				"XDG_CURRENT_DESKTOP":         "GNOME",
				"XDG_SESSION_DESKTOP":         "",
			},
			expected: GNOME,
		},
		{
			name: "GNOME detection via XDG_SESSION_DESKTOP",
			envVars: map[string]string{
				"HYPRLAND_INSTANCE_SIGNATURE": "",
				"SWAYSOCK":                    "",
				"XDG_CURRENT_DESKTOP":         "",
				"XDG_SESSION_DESKTOP":         "gnome",
			},
			expected: GNOME,
		},
		{
			name: "Unknown when no match",
			envVars: map[string]string{
				"HYPRLAND_INSTANCE_SIGNATURE": "",
				"SWAYSOCK":                    "",
				"XDG_CURRENT_DESKTOP":         "unity",
				"XDG_SESSION_DESKTOP":         "",
			},
			expected: func() Type {
				if runtime.GOOS == "darwin" {
					return MacOS
				}
				return Unknown
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip non-macOS tests on macOS since it always returns MacOS
			if runtime.GOOS == "darwin" && tt.expected != MacOS {
				t.Skip("Skipping non-macOS test on macOS")
			}

			// Set environment variables
			for key, value := range tt.envVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}

			result := Detect()
			if result != tt.expected {
				t.Errorf("Detect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestType_String(t *testing.T) {
	tests := []struct {
		input    Type
		expected string
	}{
		{Hyprland, "Hyprland"},
		{Sway, "Sway"},
		{GNOME, "GNOME"},
		{MacOS, "macOS"},
		{Unknown, "Unknown"},
		{Type(999), "Unknown"}, // Invalid type
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("Type.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMacOSDetection(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS detection test only runs on macOS")
	}

	// On macOS, it should always return MacOS regardless of other env vars
	os.Setenv("XDG_CURRENT_DESKTOP", "GNOME")
	os.Setenv("SWAYSOCK", "/some/socket")
	defer func() {
		os.Unsetenv("XDG_CURRENT_DESKTOP")
		os.Unsetenv("SWAYSOCK")
	}()

	result := Detect()
	if result != MacOS {
		t.Errorf("On macOS, Detect() should return MacOS, got %v", result)
	}
}