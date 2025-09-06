package providers

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/alde/yawi/pkg/window"
)

// HyprlandProvider implements window information retrieval for Hyprland
type HyprlandProvider struct{}

// Name returns the provider name
func (h *HyprlandProvider) Name() string {
	return "Hyprland"
}

// hyprlandWindow represents the JSON structure returned by Hyprland's activewindow command
type hyprlandWindow struct {
	Address   string `json:"address"`
	Mapped    bool   `json:"mapped"`
	Hidden    bool   `json:"hidden"`
	At        [2]int `json:"at"`
	Size      [2]int `json:"size"`
	Workspace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"workspace"`
	Floating       bool   `json:"floating"`
	Monitor        int    `json:"monitor"`
	Class          string `json:"class"`
	Title          string `json:"title"`
	InitialClass   string `json:"initialClass"`
	InitialTitle   string `json:"initialTitle"`
	PID            int    `json:"pid"`
	XWayland       bool   `json:"xwayland"`
	Pinned         bool   `json:"pinned"`
	Fullscreen     bool   `json:"fullscreen"`
	FakeFullscreen bool   `json:"fakeFullscreen"`
}

// GetActiveWindow retrieves the currently active window from Hyprland
func (h *HyprlandProvider) GetActiveWindow() (*window.WindowInfo, error) {
	hyprSignature := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if hyprSignature == "" {
		return nil, fmt.Errorf("HYPRLAND_INSTANCE_SIGNATURE not found - are we really running under Hyprland?")
	}

	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}

	socketPath := filepath.Join(runtimeDir, "hypr", hyprSignature, ".socket2.sock")

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Hyprland socket: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("activewindow"))
	if err != nil {
		return nil, fmt.Errorf("failed to send activewindow request: %w", err)
	}

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read Hyprland response: %w", err)
	}

	response := strings.TrimSpace(string(buffer[:n]))

	if response == "Invalid" || response == "" {
		return nil, fmt.Errorf("no active window found in Hyprland")
	}

	var hyprWindow hyprlandWindow
	if err := json.Unmarshal([]byte(response), &hyprWindow); err != nil {
		return nil, fmt.Errorf("failed to decode Hyprland JSON response: %w", err)
	}

	// Use workspace name if available, otherwise fall back to ID
	workspaceName := hyprWindow.Workspace.Name
	if workspaceName == "" {
		workspaceName = fmt.Sprintf("%d", hyprWindow.Workspace.ID)
	}

	return &window.WindowInfo{
		Title:     hyprWindow.Title,
		Class:     hyprWindow.Class,
		PID:       hyprWindow.PID,
		Workspace: workspaceName,
	}, nil
}