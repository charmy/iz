package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmy/iz/internal/config"
	"github.com/charmy/iz/internal/tree"
)

// Update handles messages and updates the application state
func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" {
			return m, tea.Quit
		}
		if key == "?" {
			m.ShowHelp = !m.ShowHelp
			return m, nil
		}
		if key == "e" && !m.ShowInputs && !m.ShowConfirm && !m.ShowHelp {
			return m, m.openConfigInEditor()
		}
		if key == "esc" {
			// ESC behavior: close dialog if open, otherwise quit
			if m.ShowHelp {
				m.ShowHelp = false
				return m, nil
			} else if m.ShowInputs {
				// Let handleInputKeys handle ESC for input dialog
				return m.handleKeyPress(msg)
			} else if m.ShowConfirm {
				// Let handleConfirmKeys handle ESC for confirm dialog
				return m.handleKeyPress(msg)
			} else {
				// No dialog open, quit the application
				return m, tea.Quit
			}
		}
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Help.Width = msg.Width
	case tree.CommandFinishedMsg:
		return m, nil
	}
	return m, nil
}

func (m App) handleKeyPress(msg tea.Msg) (App, tea.Cmd) {
	if m.ShowInputs {
		return m.handleInputKeys(msg)
	}

	if m.ShowConfirm {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			return m.handleConfirmKeys(keyMsg.String())
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		key := keyMsg.String()
		switch key {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			visibleNodes := m.getVisibleNodes()
			if m.Cursor < len(visibleNodes)-1 {
				m.Cursor++
			}
		case "enter", "r":
			return m.handleEnter()
		}
	}
	return m, nil
}

func (m App) handleConfirmKeys(key string) (App, tea.Cmd) {
	// Handle quit keys in confirm mode
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	switch key {
	case "left", "h", "right", "l":
		m.ConfirmYes = !m.ConfirmYes
	case "enter":
		if m.ConfirmYes {
			m.ShowConfirm = false
			return m, tree.RunCommandInTerminal(m.PendingCommand)
		} else {
			m.ShowConfirm = false
		}
	case "esc":
		m.ShowConfirm = false
	}
	return m, nil
}

func (m App) handleEnter() (App, tea.Cmd) {
	visibleNodes := m.getVisibleNodes()
	if m.Cursor < len(visibleNodes) {
		node := visibleNodes[m.Cursor]
		if node.IsFolder {
			node.Expanded = !node.Expanded
			return m, nil
		} else if node.Command != "" {
			// Check if command has variables
			variables := tree.ExtractVariables(node.Command)
			if len(variables) > 0 {
				// Show input dialog for variables
				m.ShowInputs = true
				m.InputCursor = 0
				m.InputFields = []InputField{}

				for _, varName := range variables {
					// Check if this variable has predefined options
					var varConfig *config.VariableConfig
					for _, vc := range node.Variables {
						if vc.Name == varName {
							varConfig = &vc
							break
						}
					}

					if varConfig != nil && len(varConfig.Options) > 0 {
						// Create choice field
						defaultChoice := 0
						selectedValue := varConfig.Options[0].Value
						if varConfig.Default != "" {
							for i, opt := range varConfig.Options {
								if opt.Value == varConfig.Default {
									defaultChoice = i
									selectedValue = opt.Value
									break
								}
							}
						}

						// Create custom input for "custom" option
						customInput := textinput.New()
						customInput.Placeholder = fmt.Sprintf("Enter custom %s", varName)
						customInput.Width = 30
						customInput.CharLimit = 100

						m.InputFields = append(m.InputFields, InputField{
							Name:          varName,
							Placeholder:   fmt.Sprintf("Select %s", varName),
							IsChoice:      true,
							Options:       varConfig.Options,
							Choice:        defaultChoice,
							SelectedValue: selectedValue,
							CustomInput:   customInput,
						})
					} else {
						// Create text input field
						ti := textinput.New()
						ti.Placeholder = fmt.Sprintf("Enter %s", varName)
						ti.Width = 40
						ti.CharLimit = 100
						if varConfig != nil && varConfig.Default != "" {
							ti.SetValue(varConfig.Default)
						}

						m.InputFields = append(m.InputFields, InputField{
							Name:        varName,
							Placeholder: fmt.Sprintf("Enter %s", varName),
							TextInput:   ti,
						})
					}
				}

				// Focus first input
				if len(m.InputFields) > 0 {
					firstField := &m.InputFields[0]
					if !firstField.IsChoice {
						firstField.TextInput.Focus()
					}
				}
				m.PendingCommand = node.Command
				return m, nil
			} else {
				// No variables, proceed as normal
				if node.Confirm {
					m.ShowConfirm = true
					m.ConfirmYes = true
					m.PendingCommand = node.Command
					return m, nil
				} else {
					return m, tree.RunCommandInTerminal(node.Command)
				}
			}
		}
	}
	return m, nil
}

