package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type SwayNode struct {
	ID                 int         `json:"id"`
	Name               *string     `json:"name"`
	Type               string      `json:"type"`
	Border             string      `json:"border"`
	CurrentBorderWidth int         `json:"current_border_width"`
	Layout             string      `json:"layout"`
	Orientation        string      `json:"orientation"`
	Percent            *float64    `json:"percent"`
	Rect               SwayRect    `json:"rect"`
	WindowRect         SwayRect    `json:"window_rect"`
	DecoRect           SwayRect    `json:"deco_rect"`
	Geometry           SwayRect    `json:"geometry"`
	Urgent             bool        `json:"urgent"`
	Focused            bool        `json:"focused"`
	Focus              []int       `json:"focus"`
	Nodes              []*SwayNode `json:"nodes"`
	FloatingNodes      []*SwayNode `json:"floating_nodes"`
	Sticky             bool        `json:"sticky"`
	Representation     *string     `json:"representation"`
	AppID              *string     `json:"app_id"`
	WindowProperties   *struct {
		Class        *string `json:"class"`
		Instance     *string `json:"instance"`
		Title        *string `json:"title"`
		TransientFor *int    `json:"transient_for"`
	} `json:"window_properties"`
	PID *int `json:"pid"`
}

type SwayRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func GetSwayActiveWindow() (*SwayNode, error) {
	socketPath := os.Getenv("SWAYSOCK")
	if socketPath == "" {
		return nil, fmt.Errorf("SWAYSOCK environment variable not found")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Sway: %w", err)
	}
	defer conn.Close()

	message := []byte("i3-ipc")
	message = append(message, 0, 0, 0, 0)
	message = append(message, 4, 0, 0, 0)

	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	header := make([]byte, 14)
	_, err = conn.Read(header)
	if err != nil {
		return nil, fmt.Errorf("failed to read Sway's response header: %w", err)
	}

	payloadLength := int(header[6]) | int(header[7])<<8 | int(header[8])<<16 | int(header[9])<<24

	payload := make([]byte, payloadLength)
	_, err = conn.Read(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to read Sway's JSON: %w", err)
	}

	var root SwayNode
	if err := json.Unmarshal(payload, &root); err != nil {
		return nil, fmt.Errorf("failed to decode Sway's JSON: %w", err)
	}

	focused := findFocusedNode(&root)
	if focused == nil {
		return nil, fmt.Errorf("no focused window found in Sway")
	}

	return focused, nil
}

func findFocusedNode(node *SwayNode) *SwayNode {

	if node.Focused && node.WindowProperties != nil {
		return node
	}

	for _, child := range node.Nodes {
		if found := findFocusedNode(child); found != nil {
			return found
		}
	}

	for _, floating := range node.FloatingNodes {
		if found := findFocusedNode(floating); found != nil {
			return found
		}
	}

	return nil
}
