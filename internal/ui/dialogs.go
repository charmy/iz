package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m App) renderWithConfirmDialog(mainView string) string {
	modalWidth := 50
	modalHeight := 10

	// Create modal content
	visibleNodes := m.getVisibleNodes()
	commandName := "Unknown"
	if m.Cursor < len(visibleNodes) {
		commandName = visibleNodes[m.Cursor].Name
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Width(modalWidth-4).
		Padding(0, 1).
		Render("Run Command?")

	commandText := "Unknown command"
	if m.Cursor < len(visibleNodes) && visibleNodes[m.Cursor].Command != "" {
		commandText = visibleNodes[m.Cursor].Command
	}

	nameText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Align(lipgloss.Center).
		Width(modalWidth - 4).
		Render(fmt.Sprintf("Task: %s", commandName))

	commandDisplay := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("236")).
		Align(lipgloss.Center).
		Width(modalWidth-4).
		Padding(0, 1).
		Render(fmt.Sprintf("$ %s", commandText))

	// Yes/No buttons
	yesStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("40")).
		Bold(true)
	noStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	if m.ConfirmYes {
		yesStyle = yesStyle.Width(10).Align(lipgloss.Center).Background(lipgloss.Color("40")).Foreground(lipgloss.Color("0"))
	} else {
		noStyle = noStyle.Width(10).Align(lipgloss.Center).Background(lipgloss.Color("196")).Foreground(lipgloss.Color("0"))
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		yesStyle.Render(" YES "),
		"  ",
		noStyle.Render(" NO "),
	)

	modalContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title, "",
		nameText, "",
		commandDisplay, "",
		buttons,
	)

	// Create modal box
	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1).
		Align(lipgloss.Center)

	modal := modalStyle.Render(modalContent)

	// Overlay modal on dimmed background
	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("234")),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("240")),
	)
}

func (m App) renderWithInputDialog(mainView string) string {
	modalWidth := 60
	modalHeight := 6 + len(m.InputFields)*2

	// Create modal content
	visibleNodes := m.getVisibleNodes()
	commandName := "Unknown"
	if m.Cursor < len(visibleNodes) {
		commandName = visibleNodes[m.Cursor].Name
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Width(modalWidth-4).
		Padding(0, 1).
		Render("Enter Parameters")

	taskText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Align(lipgloss.Center).
		Width(modalWidth - 4).
		Render(fmt.Sprintf("Command: %s", commandName))

	commandDisplay := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("236")).
		Align(lipgloss.Center).
		Width(modalWidth-4).
		Padding(0, 1).
		Render(fmt.Sprintf("$ %s", m.PendingCommand))

	var inputs []string
	for i, field := range m.InputFields {
		label := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true).
			Width(modalWidth - 8).
			Render(fmt.Sprintf("%s:", field.Name))

		inputs = append(inputs, label)

		if field.IsChoice {
			// Render choice selector
			var choices []string
			for j, option := range field.Options {
				style := lipgloss.NewStyle()
				label := option.Label

				// If this is the custom option and we have a custom value, show it
				if option.Value == "custom" && field.CustomInput.Value() != "" {
					label = fmt.Sprintf("Custom: %s", field.CustomInput.Value())
				}

				if j == field.Choice {
					if i == m.InputCursor && !field.ShowCustomInput {
						// Active choice, current field
						style = style.Foreground(lipgloss.Color("0")).Background(lipgloss.Color("39")).Bold(true)
					} else {
						// Selected choice, inactive field
						style = style.Foreground(lipgloss.Color("39")).Bold(true)
					}
					choices = append(choices, style.Render(fmt.Sprintf("● %s", label)))
				} else {
					style = style.Foreground(lipgloss.Color("240"))
					choices = append(choices, style.Render(fmt.Sprintf("○ %s", label)))
				}
			}

			choiceStyle := lipgloss.NewStyle().
				Width(modalWidth-8).
				Border(lipgloss.RoundedBorder()).
				Padding(0, 1)

			if i == m.InputCursor && !field.ShowCustomInput {
				choiceStyle = choiceStyle.BorderForeground(lipgloss.Color("39"))
			} else {
				choiceStyle = choiceStyle.BorderForeground(lipgloss.Color("240"))
			}

			inputs = append(inputs, choiceStyle.Render(strings.Join(choices, "\n")))

			// Show custom input if "custom" is selected
			if field.ShowCustomInput {
				inputs = append(inputs, field.CustomInput.View())
			}
		} else {
			// Render text input
			inputs = append(inputs, field.TextInput.View())
		}
	}

	modalContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		taskText,
		"",
		commandDisplay,
		"",
		strings.Join(inputs, "\n"),
	)

	// Create modal box
	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1).
		Align(lipgloss.Center)

	modal := modalStyle.Render(modalContent)

	// Overlay modal on dimmed background
	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("234")),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("240")),
	)
}

func (m App) renderWithHelpDialog(mainView string) string {
	helpView := m.Help.View(m.Keys)

	// Create a nice help box
	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(m.Width - 4).
		Align(lipgloss.Center)

	helpContent := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Align(lipgloss.Center).
			Render("🔍 iz - Interactive Command Manager"),
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Align(lipgloss.Center).
			Render("A powerful TUI for managing and executing commands with variables"),
		"",
		helpView,
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Align(lipgloss.Center).
			Render("Press ? again or Esc to close help"),
	)

	helpBox := helpStyle.Render(helpContent)

	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		helpBox,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("234")),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("240")),
	)
}
