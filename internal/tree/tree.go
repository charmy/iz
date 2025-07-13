package tree

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmy/iz/internal/config"
)

// TreeNode represents a hierarchical command structure
type TreeNode struct {
	Name        string
	Children    []*TreeNode
	Expanded    bool
	IsFolder    bool
	Level       int
	Command     string
	Description string
	Confirm     bool
	Variables   []config.VariableConfig
}

// ConvertConfigToTree converts configuration to tree structure
func ConvertConfigToTree(cfg *config.ConfigNode, defaultConfirm bool) *TreeNode {
	confirmSetting := defaultConfirm
	if cfg.Confirm != nil {
		confirmSetting = *cfg.Confirm
	}

	node := &TreeNode{
		Name:        cfg.Name,
		Expanded:    cfg.Expanded,
		IsFolder:    len(cfg.Children) > 0,
		Command:     cfg.Command,
		Description: cfg.Description,
		Confirm:     confirmSetting,
		Variables:   cfg.Variables,
	}

	for i := range cfg.Children {
		node.Children = append(node.Children, ConvertConfigToTree(&cfg.Children[i], defaultConfirm))
	}

	return node
}

// BuildTreeFromConfig creates tree structure from new config format
func BuildTreeFromConfig(cfg *config.Config) *TreeNode {
	defaultConfirm := true
	if cfg.Settings.Confirm {
		defaultConfirm = cfg.Settings.Confirm
	}

	root := &TreeNode{
		Name:        cfg.Name,
		Expanded:    true,
		IsFolder:    true,
		Description: cfg.Description,
		Confirm:     defaultConfirm,
	}

	for i := range cfg.Commands {
		root.Children = append(root.Children, ConvertConfigToTree(&cfg.Commands[i], defaultConfirm))
	}

	return root
}

// ExtractVariables extracts unique variable placeholders from command string
// Supports {variable} format and returns deduplicated list
func ExtractVariables(command string) []string {
	re := regexp.MustCompile(`\{(\w+)\}`)
	matches := re.FindAllStringSubmatch(command, -1)

	var variables []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			varName := match[1]
			if !seen[varName] {
				variables = append(variables, varName)
				seen[varName] = true
			}
		}
	}

	return variables
}

// ReplaceVariables substitutes variable placeholders with provided values
func ReplaceVariables(command string, values map[string]string) string {
	result := command
	for varName, value := range values {
		placeholder := fmt.Sprintf("{%s}", varName)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// CommandFinishedMsg signals completion of command execution
type CommandFinishedMsg struct{}

// RunCommandInTerminal executes command in terminal with user prompt to continue
func RunCommandInTerminal(command string) tea.Cmd {
	fullCommand := command + "; echo '\nPress Enter to continue...'; read"
	return tea.ExecProcess(
		exec.Command("sh", "-c", fullCommand),
		func(err error) tea.Msg {
			return CommandFinishedMsg{}
		},
	)
}
