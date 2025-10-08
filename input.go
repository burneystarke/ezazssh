package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	myInputKeyMap struct {
		set  key.Binding
		back key.Binding
	}
	myInput struct {
		input      textinput.Model
		inputStyle lipgloss.Style
		style      lipgloss.Style
		help       lipgloss.Style
	}
	setUser struct{ user string }
)

func (myInput *myInput) View() string {
	return myInput.style.Render(
		lipgloss.JoinVertical(0, myInput.inputStyle.Render(myInput.input.View()),
			myInput.help.Render(),
		))
}

func (myInput *myInput) Resize(width int, height int) {
	_, hheight, _ := getStringsDims(myInput.help.String())
	myInput.style = myInput.style.Height(height - myInput.style.GetVerticalFrameSize() - 2).Width(width - myInput.style.GetHorizontalBorderSize())

	myInput.help = myInput.help.Width(myInput.style.GetWidth() - myInput.help.GetHorizontalFrameSize()).Height(hheight)
	myInput.inputStyle = myInput.inputStyle.Width(myInput.style.GetWidth() - myInput.inputStyle.GetHorizontalFrameSize()).Height(myInput.style.GetHeight() - myInput.style.GetVerticalPadding() - hheight)

}

func (myInput *myInput) Update(msg tea.Msg) tea.Cmd {
	newInput, cmd := myInput.input.Update(msg)
	myInput.input = newInput
	return cmd
}

func newmyInputPrompt(colorSet colorSet, width int, height int, defaultvalue string, prompt string, help string) myInput {
	myInput := myInput{
		input: textinput.New(),
		style: lipgloss.NewStyle().Width(width).Height(height).Padding(1, 1, 0, 1).Border(lipgloss.RoundedBorder(), false, true, true, true).
			Background(colorSet.bgColor).
			BorderBackground(colorSet.bgColor).
			Foreground(colorSet.borderColor).
			BorderForeground(colorSet.borderColor),
		inputStyle: lipgloss.NewStyle().Background(colorSet.bgColor).
			BorderBackground(colorSet.bgColor).
			Foreground(colorSet.borderColor).
			BorderForeground(colorSet.borderColor),
		help: lipgloss.NewStyle().Background(colorSet.bgColor).
			BorderBackground(colorSet.bgColor).
			Foreground(colorSet.helpColor).
			BorderForeground(colorSet.borderColor).SetString(help)}
	myInput.input.Cursor.Style = lipgloss.NewStyle().Foreground(colorSet.selectedItemColor).Background(colorSet.bgColor)
	myInput.input.CharLimit = 36
	myInput.input.Placeholder = defaultvalue
	myInput.input.Focus()
	myInput.input.Prompt = prompt
	myInput.input.PromptStyle = lipgloss.NewStyle().Foreground(colorSet.listItemColor).Background(colorSet.bgColor)
	myInput.input.TextStyle = lipgloss.NewStyle().Foreground(colorSet.selectedItemColor).Background(colorSet.bgColor)
	myInput.input.Cursor.Blink = true
	myInput.input.Cursor.Focus()
	myInput.input.Cursor.TextStyle = myInput.input.TextStyle
	myInput.input.SetValue(defaultvalue)
	return myInput
}
