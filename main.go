package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	config = Config{
		ConfigPath: "",
		ConfigFile: "",
		Debug:      false,
	}
)

type (
	Config struct {
		ConfigPath string
		ConfigFile string
		Debug      bool
	}

	Styles struct {
		app lipgloss.Style
	}

	cmdErr struct{ err error }

	// Model for Bubble Tea
	model struct {
		views        map[string]*view
		colorSet     colorSet
		width        int
		height       int
		view         string
		quitting     bool
		Style        Styles
		windowsuser  string
		subscription SubscriptionInfo
	}
)

func (m *model) Resize() {

	m.Style.app = m.Style.app.Width(m.width - m.Style.app.GetHorizontalBorderSize())
	m.Style.app = m.Style.app.Height(m.height - m.Style.app.GetVerticalBorderSize())

	for _, v := range m.views {
		v.Resize(m.Style.app.GetWidth(), m.Style.app.GetHeight())
	}

}

// Init initializes the Bubble Tea program
func (m model) Init() tea.Cmd {
	CallClear()
	return nil
}

func (m *model) SwitchView(v string, refresh bool) {
	if v != m.view {
		if _, ok :=m.views[v];!ok || refresh  {
			switch v {
			case "machines":
				m.views[v] = m.newMachineView(m.subscription)
			case "subscriptions":
				m.views[v] = m.newSubscriptionView(true)
			case "setuser":
				m.views[v] = m.newSetUserView()
			}
		}
		if _, ok := m.views[v]; ok {
			m.view = v
		}
	}
	m.Resize()
}

// Update handles view changes and user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Resize()
	case switchView:
		m.SwitchView(msg.view, msg.refresh)
	case setSub:
		if msg.setDefault {
			msg.sub.setDefault()
		}
		m.subscription = msg.sub
	case setUser:
		m.windowsuser = msg.user
		SaveConfig("windowsuser", msg.user)
	case sshMachine:
		return m, tea.Sequence(teaCallClear(), msg.SSHToMachine(m.windowsuser))
	case tea.KeyMsg:
		return m, m.views[m.view].Update(msg)
	}
	return m, nil
}

// View renders the TUI
func (m model) View() string {
	//var elements []string
	if m.quitting {
		CallClear()
		if config.Debug {
			SaveConfig("styledump", fmt.Sprintf("%+v\n%d\n%d\n", m.Style, m.height, m.width))
		}
		return ""
	}
	if m.view != "" {
		return m.Style.app.Render(m.views[m.view].View())
	}
	return "Unknown view " + m.view
}

func teaCallClear() tea.Cmd {
	var clear string = "clear"
	if runtime.GOOS == "windows" {
		clear = "cls"
	}
	cmd := osCommand(clear)
	cmd.Stdout = os.Stdout
	return tea.ExecProcess(&cmd, func(err error) tea.Msg { return cmdErr{err} })
}

func newModel() model {
	m := model{
		Style: Styles{
			app: lipgloss.NewStyle(),
		},
		colorSet: colorSet{
			bgColor:            lipgloss.Color("#101010"),
			selectedItemColor:  lipgloss.Color("#7ae582"),
			listItemColor:      lipgloss.Color("#d8ddae"),
			borderColor:        lipgloss.Color("#7ae582"),
			helpColor:          lipgloss.Color("#e2914f"),
			gradientStartColor: lipgloss.Color("#2c7255"),
			gradientEndColor:   lipgloss.Color("#52b788"),
		},
	}
	m.windowsuser = LoadConfig("windowsuser")
	if m.windowsuser == "" {
		m.windowsuser = "administrator"
		SaveConfig("windowsuser", "administrator")
	}
	m.subscription = GetDefaultSubscription()
	m.views = make(map[string]*view)
	m.SwitchView("machines", false)
	return m
}

func main() {

	CallClear()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to determine user home directory: " + err.Error())
	}
	// Initialize the Config struct
	config.ConfigPath = filepath.Join(homeDir, ".ezazssh")
	config.ConfigFile = filepath.Join(config.ConfigPath, "config")

	m := newModel()
	m.SwitchView("machines", true)
	p := tea.NewProgram(m) //tea.WithAltScreen(),

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting TUI: %v", err)
	}
}
