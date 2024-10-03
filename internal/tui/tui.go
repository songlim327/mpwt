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

// View extends tea.Model interface
// Add width and height setter
type View interface {
	setWidth(int)
	setHeight(int)

	// tea.Model interface functions
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

// tui represents the state of main tui window
type tui struct {
	width          int
	height         int
	viewStr        string
	view           View
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

// viewStrMsg represents a message struct to trigger main window view changes
type viewStrMsg struct {
	viewStr string
}

// reloadMsg represents a message struct to trigger main window view reload after config changes
type reloadMsg struct{}

// sendViewStrUpdate send viewportMsg which to be captured by main window
func sendViewStrUpdate(viewStr string) func() tea.Msg {
	return func() tea.Msg {
		return viewStrMsg{viewStr: viewStr}
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
	o := newOption()
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
		viewStr:        MainView,
		view:           o, // default view: option
		TuiConfig:      tuiConf,
		status:         newStatus(""),
		footer:         newFooter(),
		option:         o,
		execute:        newExecute(tuiConf),
		history:        h,
		favourite:      f,
		favouriteInput: newFavouriteInput(tuiConf),
		settings:       s,
	}, nil
}

// mapViewStrToView is a helper function which maps the viewport to a View
func (t *tui) mapViewStrToView(viewStr string) View {
	switch viewStr {
	case ExecuteView:
		return t.execute
	case HistoryView:
		return t.history
	case FavouriteView:
		return t.favourite
	case FavouriteInputView:
		return t.favouriteInput
	case SettingsView:
		return t.settings
	default:
		return t.option
	}
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

	case viewStrMsg:
		t.viewStr = msg.viewStr
		t.view = t.mapViewStrToView(msg.viewStr)
		if t.viewStr == ExecuteView {
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
		// Forward keypress message to active view
		v, cmd := t.view.Update(msg)
		t.view = v.(View)
		return t, cmd
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

	// Set view width and height
	t.view.setWidth(boxWidth - padding*2)
	t.view.setHeight(boxHeight - padding - t.footer.style.GetHeight())

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
					t.view.View(),
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
