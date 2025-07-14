package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Settings represents global application settings
type Settings struct {
	Confirm bool `yaml:"confirm"`
}

// VariableOption represents a predefined option for a variable
type VariableOption struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
}

// VariableConfig represents configuration for a command variable
type VariableConfig struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description,omitempty"`
	Default     string           `yaml:"default,omitempty"`
	Options     []VariableOption `yaml:"options,omitempty"`
}

// ConfigNode represents a node in the command tree from YAML
type ConfigNode struct {
	Name        string           `yaml:"name"`
	Expanded    bool             `yaml:"expanded"`
	Command     string           `yaml:"command,omitempty"`
	Description string           `yaml:"description,omitempty"`
	Confirm     *bool            `yaml:"confirm,omitempty"`
	Variables   []VariableConfig `yaml:"variables,omitempty"`
	Children    []ConfigNode     `yaml:"children,omitempty"`
}

// Config represents the main configuration structure
type Config struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description,omitempty"`
	Settings    Settings         `yaml:"settings,omitempty"`
	Variables   []VariableConfig `yaml:"variables,omitempty"`
	Commands    []ConfigNode     `yaml:"commands"`
}

// GetConfigPath returns the configuration file path
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "iz")
	configFile := filepath.Join(configDir, "config.yaml")

	return configFile, nil
}

// EnsureConfigExists creates default config if it doesn't exist
func EnsureConfigExists() error {
	configFile, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Check if config exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create config directory
		configDir := filepath.Dir(configFile)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("could not create config directory: %w", err)
		}

		// Create comprehensive default config
		defaultConfig := `name: "🔍iz - Interactive Command Manager"
description: "A powerful TUI for managing and executing commands"

settings:
  confirm: true

# Global variables - available to all commands
variables:
  - name: "host"
    description: "Default target hostname or IP"
    default: "localhost"
    options:
      - label: "Local (localhost)"
        value: "localhost"
      - label: "Production Server"
        value: "prod.example.com"
      - label: "Development Server"
        value: "dev.example.com"
      - label: "Custom..."
        value: "custom"
  - name: "user"
    description: "Default username"
    default: "root"

commands:
  - name: "System"
    expanded: false
    children:
      - name: "List Files"
        command: "ls -la"
        description: "List all files in current directory"
        confirm: false
      
      - name: "Current Directory"
        command: "pwd"
        description: "Show current working directory"
        confirm: false
        
      - name: "Disk Usage"
        command: "df -h"
        description: "Show disk space usage"
        confirm: false
        
      - name: "Memory Info"
        command: "free -h || top -l 1 | head -n 10"
        description: "Show memory usage (Linux/macOS compatible)"
        confirm: false

  - name: "Network"
    expanded: false
    children:
      - name: "Ping Host"
        command: "ping -c {count} {host}"
        description: "Ping a host with specified count"
        confirm: false
        variables:
          - name: "count"
            description: "Number of ping attempts"
            default: "4"
            options:
              - label: "Quick (1 ping)"
                value: "1"
              - label: "Normal (4 pings)"
                value: "4"
              - label: "Extended (10 pings)"
                value: "10"
              - label: "Custom..."
                value: "custom"
          - name: "host"
            description: "Target hostname or IP"
            default: "google.com"
            
      - name: "Check Port"
        command: "nc -zv {host} {port}"
        description: "Check if a port is open on a host"
        confirm: false
        variables:
          - name: "host"
            description: "Target hostname or IP"
            default: "localhost"
          - name: "port"
            description: "Port number to check"
            default: "80"
            options:
              - label: "HTTP (80)"
                value: "80"
              - label: "HTTPS (443)"
                value: "443"
              - label: "SSH (22)"
                value: "22"
              - label: "Custom..."
                value: "custom"

      - name: "SSH Connect"
        command: "ssh {user}@{host}"
        description: "Connect to a server via SSH (uses global variables)"
        confirm: false
        # This command uses global variables: host and user
        # No local variables defined, so it uses global defaults

      - name: "SSH with Custom Port"
        command: "ssh -p {port} {user}@{host}"
        description: "Connect via SSH with custom port"
        confirm: false
        variables:
          - name: "port"
            description: "SSH port number"
            default: "22"
        # host and user come from global variables, port is local

  - name: "Development"
    expanded: false
    children:
      - name: "Git Status"
        command: "git status"
        description: "Show git repository status"
        confirm: false
        
      - name: "Git Log"
        command: "git log --oneline -n {count}"
        description: "Show recent git commits"
        confirm: false
        variables:
          - name: "count"
            description: "Number of commits to show"
            default: "10"
            options:
              - label: "Last 5 commits"
                value: "5"
              - label: "Last 10 commits"
                value: "10"
              - label: "Last 20 commits"
                value: "20"
              - label: "Custom..."
                value: "custom"
`

		if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("could not create default config: %w", err)
		}

		fmt.Printf("Created default config at: %s\n", configFile)
	}

	return nil
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadConfig loads configuration with auto-creation
func LoadConfig() (*Config, error) {
	// Ensure config exists
	if err := EnsureConfigExists(); err != nil {
		return nil, err
	}

	// Get config path
	configFile, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Load config
	return LoadFromFile(configFile)
}

// GetFallbackConfig returns a default configuration when file loading fails
func GetFallbackConfig() *Config {
	return &Config{
		Name:        "iz - Command Manager",
		Description: "Fallback configuration",
		Settings:    Settings{Confirm: true},
		Commands: []ConfigNode{
			{
				Name:     "System",
				Expanded: false,
				Children: []ConfigNode{
					{
						Name:        "List Files",
						Command:     "ls -la",
						Description: "List all files",
					},
					{
						Name:        "Current Dir",
						Command:     "pwd",
						Description: "Show current directory",
					},
				},
			},
		},
	}
}
