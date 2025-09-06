package window

import "fmt"

// WindowInfo represents information about a window across different compositors
type WindowInfo struct {
	Title     string `json:"title"`
	Class     string `json:"class"`
	PID       int    `json:"pid"`
	Workspace string `json:"workspace"`
}

// String returns a friendly string representation of the window
func (w WindowInfo) String() string {
	if w.Title == "" && w.Class == "" {
		return "No active window found!"
	}
	return fmt.Sprintf("🪟 %s (%s)", w.Title, w.Class)
}

// Provider defines the interface for getting window information from different compositors
type Provider interface {
	// GetActiveWindow returns information about the currently active window
	GetActiveWindow() (*WindowInfo, error)
	
	// Name returns the human-readable name of this provider
	Name() string
}