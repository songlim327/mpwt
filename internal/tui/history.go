package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// history represents the state of history component
type history struct {
	width     int
	height    int
	list      list.Model
	keys      *historyDelegateKeyMap
	tuiConfig *TuiConfig
}

// newHistory creates a new history view
// It reads the history data from the database and populates the list
func newHistory(tuiConf *TuiConfig) (*history, error) {
	items := []list.Item{}

	histories, err := tuiConf.Repository.ReadHistory()
	if err != nil {
		return nil, fmt.Errorf("failed to read history: %v", err)
	}

	for _, h := range histories {
		maxCmdsLength := 20
		shortCmds := h.Cmds
		if len(h.Cmds) > maxCmdsLength {
			shortCmds = h.Cmds[:maxCmdsLength]
		}
		items = append(items, cmdItem{
			title: fmt.Sprintf("(%d panes) %s...", h.PaneCount, shortCmds),
			desc:  h.ExecutedAt.Format("02/01/2006 15:04:00"),
			cmds:  h.Cmds,
			wtCmd: h.Wtcmd,
		})
	}

	keys := newHistoryDelegateKeyMap()
	l := list.New(items, newHistoryDelegate(keys), 0, 0)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)

	return &history{
		list:      l,
		keys:      keys,
		tuiConfig: tuiConf,
	}, nil
}

// setWidth sets the width of the history component
func (h *history) setWidth(width int) {
	h.width = width
}

// setHeight sets the height of the history component
func (h *history) setHeight(height int) {
	h.height = height
}

// Init is the bubbletea package ELM architecture specific functions
func (h *history) Init() tea.Cmd {
	return nil
}

// Update is the bubbletea package ELM architecture specific functions
func (h *history) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.keys.back):
			return h, tea.Batch(
				sendViewStrUpdate(MainView),
				sendStatusUpdate(""),
			)

		case key.Matches(msg, h.keys.favourite):
			i, ok := h.list.SelectedItem().(cmdItem)
			if ok {
				// Show favourite input view
				return h, tea.Batch(
					sendFavouriteInputUpdate(i.wtCmd, strings.Split(i.cmds, ",")),
					sendViewStrUpdate(FavouriteInputView),
					sendStatusUpdate(""),
				)
			}

		case key.Matches(msg, h.keys.launch):
			i, ok := h.list.SelectedItem().(cmdItem)
			if ok {
				// Execute the command
				cmd := exec.Command("cmd", "/C", i.wtCmd)
				if err := cmd.Run(); err != nil {
					return h, sendStatusUpdate(err.Error())
				}

				// Add command history to database
				err := h.tuiConfig.Repository.InsertHistory(i.wtCmd, strings.Split(i.cmds, ","))
				if err != nil {
					return h, sendStatusUpdate(err.Error())
				}

				return h, tea.Quit
			}
			return h, tea.Quit

		}
	}

	var cmd tea.Cmd
	h.list, cmd = h.list.Update(msg)
	return h, cmd
}

// View is the bubbletea package ELM architecture specific functions
func (h *history) View() string {
	h.list.SetSize(h.width, h.height)
	return h.list.View()
}
