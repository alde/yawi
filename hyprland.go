package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type HyprlandWindow struct {
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

func GetHyprlandActiveWindow() (*HyprlandWindow, error) {

	hyprSignature := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if hyprSignature == "" {
		return nil, fmt.Errorf("HYPRLAND_INSTANCE_SIGNATURE not found")
	}

	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}

	socketPath := filepath.Join(runtimeDir, "hypr", hyprSignature, ".socket2.sock")

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Hyprland: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("activewindow"))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read Hyprland's response: %w", err)
	}

	response := strings.TrimSpace(string(buffer[:n]))

	if response == "Invalid" || response == "" {
		return nil, fmt.Errorf("no active window found in Hyprland")
	}

	var window HyprlandWindow
	if err := json.Unmarshal([]byte(response), &window); err != nil {
		return nil, fmt.Errorf("failed to decode Hyprland's JSON: %w", err)
	}

	return &window, nil
}
