package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// cmdItem represents custom item for list.Model (used in history, favourite)
type cmdItem struct {
	id                       int
	title, desc, cmds, wtCmd string
}

func (i cmdItem) Title() string       { return i.title }
func (i cmdItem) Description() string { return i.desc }
func (i cmdItem) FilterValue() string { return i.title }

// optionItem represents custom item for list.Model (used in option)
type optionItem struct {
	title, desc string
}

func (i optionItem) FilterValue() string { return i.title }

// Custom styling for list.Model
var (
	selectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color(SelectionColor)).
				Foreground(lipgloss.Color(SelectionColor)).
				Padding(0, 0, 0, 1)
	selectedDescStyle       = selectedTitleStyle.Foreground(lipgloss.Color(RosewaterColor))
	simpleItemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	simpleSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(SelectionColor))
)

// optionDelegate is a custom list.Delegate for option view
type optionDelegate struct{}

func (d optionDelegate) Height() int                             { return 1 }
func (d optionDelegate) Spacing() int                            { return 0 }
func (d optionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d optionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(optionItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := simpleItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return simpleSelectedItemStyle.Render("🍊  " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// newOptionDelegate creates a new list.Model item delegate for option view
func newOptionDelegate() *optionDelegate {
	return &optionDelegate{}
}

// newHistoryDelegate creates a new history delegate from default item delegate with given key bindings
func newHistoryDelegate(keys *historyDelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// Custom selected item styles
	d.Styles.SelectedTitle = selectedTitleStyle
	d.Styles.SelectedDesc = selectedDescStyle

	// Custom help bindings for the history item delegate
	help := []key.Binding{keys.launch, keys.favourite, keys.back}

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
	back      key.Binding
	launch    key.Binding
	favourite key.Binding
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
		favourite: key.NewBinding(
			key.WithKeys("ctrl+f"),
			key.WithHelp("ctrl+f", "favourite"),
		),
	}
}

// newFavouriteDelegate creates a new favourite delegate from default item delegate
func newFavouriteDelegate(keys *favouriteDelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// Custom selected item styles
	d.Styles.SelectedTitle = selectedTitleStyle
	d.Styles.SelectedDesc = selectedDescStyle

	// Custom help bindings for the history item delegate
	help := []key.Binding{keys.launch, keys.delete, keys.back}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

// favouriteDelegateKeyMap is a map of key bindings for the favourite item delegate
type favouriteDelegateKeyMap struct {
	back   key.Binding
	launch key.Binding
	delete key.Binding
}

// newFavouriteDelegateKeyMap creates a new favouriteDelegateKeyMap with default bindings
func newFavouriteDelegateKeyMap() *favouriteDelegateKeyMap {
	return &favouriteDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to main menu"),
		),
		launch: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "launch"),
		),
		delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "delete favourite"),
		),
	}
}
