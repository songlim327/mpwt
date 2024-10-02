package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// favourite represents the state of favourite component
type favourite struct {
	width     int
	height    int
	list      list.Model
	keys      *favouriteDelegateKeyMap
	tuiConfig *TuiConfig
}

// favouriteMsg updates the favourite list items by refetching the items from database
type favouriteMsg struct{}

// newFavourite creates a new favourite view
// It reads the favourite data from the database and populates the list
func newFavourite(tuiConf *TuiConfig) (*favourite, error) {
	items, err := loadItems(tuiConf)
	if err != nil {
		return nil, err
	}
	keys := newFavouriteDelegateKeyMap()
	l := list.New(items, newFavouriteDelegate(keys), 0, 0)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)

	return &favourite{
		list:      l,
		keys:      keys,
		tuiConfig: tuiConf,
	}, nil
}

// sendFavouriteUpdate send favouriteMsg which to be captured by the favourite component
func sendFavouriteUpdate() func() tea.Msg {
	return func() tea.Msg {
		return favouriteMsg{}
	}
}

// loadItems loads the favourite data from the database and populates the list
func loadItems(tuiConf *TuiConfig) ([]list.Item, error) {
	items := []list.Item{}

	favourites, err := tuiConf.Repository.ReadFavourite()
	if err != nil {
		return nil, err
	}

	for _, f := range favourites {
		items = append(items, cmdItem{
			title: f.Name,
			desc:  fmt.Sprintf("(%d panes) %s", len(strings.Split(f.Cmds, ",")), f.Cmds),
			cmds:  f.Cmds,
			wtCmd: f.Wtcmd,
		})
	}

	return items, nil
}

// Init is the bubbletea package ELM architecture specific functions
func (f *favourite) Init() tea.Cmd {
	return nil
}

// Update is the bubbletea package ELM architecture specific functions
func (f *favourite) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case favouriteMsg:
		// Reload favourite items when changes triggered
		items, err := loadItems(f.tuiConfig)
		if err != nil {
			return f, sendStatusUpdate(err.Error())
		}
		f.list.SetItems(items)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.back):
			return f, tea.Batch(
				sendViewportUpdate(MainView),
				sendStatusUpdate(""),
			)

		case key.Matches(msg, f.keys.delete):
			i, ok := f.list.SelectedItem().(cmdItem)
			if ok {
				// Delete favourite
				err := f.tuiConfig.Repository.DeleteFavourite(i.title)
				if err != nil {
					return f, sendStatusUpdate(err.Error())
				}
				return f, tea.Batch(
					sendStatusUpdate("favourite deleted"),
					sendFavouriteUpdate(),
				)
			}

		case key.Matches(msg, f.keys.launch):
			i, ok := f.list.SelectedItem().(cmdItem)
			if ok {
				// Execute the command
				cmd := exec.Command("cmd", "/C", i.wtCmd)
				if err := cmd.Run(); err != nil {
					return f, sendStatusUpdate(err.Error())
				}

				// Add command history to database
				err := f.tuiConfig.Repository.InsertHistory(i.wtCmd, strings.Split(i.cmds, ","))
				if err != nil {
					return f, sendStatusUpdate(err.Error())
				}

				return f, tea.Quit
			}
			return f, tea.Quit
		}
	}

	var cmd tea.Cmd
	f.list, cmd = f.list.Update(msg)
	return f, cmd
}

// View is the bubbletea package ELM architecture specific functions
func (f *favourite) View() string {
	f.list.SetSize(f.width, f.height)
	return f.list.View()
}
