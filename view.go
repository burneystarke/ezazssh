package main

/*
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"*/
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	view struct {
		title    title
		content  content
		Style    lipgloss.Style
		colorSet colorSet
	}
	colorSet struct {
		bgColor            lipgloss.Color
		selectedItemColor  lipgloss.Color
		listItemColor      lipgloss.Color
		borderColor        lipgloss.Color
		helpColor          lipgloss.Color
		gradientStartColor lipgloss.Color
		gradientEndColor   lipgloss.Color
	}
	content interface {
		View() string
		Update(tea.Msg) tea.Cmd
		Resize(width int, height int)
	}
	switchView struct {
		view    string
		refresh bool
	}
)

func (view *view) View() string {
	return lipgloss.JoinVertical(
		0,
		view.title.View(),
		view.content.View(),
	)
}

func (view *view) Update(msg tea.Msg) tea.Cmd {
	return view.content.Update(msg)
}

func (view *view) Resize(width int, height int) {
	view.Style = view.Style.Width(width).Height(height)
	view.title.Resize(width, height)
	view.content.Resize(width, height-view.title.style.GetHeight()-view.title.style.GetVerticalBorderSize())
}
