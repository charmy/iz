package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmy/iz/internal/tree"
)

// View renders the main UI
func (m App) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	contentHeight := m.Height - 3
	paneWidth := (m.Width - 4) / 2

	panes := []string{
		m.renderPane("Commands", m.renderTree(), paneWidth, contentHeight),
		m.renderPane("Details", m.renderDetails(), paneWidth, contentHeight),
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, panes...)
	mainView := lipgloss.JoinVertical(lipgloss.Left, content, m.renderStatusBar())

	if m.ShowHelp {
		return m.renderWithHelpDialog(mainView)
	}

	if m.ShowInputs {
		return m.renderWithInputDialog(mainView)
	}

	if m.ShowConfirm {
		return m.renderWithConfirmDialog(mainView)
	}

	return mainView
}

func (m App) renderPane(title, content string, width, height int) string {
	paneStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	finalContent := titleStyle.Render(title) + "\n\n" + content
	return paneStyle.Render(finalContent)
}

func (m App) renderTree() string {
	visibleNodes := m.getVisibleNodes()
	var lines []string

	for i, node := range visibleNodes {
		lines = append(lines, m.renderNode(node, i == m.Cursor))
	}

	return strings.Join(lines, "\n")
}

func (m App) renderDetails() string {
	visibleNodes := m.getVisibleNodes()
	if m.Cursor >= len(visibleNodes) {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Align(lipgloss.Center)
		return emptyStyle.Render("No selection")
	}

	selected := visibleNodes[m.Cursor]
	var content []string

	// Name with highlight
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	content = append(content, nameStyle.Render(fmt.Sprintf("📋 %s", selected.Name)))
	content = append(content, "")

	if selected.IsFolder {
		// Folder details
		typeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			Bold(true)
		content = append(content, typeStyle.Render("📁 FOLDER"))

		childrenStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))
		content = append(content, childrenStyle.Render(fmt.Sprintf("└─ %d items", len(selected.Children))))

		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))
		if selected.Expanded {
			content = append(content, statusStyle.Render("✓ Expanded"))
		} else {
			statusStyle = statusStyle.Foreground(lipgloss.Color("202"))
			content = append(content, statusStyle.Render("⊕ Collapsed"))
		}
	} else {
		// Command details
		typeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			Bold(true)
		content = append(content, typeStyle.Render("⚡ COMMAND"))
		content = append(content, "")

		if selected.Command != "" {
			commandStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Background(lipgloss.Color("237")).
				Padding(0, 1)
			content = append(content, "Command:")
			content = append(content, commandStyle.Render(fmt.Sprintf("$ %s", selected.Command)))
			content = append(content, "")
		}

		if selected.Description != "" {
			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("250")).
				Background(lipgloss.Color("238")).
				Padding(0, 1).
				Italic(true)
			content = append(content, "Description:")
			content = append(content, descStyle.Render(selected.Description))
		}

		// Add action hint
		content = append(content, "")
		hintStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Bold(true)
		content = append(content, hintStyle.Render("💡 Press Enter to run"))
	}

	return strings.Join(content, "\n")
}

func (m App) renderStatusBar() string {
	statusStyle := lipgloss.NewStyle().
		Width(m.Width).
		Background(lipgloss.Color("237")).
		Foreground(lipgloss.Color("250")).
		Padding(0, 1)

	if m.ShowInputs {
		return statusStyle.Render("↑/↓ or Tab/Shift+Tab to switch fields • Enter when all filled • ESC to go back")
	}

	if m.ShowConfirm {
		return statusStyle.Render("Use ←/→ to select • Enter to confirm • ESC to go back")
	}

	if m.ShowHelp {
		return statusStyle.Render("Keyboard shortcuts • ESC to go back")
	}

	return statusStyle.Render("Use ↑/↓ to navigate • Enter/r to run • e to edit config • ? for help • ESC to quit")
}

func (m App) renderNode(node *tree.TreeNode, selected bool) string {
	indent := strings.Repeat("  ", node.Level)

	var prefix string
	if node.IsFolder {
		if node.Expanded {
			prefix = "▼ "
		} else {
			prefix = "▶ "
		}
	} else {
		prefix = "• "
	}

	line := indent + prefix + node.Name

	if selected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Render(line)
	}

	return line
}
