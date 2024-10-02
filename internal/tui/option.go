package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Underline(true)
)

// option represents the the main menu selection component
type option struct {
	list   list.Model
	width  int
	height int
}

// newOption creates a new option
func newOption() *option {
	items := []list.Item{
		optionItem{title: ExecuteView, desc: ExecuteViewDesc},
		optionItem{title: HistoryView, desc: HistoryViewDesc},
		optionItem{title: ExitView, desc: ExitViewDesc},
	}

	l := list.New(items, newOptionDelegate(), 0, 0)
	l.Title = "Welcome to MPWT, choose your option below:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	return &option{list: l}
}

// Init is the bubbletea package ELM architecture specific functions
func (o *option) Init() tea.Cmd {
	return nil
}

// Update is the bubbletea package ELM architecture specific functions
func (o *option) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return o, tea.Quit

		case "up":
			o.list.CursorUp()
			i, ok := o.list.SelectedItem().(optionItem)
			if ok {
				return o, func() tea.Msg {
					return statusMsg{message: i.desc}
				}
			}
			return o, nil

		case "down":
			o.list.CursorDown()
			i, ok := o.list.SelectedItem().(optionItem)
			if ok {
				return o, func() tea.Msg {
					return statusMsg{message: i.desc}
				}
			}
			return o, nil

		case "enter":
			i, ok := o.list.SelectedItem().(optionItem)
			if ok {
				if i.title == ExitView {
					return o, tea.Quit
				}
				return o, func() tea.Msg {
					return viewportMsg{viewport: i.title}
				}
			}
			return o, tea.Quit
		}
	}

	var cmd tea.Cmd
	o.list, cmd = o.list.Update(msg)
	return o, cmd
}

// View is the bubbletea package ELM architecture specific functions
func (o *option) View() string {
	o.list.SetWidth(o.width)
	o.list.SetHeight(o.height)
	return o.list.View()
}
