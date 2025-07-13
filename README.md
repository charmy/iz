# 🔍 iz - Interactive Command Manager

A terminal-based interactive command manager that organizes and executes commands with variables.

## Installation

```bash
go install github.com/charmy/iz/cmd/iz@latest
```

## Usage

```bash
iz
```

On first run, `~/.config/iz/config.yaml` is automatically created.

## Features

- 📋 Hierarchical command organization
- 🔧 Variable support in commands
- ✅ Security confirmation system
- 🎨 Terminal UI

## Configuration

Edit `~/.config/iz/config.yaml` to add your commands:

```yaml
settings:
  confirm: true

commands:
  - name: "System"
    children:
      - name: "List Files"
        command: "ls -la"
        confirm: false
      
      - name: "Ping Host"
        command: "ping -c {count} {host}"
        variables:
          - name: "count"
            default: "4"
          - name: "host"
            default: "google.com"
```

## Keyboard Shortcuts

- `↑/↓` or `j/k` - Navigate
- `Enter/r` - Run command
- `e` - Edit config
- `?` - Help
- `q` - Quit

### Config Editor

- `Ctrl+S` - Save config
- `Esc` - Cancel editing
