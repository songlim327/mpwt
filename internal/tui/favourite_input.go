package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// favouriteInputKeyMap defines a set of keybindings for favourite input component
type favouriteInputKeyMap struct {
	save key.Binding
	back key.Binding
	quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
// It is part of the key.Map interface
func (k favouriteInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.save, k.back, k.quit}
}

// FullHelp returns keybindings to be shown in the full help view
// It is part of the key.Map interface
func (k favouriteInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.save, k.back, k.quit},
	}
}

// favouriteInputMsg represents a message struct to be displayed in the favourite input component
type favouriteInputMsg struct {
	cmds  []string
	wtCmd string
}

// favouriteInput represents the state of favourite input component
type favouriteInput struct {
	width     int
	height    int
	wtCmd     string
	cmds      []string
	input     textinput.Model
	help      help.Model
	keys      favouriteInputKeyMap
	textStyle lipgloss.Style
	tuiConfig *TuiConfig
}

// newFavouriteInput returns a new favourite input component
func newFavouriteInput(tuiConf *TuiConfig) *favouriteInput {
	ti := textinput.New()
	ti.Placeholder = "Enter the name"
	ti.Focus()
	ti.CharLimit = 100

	keys := favouriteInputKeyMap{
		save: key.NewBinding(
			key.WithKeys("enter", "ctrl+s"),
			key.WithHelp("enter/ctrl+s", "save"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to main menu"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}

	return &favouriteInput{
		input:     ti,
		help:      help.New(),
		tuiConfig: tuiConf,
		keys:      keys,
		textStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(TextColor)),
	}
}

// sendFavouriteInputUpdate sends favouriteInputMsg to be captured by the favourite input component
func sendFavouriteInputUpdate(wtCmd string, cmds []string) func() tea.Msg {
	return func() tea.Msg {
		return favouriteInputMsg{
			cmds:  cmds,
			wtCmd: wtCmd,
		}
	}
}

// setWidth sets the width of the favouriteInput component
func (f *favouriteInput) setWidth(width int) {
	f.width = width
}

// setHeight sets the height of the favouriteInput component
func (f *favouriteInput) setHeight(height int) {
	f.height = height
}

// Init is the bubbletea package ELM architecture specific functions
func (f *favouriteInput) Init() tea.Cmd {
	return nil
}

// Update is the bubbletea package ELM architecture specific functions
func (f *favouriteInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case favouriteInputMsg:
		f.cmds = msg.cmds
		f.wtCmd = msg.wtCmd

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.quit):
			return f, tea.Quit

		case key.Matches(msg, f.keys.back):
			return f, tea.Batch(
				sendViewStrUpdate(MainView),
				sendStatusUpdate(""),
			)

		case key.Matches(msg, f.keys.save):
			name := f.input.Value()
			err := f.tuiConfig.Repository.InsertFavourite(name, f.wtCmd, f.cmds)
			if err != nil {
				return f, sendStatusUpdate(err.Error())
			} else {
				return f, tea.Batch(
					sendFavouriteUpdate(),
					sendViewStrUpdate(MainView),
					sendStatusUpdate("Favourite saved successfully"),
				)
			}
		}
	}

	var cmd tea.Cmd
	f.input, cmd = f.input.Update(msg)
	return f, cmd
}

// View is the bubbletea package ELM architecture specific functions
func (f *favouriteInput) View() string {
	f.input.Width = f.width
	emptyHeight := f.height - 4 // height of each textStyle (1x2), input.Model(1), help.Model(1)
	empty := lipgloss.NewStyle().Height(emptyHeight).Render("")

	return lipgloss.JoinVertical(lipgloss.Left,
		f.textStyle.Render(fmt.Sprintf("Panes: %d", len(f.cmds))),
		f.textStyle.Render(fmt.Sprintf("Commands: %s", strings.Join(f.cmds, ","))),
		f.input.View(),
		empty,
		f.help.View(f.keys),
	)
}
