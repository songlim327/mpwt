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
	action string
	item   cmdItem
}

// favouriteInput represents the state of favourite input component
type favouriteInput struct {
	width     int
	height    int
	item      cmdItem
	action    string
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
func sendFavouriteInputUpdate(action string, item cmdItem) func() tea.Msg {
	return func() tea.Msg {
		return favouriteInputMsg{
			action: action,
			item:   item,
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
		f.item = msg.item
		f.action = msg.action

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
			switch f.action {
			case FavouriteInputAdd:
				err := f.tuiConfig.Repository.InsertFavourite(name, f.item.wtCmd, strings.Split(f.item.cmds, ","))
				if err != nil {
					return f, sendStatusUpdate(err.Error())
				} else {
					f.input.SetValue("")
					return f, tea.Batch(
						sendFavouriteUpdate(),
						sendViewStrUpdate(MainView),
						sendStatusUpdate("Favourite added successfully"),
					)
				}
			case FavouriteInputEdit:
				err := f.tuiConfig.Repository.UpdateFavourite(f.item.id, name, f.item.wtCmd, strings.Split(f.item.cmds, ","))
				if err != nil {
					return f, sendStatusUpdate(err.Error())
				} else {
					f.input.SetValue("")
					return f, tea.Batch(
						sendFavouriteUpdate(),
						sendViewStrUpdate(FavouriteView),
						sendStatusUpdate("Favourite updated successfully"),
					)
				}
			default:
				// This case should not be hit
				return f, sendStatusUpdate("Unknown action")
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
		f.textStyle.Render(fmt.Sprintf("Panes: %d", len(strings.Split(f.item.cmds, ",")))),
		f.textStyle.Render(fmt.Sprintf("Commands: %s", f.item.cmds)),
		f.input.View(),
		empty,
		f.help.View(f.keys),
	)
}
