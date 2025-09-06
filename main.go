package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type WindowInfo struct {
	Title     string `json:"title"`
	Class     string `json:"class"`
	PID       int    `json:"pid"`
	Workspace string `json:"workspace"`
}

func (w WindowInfo) String() string {
	if w.Title == "" && w.Class == "" {
		return "No active window found!"
	}
	return fmt.Sprintf("🪟 %s (%s)", w.Title, w.Class)
}

func main() {
	var (
		showFull       = flag.Bool("full", false, "Output full data")
		showHelp       = flag.Bool("help", false, "Show help")
		showCompositor = flag.Bool("compositor", false, "Show which compositor is detected")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "YAWI - Yet Another Window Inspector!\n")
		fmt.Fprintf(os.Stderr, "A simple tool to get the currently active window's class in Wayland\n\n")
		fmt.Fprintf(os.Stderr, "Usage: yawi [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSupported Compositors: Hyprland, Sway, GNOME\n")
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	compositor := DetectCompositor()

	if *showCompositor {
		fmt.Printf("Current compositor: %s\n", compositor)
		return
	}

	windowInfo, err := getActiveWindowInfo(compositor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	if *showFull {
		jsonData, err := json.MarshalIndent(windowInfo, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to create JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Println(windowInfo.Class)
	}
}

func getActiveWindowInfo(compositor Compositor) (*WindowInfo, error) {
	switch compositor {
	case Hyprland:
		return getHyprlandWindowInfo()
	case Sway:
		return getSwayWindowInfo()
	case GNOME:
		return GetGNOMEActiveWindow()
	default:
		return nil, fmt.Errorf("unsupported or unknown compositor: %s\nsupported: Hyprland, Sway, GNOME-Shell", compositor)
	}
}

func getHyprlandWindowInfo() (*WindowInfo, error) {
	window, err := GetHyprlandActiveWindow()
	if err != nil {
		return nil, err
	}

	workspaceName := window.Workspace.Name
	if workspaceName == "" {
		workspaceName = fmt.Sprintf("%d", window.Workspace.ID)
	}

	return &WindowInfo{
		Title:     window.Title,
		Class:     window.Class,
		PID:       window.PID,
		Workspace: workspaceName,
	}, nil
}

func getSwayWindowInfo() (*WindowInfo, error) {
	node, err := GetSwayActiveWindow()
	if err != nil {
		return nil, err
	}

	var title, class string
	var pid int

	if node.WindowProperties != nil {
		if node.WindowProperties.Title != nil {
			title = *node.WindowProperties.Title
		}
		if node.WindowProperties.Class != nil {
			class = *node.WindowProperties.Class
		}
	}

	if node.PID != nil {
		pid = *node.PID
	}

	workspace := ""
	if node.Representation != nil {
		workspace = *node.Representation
	} else if node.Name != nil {
		workspace = *node.Name
	}

	return &WindowInfo{
		Title:     title,
		Class:     class,
		PID:       pid,
		Workspace: workspace,
	}, nil
}
