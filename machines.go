package main

import (
	"fmt"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	MachineInfo struct {
		Name           string
		ResourceGroup  string
		SubscriptionID string
		OS             string
		IsArc          bool
	}
	machineDelegateKeyMap struct {
		Select      key.Binding
		SwitchView  key.Binding
		WindowsUser key.Binding
	}
	machineList struct {
		myList
	}
	sshMachine struct {
		MachineInfo
	}
)

var (
	mdKeyMap = machineDelegateKeyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("[↵ Enter]", "SSH to machine"),
		),
		SwitchView: key.NewBinding(
			key.WithKeys("s", "S"),
			key.WithHelp("[S]", "Change Subscription"),
		),
		WindowsUser: key.NewBinding(
			key.WithKeys("w", "W"),
			key.WithHelp("[W]", "Set Windows User"),
		),
	}
	machineTitle = titleSet{
		mediumtitle: `██╗   ██╗███╗   ███╗███████╗
██║   ██║████╗ ████║██╔════╝
██║   ██║██╔████╔██║███████╗
╚██╗ ██╔╝██║╚██╔╝██║╚════██║
 ╚████╔╝ ██║ ╚═╝ ██║███████║
  ╚═══╝  ╚═╝     ╚═╝╚══════╝`,
		largetitle: `███████╗███████╗ █████╗ ███████╗███████╗███████╗██╗  ██╗         ██╗   ██╗███╗   ███╗███████╗
██╔════╝╚══███╔╝██╔══██╗╚══███╔╝██╔════╝██╔════╝██║  ██║   ███╗  ██║   ██║████╗ ████║██╔════╝
█████╗    ███╔╝ ███████║  ███╔╝ ███████╗███████╗███████║  █████╗ ██║   ██║██╔████╔██║███████╗
██╔══╝   ███╔╝  ██╔══██║ ███╔╝  ╚════██║╚════██║██╔══██║  ╚███╔╝ ╚██╗ ██╔╝██║╚██╔╝██║╚════██║
███████╗███████╗██║  ██║███████╗███████║███████║██║  ██║   ╚══╝   ╚████╔╝ ██║ ╚═╝ ██║███████║
╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝           ╚═══╝  ╚═╝     ╚═╝╚══════╝`,
		smalltitle: "VMS",
	}
)

func (i MachineInfo) Title() string { return i.Name }
func (i MachineInfo) Description() string {
	var arc string = "Azure"
	if i.IsArc {
		arc = "Arc"
	}
	return fmt.Sprintf("%s • %s", i.ResourceGroup, arc)
}
func (i MachineInfo) FilterValue() string {
	var arc string = "Azure"
	if i.IsArc {
		arc = "Arc"
	}
	return fmt.Sprintf("%s %s %s", i.Name, i.ResourceGroup, arc)
}

func (machineList *machineList) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if machineList.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, mdKeyMap.Select):
			if i, ok := machineList.list.SelectedItem().(MachineInfo); ok {
				return tea.Sequence(func() tea.Msg { return sshMachine{i} })
			}
		case key.Matches(msg, mdKeyMap.SwitchView):
			return func() tea.Msg { return switchView{view: "subscriptions", refresh: true} }

		case key.Matches(msg, mdKeyMap.WindowsUser):
			return func() tea.Msg { return switchView{view: "setuser", refresh: false} }
		case key.Matches(msg, machineList.list.KeyMap.Quit):
			return tea.Sequence(teaCallClear(), tea.Quit)
		}
	}
	var cmd tea.Cmd
	machineList.list, cmd = machineList.list.Update(msg)
	return cmd

}

func newMachineListModel(sub SubscriptionInfo, colorSet colorSet, width int, height int) *machineList {
	items, _ := sub.getMachineList()
	listModel := newListModel(items, colorSet, width, height)
	listModel.list.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{mdKeyMap.Select, mdKeyMap.SwitchView, mdKeyMap.WindowsUser} }
	listModel.list.AdditionalShortHelpKeys = func() []key.Binding { return []key.Binding{mdKeyMap.Select, mdKeyMap.SwitchView, mdKeyMap.WindowsUser} }
	return &machineList{listModel}
}

func (m model) newMachineView(sub SubscriptionInfo) *view {
	machineList := newMachineListModel(sub, m.colorSet, m.width, m.height)
	return &view{
		title:    NewTitle(m.colorSet, machineTitle),
		content:  machineList,
		Style:    lipgloss.NewStyle(),
		colorSet: m.colorSet,
	}
}

func (s SubscriptionInfo) getMachineList() ([]list.Item, error) {

	machines, err := GetMachines(s.ID)
	if err != nil {
		return nil, err
	}
	items := make([]list.Item, len(machines)) // Preallocate the slice
	for i, machine := range machines {
		items[i] = list.Item(machine) // Assign directly
	}
	return items, nil
}

// SSHToMachine starts an interactive SSH session
func (machine MachineInfo) SSHToMachine(windowsuser string) tea.Cmd {
	var (
		machinetype string = "arc"
		extraopts   string = " -y"
		opensshopts string = "-o StrictHostKeyChecking=no -q"
	)

	if !machine.IsArc {
		machinetype = "vm"
		extraopts += " --prefer-private-ip"
	}
	if machine.OS == "windows" {
		extraopts += " --local-user " + windowsuser
	}
	var cmdstring string = fmt.Sprintf("az ssh %s --subscription %s -g %s -n %s %s", machinetype, machine.SubscriptionID, machine.ResourceGroup, machine.Name, extraopts)
	if runtime.GOOS == "windows" {
		cmdstring = "set term=xterm-color&" + cmdstring
	}
	cmdstring = fmt.Sprintf("%s -- %s", cmdstring, opensshopts)
	cmd := osCommand(cmdstring)
	return tea.ExecProcess(&cmd, func(err error) tea.Msg {
		return cmdErr{err}
	})
}
