package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// newHistoryDelegate creates a new history delegate from default item delegate with given key bindings
func newHistoryDelegate(keys *historyDelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// Custom selected item styles
	d.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(SelectionColor)).
		Foreground(lipgloss.Color(SelectionColor)).
		Padding(0, 0, 0, 1)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Foreground(lipgloss.Color(RosewaterColor))

	// Custom help bindings for the history item delegate
	help := []key.Binding{keys.launch, keys.back}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

// historyDelegateKeyMap is a map of key bindings for the history item delegate
type historyDelegateKeyMap struct {
	back   key.Binding
	launch key.Binding
}

// newHistoryDelegateKeyMap creates a new historyDelegateKeyMap with default bindings
func newHistoryDelegateKeyMap() *historyDelegateKeyMap {
	return &historyDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to main menu"),
		),
		launch: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "launch"),
		),
	}
}
