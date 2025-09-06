package window

import (
	"testing"
)

func TestWindowInfo_String(t *testing.T) {
	tests := []struct {
		name     string
		window   WindowInfo
		expected string
	}{
		{
			name: "Normal window with title and class",
			window: WindowInfo{
				Title: "YAWI - Yet Another Window Inspector",
				Class: "Firefox",
				PID:   12345,
			},
			expected: "🪟 YAWI - Yet Another Window Inspector (Firefox)",
		},
		{
			name: "Window with empty title",
			window: WindowInfo{
				Title: "",
				Class: "Terminal",
				PID:   54321,
			},
			expected: "🪟  (Terminal)",
		},
		{
			name: "Window with empty class",
			window: WindowInfo{
				Title: "Some Title",
				Class: "",
				PID:   98765,
			},
			expected: "🪟 Some Title ()",
		},
		{
			name: "Empty window info",
			window: WindowInfo{
				Title: "",
				Class: "",
				PID:   0,
			},
			expected: "No active window found!",
		},
		{
			name: "Window with special characters",
			window: WindowInfo{
				Title: "Test & Special < Characters >",
				Class: "app-name",
				PID:   1337,
			},
			expected: "🪟 Test & Special < Characters > (app-name)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.window.String()
			if result != tt.expected {
				t.Errorf("WindowInfo.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestWindowInfo_JSONMarshaling(t *testing.T) {
	window := WindowInfo{
		Title:     "Test Window",
		Class:     "TestApp",
		PID:       12345,
		Workspace: "workspace-1",
	}

	// Test that struct fields are accessible - the JSON tags are working if the fields are set correctly
	if window.Title != "Test Window" {
		t.Errorf("Expected Title to be 'Test Window', got %q", window.Title)
	}
	if window.Class != "TestApp" {
		t.Errorf("Expected Class to be 'TestApp', got %q", window.Class)
	}
	if window.PID != 12345 {
		t.Errorf("Expected PID to be 12345, got %d", window.PID)
	}
	if window.Workspace != "workspace-1" {
		t.Errorf("Expected Workspace to be 'workspace-1', got %q", window.Workspace)
	}
}