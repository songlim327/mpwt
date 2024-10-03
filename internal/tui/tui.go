package tui

import (
	"mpwt/internal/config"
	"mpwt/internal/core"
	"mpwt/internal/repository"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TuiConfig represents the configuration for tui application
type TuiConfig struct {
	TerminalConfig *core.TerminalConfig
	Repository     repository.IRepository
	ConfigMgr      config.IConfigManager
}

// tui represents the state of main tui window
type tui struct {
	width          int
	height         int
	viewport       string
	TuiConfig      *TuiConfig
	status         *status
	footer         *footer
	option         *option
	execute        *execute
	history        *history
	favourite      *favourite
	favouriteInput *favouriteInput
	settings       *settings
}

// viewportMsg represents a message struct to trigger main window view changes
type viewportMsg struct {
	viewport string
}

// reloadMsg represents a message struct to trigger main window view reload after config changes
type reloadMsg struct{}

// sendViewportUpdate send viewportMsg which to be captured by main window
func sendViewportUpdate(viewport string) func() tea.Msg {
	return func() tea.Msg {
		return viewportMsg{viewport: viewport}
	}
}

// sendReloadUpdate send reloadMsg which to be captured by main window and reload application state
func sendReloadUpdate() func() tea.Msg {
	return func() tea.Msg {
		return reloadMsg{}
	}
}

// newTui creates a new tui (main window view)
func newTui(tuiConf *TuiConfig) (*tui, error) {
	h, err := newHistory(tuiConf)
	if err != nil {
		return nil, err
	}

	f, err := newFavourite(tuiConf)
	if err != nil {
		return nil, err
	}

	s, err := newSettings(tuiConf)
	if err != nil {
		return nil, err
	}

	return &tui{
		viewport:       MainView,
		TuiConfig:      tuiConf,
		status:         newStatus(""),
		footer:         newFooter(),
		option:         newOption(),
		execute:        newExecute(tuiConf),
		history:        h,
		favourite:      f,
		favouriteInput: newFavouriteInput(tuiConf),
		settings:       s,
	}, nil
}

// Init is the bubbletea package ELM architecture specific functions
func (t *tui) Init() tea.Cmd {
	return tea.SetWindowTitle("🍊 MPWT")
}

// Update is the bubbletea package ELM architecture specific functions
func (t *tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height

	case viewportMsg:
		t.viewport = msg.viewport
		if t.viewport == ExecuteView {
			s, cmd := t.status.Update(statusMsg{message: "Each line of command will spawn a new pane in terminal"})
			t.status = s.(*status)
			return t, cmd
		}

	case statusMsg:
		s, cmd := t.status.Update(msg)
		t.status = s.(*status)
		return t, cmd

	case favouriteMsg:
		f, cmd := t.favourite.Update(msg)
		t.favourite = f.(*favourite)
		return t, cmd

	case favouriteInputMsg:
		i, cmd := t.favouriteInput.Update(msg)
		t.favouriteInput = i.(*favouriteInput)
		return t, cmd

	case reloadMsg:
		// Read config from file
		conf, err := t.TuiConfig.ConfigMgr.ReadConfig()
		if err != nil {
			return t, sendStatusUpdate(err.Error())
		}

		// Reload terminal application config
		t.TuiConfig.TerminalConfig = &core.TerminalConfig{
			Maximize:     conf.Maximize,
			Direction:    conf.Direction,
			Columns:      conf.Columns,
			OpenInNewTab: conf.OpenInNewTab,
		}

		// Recreate view requiring TerminalConfig
		t.execute = newExecute(t.TuiConfig)

	case tea.KeyMsg:
		switch t.viewport {
		case ExecuteView:
			e, cmd := t.execute.Update(msg)
			t.execute = e.(*execute)
			return t, cmd
		case HistoryView:
			h, cmd := t.history.Update(msg)
			t.history = h.(*history)
			return t, cmd
		case FavouriteView:
			f, cmd := t.favourite.Update(msg)
			t.favourite = f.(*favourite)
			return t, cmd
		case FavouriteInputView:
			i, cmd := t.favouriteInput.Update(msg)
			t.favouriteInput = i.(*favouriteInput)
			return t, cmd
		case SettingsView:
			s, cmd := t.settings.Update(msg)
			t.settings = s.(*settings)
			return t, cmd
		default:
			o, cmd := t.option.Update(msg)
			t.option = o.(*option)
			return t, cmd
		}
	}
	return t, nil
}

// View is the bubbletea package ELM architecture specific functions
func (t *tui) View() string {
	margin := 2
	padding := 1
	gap := 1

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderForegroundColor))

	// Window and children size calculations
	boxWidth := t.width - margin*2
	boxHeight := t.height - margin*2 - boxStyle.GetBorderTopSize()*2 - t.status.style.GetHeight() - t.status.style.GetBorderTopSize()*2

	t.status.width = boxWidth
	t.footer.width = boxWidth - padding*2

	var view string
	switch t.viewport {
	case ExecuteView:
		t.execute.width = boxWidth - padding*2
		t.execute.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.execute.View()
	case HistoryView:
		t.history.width = boxWidth - padding*2
		t.history.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.history.View()
	case FavouriteView:
		t.favourite.width = boxWidth - padding*2
		t.favourite.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.favourite.View()
	case FavouriteInputView:
		t.favouriteInput.width = boxWidth - padding*2
		t.favouriteInput.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.favouriteInput.View()
	case SettingsView:
		t.settings.width = boxWidth - padding*2
		t.settings.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.settings.View()
	default:
		t.option.width = boxWidth - padding*2
		t.option.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.option.View()
	}

	// Content box
	return lipgloss.NewStyle().
		Margin(margin).
		Render(
			t.status.View(),
			boxStyle.
				Width(boxWidth).
				Height(boxHeight).
				Padding(padding, padding, 0, padding).
				MarginTop(gap).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					view,
					t.footer.View(),
				),
			),
		)
}

// InitTea intialize a new tea program with user interactions
func InitTea(tc *TuiConfig) error {
	t, err := newTui(tc)
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		t,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
