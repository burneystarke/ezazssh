package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	SubscriptionInfo struct {
		ID   string
		Name string
	}
	subDelegateKeyMap struct {
		setDefault key.Binding
		use        key.Binding
		back       key.Binding
	}
	subList struct {
		myList
	}
	setSub struct {
		sub        SubscriptionInfo
		setDefault bool
	}
)

var (
	sdKeyMap = subDelegateKeyMap{
		use: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("[↵ Enter]", "Use Subscription"),
		),
		back: key.NewBinding(
			key.WithKeys("b", "B"),
			key.WithHelp("[B]", "Back to VMs"),
		),
		setDefault: key.NewBinding(
			key.WithKeys("d", "D"),
			key.WithHelp("[D]", "Set Default Subscription"),
		),
	}
	subscriptionTitle = titleSet{
		mediumtitle: `███████╗██╗   ██╗██████╗ ███████╗
██╔════╝██║   ██║██╔══██╗██╔════╝
███████╗██║   ██║██████╔╝███████╗
╚════██║██║   ██║██╔══██╗╚════██║
███████║╚██████╔╝██████╔╝███████║
╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝`,
		largetitle: `███████╗███████╗ █████╗ ███████╗███████╗███████╗██╗  ██╗         ███████╗██╗   ██╗██████╗ ███████╗
██╔════╝╚══███╔╝██╔══██╗╚══███╔╝██╔════╝██╔════╝██║  ██║   ███╗  ██╔════╝██║   ██║██╔══██╗██╔════╝
█████╗    ███╔╝ ███████║  ███╔╝ ███████╗███████╗███████║  █████╗ ███████╗██║   ██║██████╔╝███████╗
██╔══╝   ███╔╝  ██╔══██║ ███╔╝  ╚════██║╚════██║██╔══██║  ╚███╔╝ ╚════██║██║   ██║██╔══██╗╚════██║
███████╗███████╗██║  ██║███████╗███████║███████║██║  ██║   ╚══╝  ███████║╚██████╔╝██████╔╝███████║
╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝         ╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝`,
		smalltitle: "SUBS",
	}
)

func (i SubscriptionInfo) Title() string       { return i.Name }
func (i SubscriptionInfo) Description() string { return i.ID }
func (i SubscriptionInfo) FilterValue() string { return fmt.Sprintf("%s %s", i.Name, i.ID) }

func (m model) newSubListModel(refresh bool) *subList {

	sub := GetDefaultSubscription()
	var items []list.Item
	if sub.Name != "" && sub.ID != "" && !refresh {
		items = []list.Item{sub}
	} else {
		subs, _ := GetSubscriptions()
		items = make([]list.Item, len(subs)) // Preallocate the slice
		for i, machine := range subs {
			items[i] = list.Item(machine) // Assign directly
		}
	}
	listModel := newListModel(items, m.colorSet, m.width, m.height)
	listModel.list.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{sdKeyMap.setDefault, sdKeyMap.use, sdKeyMap.back} }
	listModel.list.AdditionalShortHelpKeys = func() []key.Binding { return []key.Binding{sdKeyMap.setDefault, sdKeyMap.use, sdKeyMap.back} }
	listModel.list.Select(0)
	return &subList{listModel}
}

func (subList *subList) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if subList.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, sdKeyMap.use):
			if i, ok := subList.list.SelectedItem().(SubscriptionInfo); ok {
				return tea.Sequence(func() tea.Msg { return setSub{i, false} },func() tea.Msg { return switchView{view: "machines", refresh: true} })
			}
		case key.Matches(msg, sdKeyMap.setDefault):
			if i, ok := subList.list.SelectedItem().(SubscriptionInfo); ok {
				return tea.Sequence(func() tea.Msg { return setSub{i, true} },func() tea.Msg { return switchView{view: "machines", refresh: true} })
			}

		case key.Matches(msg, sdKeyMap.back):
			return func() tea.Msg { return switchView{view: "machines", refresh: false} }
		case key.Matches(msg, subList.list.KeyMap.Quit):
			return tea.Sequence(teaCallClear(), tea.Quit)
		}
	}
	var cmd tea.Cmd
	subList.list, cmd = subList.list.Update(msg)
	return cmd
}

func (m model) newSubscriptionView(refresh bool) *view {
	subList := m.newSubListModel(refresh)
	return &view{
		title:    NewTitle(m.colorSet, subscriptionTitle),
		content:  subList,
		Style:    lipgloss.NewStyle(),
		colorSet: m.colorSet,
	}
}
func (s SubscriptionInfo) setDefault() {
	SaveConfig("subscriptionId", s.ID)
	SaveConfig("subscriptionName", s.Name)
}

func GetDefaultSubscription() SubscriptionInfo {
	return SubscriptionInfo{
		ID:   LoadConfig("subscriptionId"),
		Name: LoadConfig("subscriptionName"),
	}
}
