package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(0, 2)

	tabBarStyle = lipgloss.NewStyle().
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			MarginBottom(1)

	// List item styles
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	repoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	draftStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	// CI status styles
	ciSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	ciFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	ciPendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	ciNoneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	// Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginBottom(1)

	// Error
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	// Loading
	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
)