func (m App) handleInputKeys(msg tea.Msg) (App, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit keys in input mode
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle Shift+Tab for reverse navigation
		if msg.String() == "shift+tab" {
			// Switch to previous input field
			if m.InputCursor > 0 {
				// Blur current field
				currentField := &m.InputFields[m.InputCursor]
				if currentField.IsChoice && currentField.ShowCustomInput {
					currentField.CustomInput.Blur()
				} else if !currentField.IsChoice {
					currentField.TextInput.Blur()
				}

				m.InputCursor--

				// Focus previous field
				prevField := &m.InputFields[m.InputCursor]
				if !prevField.IsChoice {
					prevField.TextInput.Focus()
				} else if prevField.IsChoice && prevField.ShowCustomInput {
					prevField.CustomInput.Focus()
				}
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyTab:
			// Switch between input fields
			if m.InputCursor < len(m.InputFields)-1 {
				// Blur current field
				currentField := &m.InputFields[m.InputCursor]
				if currentField.IsChoice && currentField.ShowCustomInput {
					currentField.CustomInput.Blur()
				} else if !currentField.IsChoice {
					currentField.TextInput.Blur()
				}

				m.InputCursor++

				// Focus next field
				nextField := &m.InputFields[m.InputCursor]
				if !nextField.IsChoice {
					nextField.TextInput.Focus()
				} else if nextField.IsChoice && nextField.ShowCustomInput {
					nextField.CustomInput.Focus()
				}
			}
		case tea.KeyUp:
			currentField := &m.InputFields[m.InputCursor]
			if currentField.IsChoice && !currentField.ShowCustomInput {
				// Navigate choice options
				if currentField.Choice > 0 {
					currentField.Choice--
					currentField.SelectedValue = currentField.Options[currentField.Choice].Value

					// Check if "custom" is selected
					if currentField.SelectedValue == "custom" {
						currentField.ShowCustomInput = true
						currentField.CustomInput.Focus()
					} else {
						currentField.ShowCustomInput = false
						currentField.CustomInput.Blur()
					}
				} else {
					// If we're at the first choice, move to previous field
					if m.InputCursor > 0 {
						m.InputCursor--
						// Focus previous field
						prevField := &m.InputFields[m.InputCursor]
						if !prevField.IsChoice {
							prevField.TextInput.Focus()
						} else if prevField.IsChoice && prevField.ShowCustomInput {
							prevField.CustomInput.Focus()
						}
					}
				}
			} else {
				// Navigate between fields
				if m.InputCursor > 0 {
					// Blur current field
					if currentField.IsChoice && currentField.ShowCustomInput {
						currentField.CustomInput.Blur()
						currentField.ShowCustomInput = false
					} else if !currentField.IsChoice {
						currentField.TextInput.Blur()
					}

					m.InputCursor--

					// Focus previous field
					prevField := &m.InputFields[m.InputCursor]
					if !prevField.IsChoice {
						prevField.TextInput.Focus()
					} else if prevField.IsChoice && prevField.ShowCustomInput {
						prevField.CustomInput.Focus()
					}
					// Choice fields don't need explicit focus, they're navigated with Up/Down
				}
			}
		case tea.KeyDown:
			currentField := &m.InputFields[m.InputCursor]
			if currentField.IsChoice && !currentField.ShowCustomInput {
				// Navigate choice options
				if currentField.Choice < len(currentField.Options)-1 {
					currentField.Choice++
					currentField.SelectedValue = currentField.Options[currentField.Choice].Value

					// Check if "custom" is selected
					if currentField.SelectedValue == "custom" {
						currentField.ShowCustomInput = true
						currentField.CustomInput.Focus()
					} else {
						currentField.ShowCustomInput = false
						currentField.CustomInput.Blur()
					}
				} else {
					// If we're at the last choice, move to next field
					if m.InputCursor < len(m.InputFields)-1 {
						m.InputCursor++
						// Focus next field
						nextField := &m.InputFields[m.InputCursor]
						if !nextField.IsChoice {
							nextField.TextInput.Focus()
						} else if nextField.IsChoice && nextField.ShowCustomInput {
							nextField.CustomInput.Focus()
						}
					}
				}
			} else {
				// Navigate between fields
				if m.InputCursor < len(m.InputFields)-1 {
					// Blur current field
					if currentField.IsChoice && currentField.ShowCustomInput {
						currentField.CustomInput.Blur()
						currentField.ShowCustomInput = false
					} else if !currentField.IsChoice {
						currentField.TextInput.Blur()
					}

					m.InputCursor++

					// Focus next field
					nextField := &m.InputFields[m.InputCursor]
					if !nextField.IsChoice {
						nextField.TextInput.Focus()
					} else if nextField.IsChoice && nextField.ShowCustomInput {
						nextField.CustomInput.Focus()
					}
					// Choice fields don't need explicit focus, they're navigated with Up/Down
				}
			}
		case tea.KeyEnter:
			// Check if all fields are filled
			allFilled := true
			for _, field := range m.InputFields {
				if field.IsChoice {
					if field.SelectedValue == "custom" && strings.TrimSpace(field.CustomInput.Value()) == "" {
						allFilled = false
						break
					}
				} else if strings.TrimSpace(field.TextInput.Value()) == "" {
					allFilled = false
					break
				}
			}

			if allFilled {
				// Populate inputValues map
				m.InputValues = make(map[string]string)
				for _, field := range m.InputFields {
					if field.IsChoice {
						if field.SelectedValue == "custom" {
							m.InputValues[field.Name] = field.CustomInput.Value()
						} else {
							m.InputValues[field.Name] = field.SelectedValue
						}
					} else {
						m.InputValues[field.Name] = field.TextInput.Value()
					}
				}

				// Replace variables in command
				finalCommand := tree.ReplaceVariables(m.PendingCommand, m.InputValues)

				// Reset input state
				m.ShowInputs = false
				m.InputFields = []InputField{}

				// Check if this command needs confirmation
				visibleNodes := m.getVisibleNodes()
				if m.Cursor < len(visibleNodes) {
					node := visibleNodes[m.Cursor]
					if node.Confirm {
						m.ShowConfirm = true
						m.ConfirmYes = true
						m.PendingCommand = finalCommand
						return m, nil
					} else {
						return m, tree.RunCommandInTerminal(finalCommand)
					}
				}
			}
		case tea.KeyEsc:
			m.ShowInputs = false
			m.InputFields = []InputField{}
		}
	}

	// Update current input field
	if len(m.InputFields) > 0 && m.InputCursor < len(m.InputFields) {
		currentField := &m.InputFields[m.InputCursor]
		if currentField.IsChoice && currentField.ShowCustomInput {
			// Update custom input
			var cmd tea.Cmd
			currentField.CustomInput, cmd = currentField.CustomInput.Update(msg)
			return m, cmd
		} else if !currentField.IsChoice {
			// Update regular text input
			var cmd tea.Cmd
			currentField.TextInput, cmd = currentField.TextInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m App) getVisibleNodes() []*tree.TreeNode {
	var nodes []*tree.TreeNode
	m.collectVisibleNodes(m.Tree, &nodes, 0)
	return nodes
}

func (m App) collectVisibleNodes(node *tree.TreeNode, nodes *[]*tree.TreeNode, level int) {
	node.Level = level
	*nodes = append(*nodes, node)

	if node.Expanded && node.IsFolder {
		for _, child := range node.Children {
			m.collectVisibleNodes(child, nodes, level+1)
		}
	}
}

// openConfigInEditor opens the config file in the default text editor
func (m App) openConfigInEditor() tea.Cmd {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return nil
	}

	// Try different editors in order of preference
	editors := []string{
		os.Getenv("VISUAL"),
		os.Getenv("EDITOR"),
		"code",  // VS Code
		"nano",  // Nano
		"vim",   // Vim
		"vi",    // Vi
		"emacs", // Emacs
	}

	for _, editor := range editors {
		if editor == "" {
			continue
		}

		// Check if editor exists
		if _, err := exec.LookPath(editor); err == nil {
			return tea.ExecProcess(
				exec.Command(editor, configPath),
				func(err error) tea.Msg {
					// After editing, we might want to reload
					return tree.CommandFinishedMsg{}
				},
			)
		}
	}

	return nil
}
