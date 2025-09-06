package providers

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"yawi/pkg/window"
)

// SwayProvider implements window information retrieval for Sway
type SwayProvider struct{}

// Name returns the provider name
func (s *SwayProvider) Name() string {
	return "Sway"
}

// swayNode represents the JSON structure of Sway's tree nodes
type swayNode struct {
	ID                 int        `json:"id"`
	Name               *string    `json:"name"`
	Type               string     `json:"type"`
	Border             string     `json:"border"`
	CurrentBorderWidth int        `json:"current_border_width"`
	Layout             string     `json:"layout"`
	Orientation        string     `json:"orientation"`
	Percent            *float64   `json:"percent"`
	Rect               swayRect   `json:"rect"`
	WindowRect         swayRect   `json:"window_rect"`
	DecoRect           swayRect   `json:"deco_rect"`
	Geometry           swayRect   `json:"geometry"`
	Urgent             bool       `json:"urgent"`
	Focused            bool       `json:"focused"`
	Focus              []int      `json:"focus"`
	Nodes              []*swayNode `json:"nodes"`
	FloatingNodes      []*swayNode `json:"floating_nodes"`
	Sticky             bool       `json:"sticky"`
	Representation     *string    `json:"representation"`
	AppID              *string    `json:"app_id"`
	WindowProperties   *struct {
		Class        *string `json:"class"`
		Instance     *string `json:"instance"`
		Title        *string `json:"title"`
		TransientFor *int    `json:"transient_for"`
	} `json:"window_properties"`
	PID *int `json:"pid"`
}

type swayRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// GetActiveWindow retrieves the currently active window from Sway
func (s *SwayProvider) GetActiveWindow() (*window.WindowInfo, error) {
	socketPath := os.Getenv("SWAYSOCK")
	if socketPath == "" {
		return nil, fmt.Errorf("SWAYSOCK environment variable not found - are we running under Sway?")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Sway socket: %w", err)
	}
	defer conn.Close()

	// Send i3-ipc GET_TREE message (magic bytes + payload length + message type)
	message := []byte("i3-ipc")
	message = append(message, 0, 0, 0, 0) // payload length (0)
	message = append(message, 4, 0, 0, 0) // GET_TREE message type

	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send i3-ipc request: %w", err)
	}

	// Read response header (14 bytes: 6 magic + 4 length + 4 type)
	header := make([]byte, 14)
	_, err = conn.Read(header)
	if err != nil {
		return nil, fmt.Errorf("failed to read Sway response header: %w", err)
	}

	// Extract payload length from header
	payloadLength := int(header[6]) | int(header[7])<<8 | int(header[8])<<16 | int(header[9])<<24

	// Read the JSON payload
	payload := make([]byte, payloadLength)
	_, err = conn.Read(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to read Sway JSON payload: %w", err)
	}

	var root swayNode
	if err := json.Unmarshal(payload, &root); err != nil {
		return nil, fmt.Errorf("failed to decode Sway JSON response: %w", err)
	}

	focused := s.findFocusedNode(&root)
	if focused == nil {
		return nil, fmt.Errorf("no focused window found in Sway tree")
	}

	var title, class string
	var pid int

	if focused.WindowProperties != nil {
		if focused.WindowProperties.Title != nil {
			title = *focused.WindowProperties.Title
		}
		if focused.WindowProperties.Class != nil {
			class = *focused.WindowProperties.Class
		}
	}

	if focused.PID != nil {
		pid = *focused.PID
	}

	// Try to get workspace information
	workspace := ""
	if focused.Representation != nil {
		workspace = *focused.Representation
	} else if focused.Name != nil {
		workspace = *focused.Name
	}

	return &window.WindowInfo{
		Title:     title,
		Class:     class,
		PID:       pid,
		Workspace: workspace,
	}, nil
}

// findFocusedNode recursively searches the Sway tree for the focused window
func (s *SwayProvider) findFocusedNode(node *swayNode) *swayNode {
	// Check if this node is focused and has window properties
	if node.Focused && node.WindowProperties != nil {
		return node
	}

	// Recursively search child nodes
	for _, child := range node.Nodes {
		if found := s.findFocusedNode(child); found != nil {
			return found
		}
	}

	// Search floating nodes as well
	for _, floating := range node.FloatingNodes {
		if found := s.findFocusedNode(floating); found != nil {
			return found
		}
	}

	return nil
}