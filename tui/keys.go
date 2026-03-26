package tui

import tea "github.com/charmbracelet/bubbletea"

type keyAction int

const (
	keyNone keyAction = iota
	keyQuit
	keyUp
	keyDown
	keyTab
	keyEnter
	keyRefresh
)

func parseKey(msg tea.KeyMsg) keyAction {
	switch msg.String() {
	case "q", "ctrl+c":
		return keyQuit
	case "up", "k":
		return keyUp
	case "down", "j":
		return keyDown
	case "tab", "l", "right":
		return keyTab
	case "shift+tab", "h", "left":
		return keyTab // toggle (only 2 tabs)
	case "enter":
		return keyEnter
	case "r":
		return keyRefresh
	}
	return keyNone
}
