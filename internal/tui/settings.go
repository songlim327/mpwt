package tui

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// settingsKeyMap defines a set of keybindings for settings component
type settingsKeyMap struct {
	save key.Binding
	back key.Binding
	quit key.Binding
}

// ShortHelp implements the mini help view
// It is part of the key.Map interface
func (k settingsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.save, k.back, k.quit}
}

// FullHelp implements the full help view
// It is part of the key.Map interface
func (k settingsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.save, k.back, k.quit},
	}
}

// settings represents the state of a settings component
type settings struct {
	width     int
	height    int
	textarea  textarea.Model
	help      help.Model
	keys      settingsKeyMap
	tuiConfig *TuiConfig
}

// newSettings creates a new settings view
func newSettings(tuiConf *TuiConfig) (*settings, error) {
	buf, err := tuiConf.ConfigMgr.ReadConfigRaw()
	if err != nil {
		return nil, fmt.Errorf("failed to read raw config file: %v", err)
	}

	// Replace CRLF to LF
	buf = bytes.ReplaceAll(buf, []byte("\r\n"), []byte("\n"))

	ta := textarea.New()
	ta.SetValue(string(buf))
	ta.Focus()

	keys := settingsKeyMap{
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

	return &settings{
		textarea:  ta,
		help:      help.New(),
		keys:      keys,
		tuiConfig: tuiConf,
	}, nil
}

// Init is the bubbletea package ELM architecture specific functions
func (s *settings) Init() tea.Cmd {
	return nil
}

// Update is the bubbletea package ELM architecture specific functions
func (s *settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.quit):
			return s, tea.Quit

		case key.Matches(msg, s.keys.back):
			return s, tea.Batch(
				sendViewportUpdate(MainView),
				sendStatusUpdate(""),
			)

		case key.Matches(msg, s.keys.save):
			// Overwrite config file
			err := s.tuiConfig.ConfigMgr.WriteConfig(s.textarea.Value())
			if err != nil {
				return s, sendStatusUpdate(err.Error())
			}
			return s, tea.Batch(
				sendStatusUpdate("settings updated"),
				sendViewportUpdate(MainView),
			)
		}
	}

	var cmd tea.Cmd
	s.textarea, cmd = s.textarea.Update(msg)
	return s, cmd
}

// View is the bubbletea package ELM architecture specific functions
func (s *settings) View() string {
	s.help.Width = s.width
	s.textarea.SetWidth(s.width)
	s.textarea.SetHeight(s.height - 1) // height of help model

	return lipgloss.JoinVertical(lipgloss.Left,
		s.textarea.View(),
		s.help.View(s.keys),
	)
}
