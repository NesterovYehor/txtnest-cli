package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	loginTab tab = iota
	registerTab
)

// AccountModel represents the account screen
type accountState struct {
	activeTab tab
	header    string
	forms     map[tab]*huh.Form
	viewport  viewport.Model
}

// NewAccountModel creates a new AccountModel
func (m model) newAccountState() accountState {
	width := int(m.widthContent)
	height := int(m.widthContent)
	view := viewport.New(width, height)

	loginForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Email").Key("email"),
			huh.NewInput().Title("Password").Key("password").EchoMode(huh.EchoModePassword),
		),
	)

	registerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Key("name"),
			huh.NewInput().Title("Email").Key("email"),
			huh.NewInput().Title("Password").Key("password").EchoMode(huh.EchoModePassword),
		),
	)

	return accountState{
		activeTab: loginTab,
		header:    "LogIn",
		viewport:  view,
		forms: map[tab]*huh.Form{
			loginTab:    loginForm,
			registerTab: registerForm,
		},
	}
}

func (m model) UpdateAccount(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "tab":
			if m.state.account.activeTab == loginTab {
				m.state.account.activeTab = registerTab
				m.state.account.header = "Register"
			} else {
				m.state.account.activeTab = loginTab
				m.state.account.header = "LogIn"
			}
			return m, m.state.account.forms[m.state.account.activeTab].Init()
		}
	}

	updatedModel, cmd := m.state.account.forms[m.state.account.activeTab].Update(msg)
	if form, ok := updatedModel.(*huh.Form); ok {
		m.state.account.forms[m.state.account.activeTab] = form
	}

	return m, cmd
}

func (m model) ViewAccount() string {
	m.viewport.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Top,
			m.state.account.header,
			m.state.account.forms[m.state.account.activeTab].View(),
		),
	)

	return m.viewport.View()
}
