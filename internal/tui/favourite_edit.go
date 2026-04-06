package tui

import (
	"mpwt/internal/core"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// favouriteEditKeyMap is a map of key bindings for the favourite item edit view
type favouriteEditKeyMap struct {
	save key.Binding
	back key.Binding
	quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k favouriteEditKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.save, k.back, k.quit}
}

// FullHelp returns keybindings to be shown in the full help view
func (k favouriteEditKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.save, k.back, k.quit},
	}
}

// favouriteEditMsg represents a message struct to be displayed in the favourite edit component
type favouriteEditMsg struct {
	item cmdItem
}

// favouriteEdit represents the state of favourite edit component
type favouriteEdit struct {
	width     int
	height    int
	item      cmdItem
	textarea  textarea.Model
	help      help.Model
	keys      favouriteEditKeyMap
	textStyle lipgloss.Style
	tuiConfig *TuiConfig
}

// newFavouriteEdit returns a new favourite edit component
func newFavouriteEdit(tuiConf *TuiConfig) *favouriteEdit {
	ta := textarea.New()
	ta.Focus()

	var keys = favouriteEditKeyMap{
		save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
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

	return &favouriteEdit{
		textarea:  ta,
		help:      help.New(),
		keys:      keys,
		tuiConfig: tuiConf,
		textStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(TextColor)),
	}
}

// sendfavouriteEditUpdate send favouriteEditMsg which to be captured by the favourite edit component
func sendfavouriteEditUpdate(item cmdItem) func() tea.Msg {
	return func() tea.Msg {
		return favouriteEditMsg{item: item}
	}
}

// setWidth sets the width of the favouriteEdit component
func (f *favouriteEdit) setWidth(width int) {
	f.width = width
}

// setHeight sets the height of the favouriteEdit component
func (f *favouriteEdit) setHeight(height int) {
	f.height = height
}

// Init is the bubbletea package ELM architecture specific functions
func (f *favouriteEdit) Init() tea.Cmd {
	return nil
}

// Update is the bubbletea package ELM architecture specific functions
func (f *favouriteEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case favouriteEditMsg:
		f.item = msg.item
		f.textarea.SetValue(strings.Join(strings.Split(msg.item.cmds, ","), "\n"))

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.quit):
			return f, tea.Quit

		case key.Matches(msg, f.keys.back):
			return f, tea.Batch(
				sendViewStrUpdate(FavouriteView),
				sendStatusUpdate(""),
			)

		case key.Matches(msg, f.keys.save):
			// Recompute the command
			cmds := strings.Split(f.textarea.Value(), "\n")
			f.tuiConfig.TerminalConfig.Commands = cmds
			cmdStr, err := core.OpenWt(f.tuiConfig.TerminalConfig)
			if err != nil {
				return f, sendStatusUpdate(err.Error())
			}

			// Update favourite in database
			updateErr := f.tuiConfig.Repository.UpdateFavourite(f.item.id, f.item.title, cmdStr, strings.Split(f.textarea.Value(), "\n"))
			if updateErr != nil {
				return f, sendStatusUpdate(updateErr.Error())
			}
			return f, tea.Batch(
				sendFavouriteUpdate(),
				sendViewStrUpdate(FavouriteView),
				sendStatusUpdate("Favourite updated"),
			)
		}
	}

	var cmd tea.Cmd
	f.textarea, cmd = f.textarea.Update(msg)
	return f, cmd
}

// View is the bubbletea package ELM architecture specific functions
func (f *favouriteEdit) View() string {
	f.help.Width = f.width
	f.textarea.SetWidth(f.width)
	f.textarea.SetHeight(f.height - 1)

	return lipgloss.JoinVertical(lipgloss.Left, f.textarea.View(), f.help.View(f.keys))
}
