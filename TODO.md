# YAWI TODO List

## Platform Support

### High Priority
- [ ] **KDE/Plasma Support** - Add support for KDE Plasma via KWin D-Bus interface
  - Research KWin's D-Bus API for active window information
  - Implement KDE provider similar to GNOME provider
  - Handle both X11 and Wayland sessions in KDE

### Medium Priority
- [ ] **Windows Support** - Add Windows window detection
  - Use Win32 API to get active window information
  - Handle different Windows versions gracefully
  - Test with various Windows applications

### Lower Priority
- [ ] **X11 Support** - Add support for traditional X11 window managers
  - Use xprop/xwininfo for window detection
  - Support common window managers (i3, dwm, awesome, etc.)
  - Graceful fallback for mixed X11/Wayland systems

## Features & Improvements

### Quality of Life
- [ ] **Configuration File** - Optional config file for custom behavior
- [ ] **Shell Completions** - Auto-generated completions for bash/zsh/fish
- [ ] **Better Error Messages** - More helpful error messages for edge cases
- [ ] **Logging** - Optional debug logging for troubleshooting

### Advanced Features
- [ ] **Window Filtering** - Filter by application name, title patterns, etc.
- [ ] **Multiple Windows** - Support for listing all windows, not just active
- [ ] **Window Actions** - Optional window manipulation (focus, minimize, etc.)
- [ ] **Watch Mode** - Continuous monitoring of active window changes

## Code Quality

### Testing
- [ ] **Unit Tests** - Add comprehensive unit tests for all providers
- [ ] **Integration Tests** - Test against real compositor environments
- [ ] **CI Testing** - Test on different Linux distributions and macOS versions

### Documentation
- [ ] **Man Page** - Traditional man page documentation
- [ ] **API Documentation** - Better inline documentation for public APIs
- [ ] **Contributing Guide** - Detailed guide for adding new platforms

## Future Distribution (Maybe Later)

### Package Managers (Low Priority)
- [ ] **Homebrew Formula** - Official Homebrew package for macOS/Linux
- [ ] **AUR Package** - Arch User Repository package

---

*Priority levels are flexible and can change based on user feedback and contributions.*