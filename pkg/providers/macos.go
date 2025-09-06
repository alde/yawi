package providers

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/alde/yawi/pkg/window"
)

// MacOSProvider implements window information retrieval for macOS
type MacOSProvider struct{}

// Name returns the provider name
func (m *MacOSProvider) Name() string {
	return "macOS"
}

// GetActiveWindow retrieves the currently active window from macOS
func (m *MacOSProvider) GetActiveWindow() (*window.WindowInfo, error) {
	// Try the simple approach first - just get the frontmost app
	windowInfo, err := m.getFrontmostApp()
	if err != nil {
		// Fallback to lsappinfo if the simple approach fails
		return m.getActiveWindowLSAppInfo()
	}
	return windowInfo, nil
}

// getFrontmostApp gets the frontmost application without requiring accessibility permissions
func (m *MacOSProvider) getFrontmostApp() (*window.WindowInfo, error) {
	// Simple AppleScript that doesn't require accessibility permissions
	script := `
tell application "System Events"
	set frontApp to first application process whose frontmost is true
	set appName to name of frontApp
	set appPID to unix id of frontApp
	return appName & "|" & (appPID as string)
end tell`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute AppleScript: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return nil, fmt.Errorf("no active application found")
	}

	// Parse the result: appName|pid
	parts := strings.Split(result, "|")
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected AppleScript output format")
	}

	appName := parts[0]
	pidStr := parts[1]
	
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		pid = 0 // Default if parsing fails
	}

	return &window.WindowInfo{
		Title:     appName,         // Use app name as title
		Class:     appName,         // On macOS, the app name serves as the "class"
		PID:       pid,
		Workspace: "main",         // Default workspace name
	}, nil
}


// Alternative implementation using lsappinfo (if AppleScript fails)
func (m *MacOSProvider) getActiveWindowLSAppInfo() (*window.WindowInfo, error) {
	// Get the frontmost app using lsappinfo
	cmd := exec.Command("/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsappinfo", 
		"info", "-only", "name,pid", "-app", "front")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute lsappinfo: %w", err)
	}

	result := string(output)
	
	// Parse the output to extract app name and PID
	nameRegex := regexp.MustCompile(`"([^"]+)"`)
	pidRegex := regexp.MustCompile(`pid=(\d+)`)
	
	nameMatch := nameRegex.FindStringSubmatch(result)
	pidMatch := pidRegex.FindStringSubmatch(result)
	
	if len(nameMatch) < 2 {
		return nil, fmt.Errorf("could not parse app name from lsappinfo output")
	}
	
	appName := nameMatch[1]
	pid := 0
	
	if len(pidMatch) >= 2 {
		if p, err := strconv.Atoi(pidMatch[1]); err == nil {
			pid = p
		}
	}

	return &window.WindowInfo{
		Title:     appName, // For lsappinfo, we only get the app name
		Class:     appName,
		PID:       pid,
		Workspace: "main", // Default workspace
	}, nil
}