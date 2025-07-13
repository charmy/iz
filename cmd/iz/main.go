package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmy/iz/internal/config"
	"github.com/charmy/iz/internal/tree"
	"github.com/charmy/iz/internal/ui"
)

func main() {
	// Load configuration with auto-creation
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("Using fallback configuration...")
		cfg = config.GetFallbackConfig()
	}

	// Convert config to tree structure
	cmdTree := tree.BuildTreeFromConfig(cfg)

	// Create UI app
	app := ui.NewApp(cmdTree, cfg.Settings.Confirm)

	// Start the program
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
