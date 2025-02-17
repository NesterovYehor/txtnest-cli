package tui

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/NesterovYehor/txtnest-cli/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// pbState represents the current view state.
type pbState int

const (
	pbListState pbState = iota
	pbEditState
)

// PasteBrowser manages the UI for listing and editing pastes.
type PasteBrowser struct {
	state      pbState
	list       *PastesList
	editor     *PasteEditor
	selectedID string
	client     *api.ApiClient
}

// NewPasteBrowser creates a new PasteBrowser instance.
func NewPasteBrowser(client *api.ApiClient, pastes []models.Metadata) PasteBrowser {
	return PasteBrowser{
		state:  pbListState,
		list:   NewPastesList(pastes),
		client: client,
	}
}

// Init returns the initial command based on the current state.
func (pb PasteBrowser) Init() tea.Cmd {
	if pb.state == pbListState {
		return pb.list.Init()
	} else if pb.state == pbEditState && pb.editor != nil {
		return pb.editor.Init()
	}
	return nil
}

// Update handles messages and switches between modes.
func (pb PasteBrowser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch pb.state {
	case pbListState:
		newList, cmd := pb.list.Update(msg)
		pb.list = newList.(*PastesList)
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "enter" || keyMsg.String() == " " {
				selected := pb.list.Selected()
				if selected == nil {
					return pb, cmd
				}
				pb.selectedID = selected.Key
				paste, err := pb.client.FetchPaste(selected.Key)
				if err != nil {
					fmt.Println("Error fetching paste:", err)
					return pb, tea.Quit
				}
				pb.editor = NewPasteEditor(paste.Content)
				pb.state = pbEditState
				return pb, pb.editor.Init()
			}
		}
		return pb, cmd

	case pbEditState:
		newEditor, cmd := pb.editor.Update(msg)
		pb.editor = newEditor.(*PasteEditor)
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "ctrl+s":
				// In a real application, send pb.editor.Content() to your API.
				fmt.Printf("Saving paste %s with content:\n%s\n", pb.selectedID, pb.editor.Content())
				pb.state = pbListState
				return pb, nil
			case "esc":
				// Cancel editing and return to the list view.
				pb.state = pbListState
				return pb, nil
			}
		}
		return pb, cmd
	}
	return pb, nil
}

// View renders the appropriate view based on the state.
func (pb PasteBrowser) View() string {
	switch pb.state {
	case pbListState:
		return pb.list.View()
	case pbEditState:
		return pb.editor.View()
	}
	return ""
}
