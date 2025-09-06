package providers

import (
	"encoding/json"
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/alde/yawi/pkg/window"
)

// GNOMEProvider implements window information retrieval for GNOME Shell
type GNOMEProvider struct{}

// Name returns the provider name
func (g *GNOMEProvider) Name() string {
	return "GNOME Shell"
}

// focusedWindowInfo represents the structure returned by the GNOME FocusedWindow extension
type focusedWindowInfo struct {
	Title              string   `json:"title,omitempty"`
	WmClass            string   `json:"wm_class,omitempty"`
	WmClassInstance    string   `json:"wm_class_instance,omitempty"`
	Pid                int      `json:"pid,omitempty"`
	Id                 int      `json:"id,omitempty"`
	Width              int      `json:"width,omitempty"`
	Height             int      `json:"height,omitempty"`
	X                  int      `json:"x,omitempty"`
	Y                  int      `json:"y,omitempty"`
	Focus              bool     `json:"focus,omitempty"`
	InCurrentWorkspace bool     `json:"in_current_workspace,omitempty"`
	Moveable           bool     `json:"moveable,omitempty"`
	Resizable          bool     `json:"resizeable,omitempty"`
	CanClose           bool     `json:"canclose,omitempty"`
	CanMaximize        bool     `json:"canmaximize,omitempty"`
	Maximized          int      `json:"maximized,omitempty"`
	CanMinimize        bool     `json:"canminimize,omitempty"`
	Display            struct{} `json:"display,omitempty"`
	FrameType          int      `json:"frame_type,omitempty"`
	WindowType         int      `json:"window_type,omitempty"`
	Layer              int      `json:"layer,omitempty"`
	Monitor            int      `json:"monitor,omitempty"`
	Role               *string  `json:"role,omitempty"`
	Area               struct{} `json:"area,omitempty"`
	AreaAll            struct{} `json:"area_all,omitempty"`
	AreaCust           struct{} `json:"area_cust,omitempty"`
}

// GetActiveWindow retrieves the currently active window from GNOME Shell
func (g *GNOMEProvider) GetActiveWindow() (*window.WindowInfo, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to D-Bus session bus: %w", err)
	}
	defer conn.Close()

	// Try the FocusedWindow GNOME Shell extension
	windowInfo, err := g.tryFocusedWindowExtension(conn)
	if err == nil {
		return windowInfo, nil
	}

	return nil, fmt.Errorf("unable to get GNOME active window - make sure the Focused Window D-Bus extension is enabled")
}

// tryFocusedWindowExtension attempts to get window info via the FocusedWindow GNOME extension
func (g *GNOMEProvider) tryFocusedWindowExtension(conn *dbus.Conn) (*window.WindowInfo, error) {
	obj := conn.Object("org.gnome.Shell", "/org/gnome/shell/extensions/FocusedWindow")
	var result string
	err := obj.Call("org.gnome.shell.extensions.FocusedWindow.Get", 0).Store(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to call FocusedWindow.Get D-Bus method: %w", err)
	}

	var info focusedWindowInfo
	err = json.Unmarshal([]byte(result), &info)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal focused window info: %w", err)
	}

	return &window.WindowInfo{
		Title:     info.Title,
		Class:     info.WmClass,
		PID:       info.Pid,
		Workspace: fmt.Sprintf("%d", info.Id), // Using window ID as workspace for now
	}, nil
}