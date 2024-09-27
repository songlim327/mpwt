package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Underline(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(SelectionColor))
)

type item struct {
	title, desc string
}

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("🍊  " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// option represents the the main menu selection component
type option struct {
	list   list.Model
	width  int
	height int
}

// newOption creates a new option
func newOption() *option {
	var items = []list.Item{
		item{title: Execute, desc: ExecuteDesc},
		item{title: Exit, desc: ExitDesc},
	}

	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = "Welcome to MPWT, choose your option below:"
	l.SetShowStatusBar(false)
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
			i, ok := o.list.SelectedItem().(item)
			if ok {
				return o, func() tea.Msg {
					return statusMsg{message: i.desc}
				}
			}
			return o, nil

		case "down":
			o.list.CursorDown()
			i, ok := o.list.SelectedItem().(item)
			if ok {
				return o, func() tea.Msg {
					return statusMsg{message: i.desc}
				}
			}
			return o, nil

		case "enter":
			i, ok := o.list.SelectedItem().(item)
			if ok {
				if i.title == Exit {
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
