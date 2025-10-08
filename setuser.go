package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	setUserInput struct {
		myInput
	}
)

var (
	suKeyMap = myInputKeyMap{
		set: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("[↵ Enter]", "Set User"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("[esc]", "Back to VMs"),
		),
	}

	setUserTitle = titleSet{
		largetitle: `███████╗███████╗ █████╗ ███████╗███████╗███████╗██╗  ██╗         ███████╗███████╗████████╗    ██╗   ██╗███████╗███████╗██████╗
██╔════╝╚══███╔╝██╔══██╗╚══███╔╝██╔════╝██╔════╝██║  ██║   ███╗  ██╔════╝██╔════╝╚══██╔══╝    ██║   ██║██╔════╝██╔════╝██╔══██╗
█████╗    ███╔╝ ███████║  ███╔╝ ███████╗███████╗███████║  █████╗ ███████╗█████╗     ██║       ██║   ██║███████╗█████╗  ██████╔╝
██╔══╝   ███╔╝  ██╔══██║ ███╔╝  ╚════██║╚════██║██╔══██║  ╚███╔╝ ╚════██║██╔══╝     ██║       ██║   ██║╚════██║██╔══╝  ██╔══██╗
███████╗███████╗██║  ██║███████╗███████║███████║██║  ██║   ╚══╝  ███████║███████╗   ██║       ╚██████╔╝███████║███████╗██║  ██║
╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝         ╚══════╝╚══════╝   ╚═╝        ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝`,
		mediumtitle: `███████╗███████╗████████╗    ██╗   ██╗███████╗███████╗██████╗ 
██╔════╝██╔════╝╚══██╔══╝    ██║   ██║██╔════╝██╔════╝██╔══██╗
███████╗█████╗     ██║       ██║   ██║███████╗█████╗  ██████╔╝
╚════██║██╔══╝     ██║       ██║   ██║╚════██║██╔══╝  ██╔══██╗
███████║███████╗   ██║       ╚██████╔╝███████║███████╗██║  ██║
╚══════╝╚══════╝   ╚═╝        ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝`,
		smalltitle: "SET USER",
	}
)

func (m model) newSetUserView() *view {
	c := m.newSetUserInputModel()
	return &view{
		title:    NewTitle(m.colorSet, setUserTitle),
		content:  c,
		Style:    lipgloss.NewStyle(),
		colorSet: m.colorSet,
	}
}
func (m model) newSetUserInputModel() *setUserInput {
	return &setUserInput{newmyInputPrompt(m.colorSet, m.width, m.height, m.windowsuser, "Enter username >", "[↵ Enter] Set Username • [esc] Back")}
}

func (su *setUserInput) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, suKeyMap.set):
			return tea.Sequence(
				func() tea.Msg { return setUser{user: su.input.Value()} },
				func() tea.Msg { return switchView{view: "machines", refresh: false} },
			)
		case key.Matches(msg, suKeyMap.back):
			return func() tea.Msg { return switchView{view: "machines", refresh: false} }
		}
	}
	var cmd tea.Cmd
	su.input, cmd = su.input.Update(msg)
	return cmd
}
