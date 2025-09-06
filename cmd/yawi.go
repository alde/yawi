package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"yawi/pkg/compositor"
	"yawi/pkg/providers"
)

var (
	version = "dev" // This will be set by goreleaser
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "yawi",
	Short: "Yet Another Window Inspector - get active window information across platforms",
	Long: `YAWI is a simple tool to get information about the currently active window
across different platforms and window managers. By default, it outputs just the
window class name, making it perfect for use in scripts and automation.

Supported platforms: Hyprland, Sway, GNOME Shell (Linux), macOS`,
	RunE: func(cmd *cobra.Command, args []string) error {
		comp := compositor.Detect()
		if comp == compositor.Unknown {
			return fmt.Errorf("unable to detect supported platform\nSupported: Hyprland, Sway, GNOME Shell (Linux), macOS")
		}

		provider, err := providers.NewProvider(comp)
		if err != nil {
			return err
		}

		windowInfo, err := provider.GetActiveWindow()
		if err != nil {
			return fmt.Errorf("failed to get active window: %w", err)
		}

		// Default behavior: just output the class name for scripts
		fmt.Println(windowInfo.Class)
		return nil
	},
	// Custom unknown command handler
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			switch args[0] {
			case "full", "json":
				return fmt.Errorf("unknown command '%s'\nDid you mean 'yawi info'?", args[0])
			case "-v":
				return fmt.Errorf("unknown command '-v'\nDid you mean 'yawi version'?")
			}
		}
		return nil
	},
}

var compositorCmd = &cobra.Command{
	Use:   "compositor",
	Short: "Show which compositor is detected",
	RunE: func(cmd *cobra.Command, args []string) error {
		comp := compositor.Detect()
		fmt.Printf("Current compositor: %s\n", comp)
		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show full window information as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		comp := compositor.Detect()
		if comp == compositor.Unknown {
			return fmt.Errorf("unable to detect supported platform\nSupported: Hyprland, Sway, GNOME Shell (Linux), macOS")
		}

		provider, err := providers.NewProvider(comp)
		if err != nil {
			return err
		}

		windowInfo, err := provider.GetActiveWindow()
		if err != nil {
			return fmt.Errorf("failed to get active window: %w", err)
		}

		jsonData, err := json.MarshalIndent(windowInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to create JSON output: %w", err)
		}
		fmt.Println(string(jsonData))

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show YAWI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("YAWI version %s\n", version)
	},
}

func init() {
	// Disable default completion command since it might confuse users
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	
	// Set up suggestion function for unknown commands
	rootCmd.SuggestionsMinimumDistance = 1
	rootCmd.SuggestFor = []string{"ful", "josn", "jsn", "inf", "vers"}
	
	// Add subcommands
	rootCmd.AddCommand(compositorCmd)
	rootCmd.AddCommand(infoCmd) 
	rootCmd.AddCommand(versionCmd)
}