package ui

import (
	"context"
	"math"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	page = int
	size = int
)

const (
	createPage page = iota
	readPage
	inputKeyPage
	accountPage
	aboutPage
)

const (
	undersized size = iota
	small
	medium
	large
)

type model struct {
	ready           bool
	page            page
	isMenuActive    bool
	state           state
	context         context.Context
	rednder         *lipgloss.Renderer
	viewportWidth   int
	viewportHeight  int
	widthContainer  int
	heightContainer int
	widthContent    int
	heightContent   int
	size            size
	accessToken     string
	viewport        viewport.Model
	hasScroll       bool
}
type state struct {
	menu        menuState
	createPaste createPasteState
	footer      footerState
	account     accountState
	keyInput    keyInputState
	readPaste   readPasteState
}

func NewModel(render *lipgloss.Renderer) (tea.Model, error) {
	ctx := context.Background()
	result := model{
		context:      ctx,
		page:         createPage,
		rednder:      render,
		isMenuActive: true,
		state: state{
			menu:    menuState{},
			account: accountState{},
			footer: footerState{
				commands: []footerCommand{
					{key: "↑/↓", value: "navigate"},
					{key: "enter", value: "select"},
				},
			},
		},
	}
	return result, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) SwitchPage(page page) model {
	m.page = page
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportHeight = msg.Height
		m.viewportWidth = msg.Width
		switch {
		case m.viewportWidth < 30 || m.viewportHeight < 10:
			m.size = undersized
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 60:
			m.size = small
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 80:
			m.size = medium
			m.widthContainer = 60
			m.heightContainer = int(math.Min(float64(msg.Height), 30))
		default:
			m.size = large
			m.widthContainer = 80
			m.heightContainer = int(math.Min(float64(msg.Height), 40))
		}
		m.widthContent = m.widthContainer - 4
		m = m.updateViewport()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlQ:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.isMenuActive {
				return m, tea.Quit
			}
			if !m.isMenuActive {
				m.isMenuActive = true
				m.updateFooter()
			}

		case tea.KeyEnter:
			if m.isMenuActive {
				m.isMenuActive = false
				m.updateFooter()
			} else {
				return m.upadateContent(msg)
			}
		default:
			return m.upadateContent(msg)
		}
	}
	return m, nil
}

func (m *model) updateFooter() {
	if m.isMenuActive {
		m.state.footer.commands = []footerCommand{
			{key: "↑/↓", value: "navigate"},
			{key: "Enter", value: "select"},
		}
	} else {
		switch m.page {
		case createPage:
			m.state.footer.commands = []footerCommand{
				{key: "↑/↓", value: "navigate"},
				{key: "Enter", value: "select"},
				{key: "Tab", value: "menu/text-area"},
				{key: "Esc", value: "exit"},
				{key: "Ctrl+Q", value: "quit"},
			}
		case readPage:
			m.state.footer.commands = []footerCommand{
				{key: "↑/↓", value: "up/down"},
				{key: "Esc", value: "back"},
				{key: "Ctrl+Q", value: "quit"},
			}
		case accountPage:
			m.state.footer.commands = []footerCommand{
				{key: "Enter", value: "next"},
				{key: "Tab", value: "change tabs"},
				{key: "Esc", value: "back"},
				{key: "Ctrl+Q", value: "quit"},
			}
		case inputKeyPage:
			m.state.footer.commands = []footerCommand{
				{key: "Enter", value: "confirm key"},
				{key: "Esc", value: "back"},
			}
		default:
			m.state.footer.commands = []footerCommand{}
		}
	}
}

func (m model) upadateContent(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isMenuActive {
		return m.UpdateMenu(msg)
	} else {
		switch m.page {
		case createPage:
			return m.UpdateCreatePaste(msg)
		case inputKeyPage:
			return m.InputKeyUpadte(msg)
		case readPage:
			return m.ReadPasteUpdate(msg)
		case accountPage:
			return m.UpdateAccount(msg)
		}
	}
	return m, nil
}

func (m model) updateViewport() model {
	verticalMargin := int(float64(m.viewportHeight) * 0.1)
	width := m.widthContainer - 4
	m.heightContent = m.heightContainer - verticalMargin
	if !m.ready {
		m.viewport = viewport.New(width, m.heightContent)
		m.viewport.HighPerformanceRendering = false
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.GotoTop()
	}

	m.state.menu.list = list.New(choices, list.NewDefaultDelegate(), m.widthContent/3, m.heightContainer)

	m.state.createPaste = createPasteState{
		inTextAreaMode: true,
		input: &createPasteInput{
			expirationMenu: expirationMenu,
			textArea:       textarea.New(),
		},
	}
	m.state.createPaste.input.textArea.SetHeight(m.heightContent)
	m.state.createPaste.input.textArea.SetWidth(m.widthContent)
	m.state.createPaste.input.textArea.ShowLineNumbers = false
	m.state.createPaste.input.textArea.Prompt = ""
	m.state.createPaste.input.textArea.Focus()
	m.state.createPaste.input.textArea.CharLimit = 0
	m.state.account = m.newAccountState()
	m.state.keyInput = m.newKeyInputModel()
	m.state.account.forms[m.state.account.activeTab].Init()
	return m
}

func (m model) View() string {
	if m.isMenuActive {
		return lipgloss.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Bottom,
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					m.MenuView(),
					m.getContent(),
				),
				m.FooterView(),
			),
		)
	} else {
		return lipgloss.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Bottom,
				m.getContent(), // Only show content when menu is inactive
				m.FooterView(),
			),
		)
	}
}

func (m model) getContent() string {
	page := ""
	switch m.page {
	case createPage:
		page = m.CreatePasteView()
	case inputKeyPage:
		page = m.InputKeyView()
	case readPage:
		page = m.ReadPasteView()
	case accountPage:
		page = m.ViewAccount()
	}
	return page
}
