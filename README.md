# YAWI - Yet Another Window Inspector

YAWI is a simple, cross-platform tool for getting information about the currently active window. Whether you're on Linux with Wayland or rocking macOS, YAWI has got you covered.

## What's in the Box?

YAWI does one thing really well: it tells you about the window that's currently stealing your attention. By default, it just outputs the window class name (perfect for scripts), but it can also give you the full story in JSON format.

## Supported Platforms

- **Hyprland** - The fancy tiling compositor that makes everything look lagom
- **Sway** - i3's Wayland cousin
- **GNOME Shell** - The desktop environment that everyone either loves or... has opinions about
- **macOS** - Because sometimes you need to know what's happening in the Apple ecosystem

### Coming Soon
- **KDE/Plasma** - The customizable desktop that lets you tweak everything (TODO)

## Installation

### The Simple Way

```bash
go install github.com/yourusername/yawi/cmd/yawi@latest
```

### Or Build from Source

```bash
git clone https://github.com/yourusername/yawi.git
cd yawi
make build
# or: go build -o yawi ./cmd/yawi.go
```

## Usage

### The Simple Way (Perfect for Scripts)

```bash
# Just get the window class - clean and simple
$ yawi
Firefox

# Check which platform YAWI detected
$ yawi compositor
Current compositor: macOS
```

### The Fancy Way (For When You Want Details)

```bash
# Get full window information as JSON
$ yawi info
{
  "title": "YAWI - Yet Another Window Inspector",
  "class": "Firefox",
  "pid": 12345,
  "workspace": "main"
}
```

### Other Useful Commands

```bash
# Check YAWI version
$ yawi version
YAWI version 1.0.0

# Get help (always handy)
$ yawi --help
```

## Platform-Specific Notes

### Linux (Wayland Compositors)

YAWI works out of the box with:
- **Hyprland**: Uses the socket API for fast, reliable window detection
- **Sway**: Communicates via the i3-ipc protocol
- **GNOME Shell**: Requires the [Focused Window D-Bus extension](https://extensions.gnome.org/extension/5592/focused-window-dbus/) to be installed and enabled

### macOS

On macOS, YAWI uses AppleScript to get the frontmost application. No additional permissions needed - it works right away. Note that on macOS, the "window title" and "class" are both set to the application name since macOS handles windows a bit differently than Linux.

## Scripting Examples

### Simple Window Class Detection

```bash
#!/bin/bash
current_window=$(yawi)
if [[ "$current_window" == "Firefox" ]]; then
    echo "Time to focus! Close that browser."
fi
```

### JSON Parsing with jq

```bash
#!/bin/bash
window_info=$(yawi info)
title=$(echo "$window_info" | jq -r '.title')
echo "You're currently looking at: $title"
```

### Cross-Platform Workspace Detection

```bash
#!/bin/bash
workspace=$(yawi info | jq -r '.workspace')
echo "Current workspace: $workspace"
```

## Building from Source

```bash
# Get the dependencies
go mod download

# Build for your current platform
go build -o yawi ./cmd/yawi.go

# Or use the included Makefile
make build

# Run tests
make test
```

## Architecture

YAWI is structured in a way that makes adding new platforms straightforward:

- `pkg/compositor/` - Platform detection logic
- `pkg/window/` - Common window information structures
- `pkg/providers/` - Platform-specific implementations
- `cmd/` - CLI application entry point

Adding support for a new platform is as simple as implementing the `window.Provider` interface and updating the factory.

## Contributing

Found a bug? Want to add support for another platform? Contributions are welcome! The code is structured to make adding new platforms pretty painless.

## License

MIT License - because sharing is caring.

---

*YAWI: Because sometimes you just need to know what window is hogging your attention.*
