package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmy/iz/internal/config"
	"github.com/charmy/iz/internal/tree"
)

// App represents the main application state
type App struct {
	// Window dimensions
	Width  int
	Height int

	// Navigation state
	Cursor int
	Tree   *tree.TreeNode

	// Dialog states
	ShowConfirm    bool
	ConfirmYes     bool
	PendingCommand string
	DefaultConfirm bool

	// Input handling
	ShowInputs  bool
	InputFields []InputField
	InputCursor int
	InputValues map[string]string

	// Help system
	ShowHelp bool
	Help     help.Model
	Keys     KeyMap
}

// InputField represents a variable input field with support for choices and text input
type InputField struct {
	Name        string
	Placeholder string
	TextInput   textinput.Model

	// Choice field support
	IsChoice        bool
	Options         []config.VariableOption
	Choice          int
	SelectedValue   string
	ShowCustomInput bool
	CustomInput     textinput.Model
}

// KeyMap defines all keyboard shortcuts for the application
type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Help  key.Binding
}

// ShortHelp returns short help
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns full help
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Back, k.Help, k.Quit},
	}
}

// NewApp creates a new application instance with initialized components
func NewApp(treeRoot *tree.TreeNode, defaultConfirm bool) App {
	// Initialize help system
	h := help.New()
	h.Width = 80

	// Initialize keyboard shortcuts
	keys := KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter", "r"),
			key.WithHelp("enter/r", "run command"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}

	return App{
		Tree:           treeRoot,
		Cursor:         0,
		DefaultConfirm: defaultConfirm,
		InputValues:    make(map[string]string),
		Help:           h,
		Keys:           keys,
	}
}

// Init initializes the application
func (m App) Init() tea.Cmd {
	return nil
}
