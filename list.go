package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	myList struct {
		style         lipgloss.Style
		list          list.Model
		delegateStyle list.DefaultDelegate
	}
)

func (myList *myList) Resize(width int, height int) {

	myList.list.SetWidth(width - myList.style.GetHorizontalFrameSize())
	myList.list.SetHeight(height - myList.style.GetVerticalFrameSize() - 2)

	myList.list.Styles.HelpStyle = myList.list.Styles.HelpStyle.Width((myList.list.Width()))
	myList.list.Styles.PaginationStyle = myList.list.Styles.PaginationStyle.Width(myList.list.Width())
	myList.list.Styles.StatusBar = myList.list.Styles.StatusBar.Width(myList.list.Width())
	myList.list.Styles.NoItems = myList.list.Styles.NoItems.Width(myList.list.Width() - myList.list.Styles.NoItems.GetHorizontalFrameSize())

	myList.delegateStyle.Styles.NormalTitle = myList.delegateStyle.Styles.NormalTitle.Width(myList.list.Width() - myList.delegateStyle.Styles.NormalTitle.GetHorizontalBorderSize())

	myList.delegateStyle.Styles.NormalDesc = myList.delegateStyle.Styles.NormalDesc.Width(myList.delegateStyle.Styles.NormalTitle.GetWidth()).Height(myList.delegateStyle.Styles.NormalTitle.GetHeight())

	myList.delegateStyle.Styles.SelectedTitle = myList.delegateStyle.Styles.SelectedTitle.Width(myList.delegateStyle.Styles.NormalTitle.GetWidth() - myList.delegateStyle.Styles.SelectedTitle.GetHorizontalMargins() - 1).Height(myList.delegateStyle.Styles.NormalTitle.GetHeight())
	myList.delegateStyle.Styles.SelectedDesc = myList.delegateStyle.Styles.SelectedDesc.Width(myList.delegateStyle.Styles.NormalDesc.GetWidth() - myList.delegateStyle.Styles.SelectedDesc.GetHorizontalMargins() - 1).Height(myList.delegateStyle.Styles.NormalDesc.GetHeight())
	

	myList.delegateStyle.Styles.DimmedTitle = myList.delegateStyle.Styles.DimmedTitle.Width(myList.delegateStyle.Styles.NormalTitle.GetWidth()).Height(myList.delegateStyle.Styles.NormalTitle.GetHeight())
	myList.delegateStyle.Styles.DimmedDesc = myList.delegateStyle.Styles.DimmedDesc.Width(myList.delegateStyle.Styles.NormalDesc.GetWidth()).Height(myList.delegateStyle.Styles.NormalDesc.GetHeight())
	myList.list.SetDelegate(myList.delegateStyle)

}

func newListModel(items []list.Item, colorSet colorSet, width int, height int) myList {
	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.BorderForeground(colorSet.selectedItemColor).Foreground(colorSet.selectedItemColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor).Margin(0, 0, 0, 2)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.BorderForeground(colorSet.selectedItemColor).Foreground(colorSet.selectedItemColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor).Margin(0, 0, 0, 2)

	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.BorderForeground(colorSet.listItemColor).Foreground(colorSet.listItemColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.BorderForeground(colorSet.listItemColor).Foreground(colorSet.listItemColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)

	delegate.Styles.DimmedTitle = delegate.Styles.DimmedTitle.BorderForeground(colorSet.borderColor).Foreground(colorSet.borderColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	delegate.Styles.DimmedDesc = delegate.Styles.DimmedDesc.BorderForeground(colorSet.borderColor).Foreground(colorSet.borderColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	delegate.Styles.FilterMatch = delegate.Styles.FilterMatch.Foreground(colorSet.selectedItemColor).Background(colorSet.bgColor)

	listModel := list.New(items, delegate, width, height)
	listModel.Styles.NoItems = listModel.Styles.NoItems.BorderForeground(colorSet.borderColor).Foreground(colorSet.borderColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor).Margin(0,0,0,2)
	listModel.SetShowStatusBar(true)
	listModel.SetShowPagination(true)
	listModel.SetShowTitle(false)
	listModel.InfiniteScrolling = true

	listModel.Styles.PaginationStyle = listModel.Styles.PaginationStyle.Background(colorSet.bgColor)
	listModel.Paginator.ActiveDot = listModel.Styles.ActivePaginationDot.Background(colorSet.bgColor).String()
	listModel.Paginator.InactiveDot = listModel.Styles.InactivePaginationDot.Background(colorSet.bgColor).String()
	listModel.Styles.ActivePaginationDot = listModel.Styles.ActivePaginationDot.Background(colorSet.bgColor).SetString("•")

	listModel.Styles.StatusBar = listModel.Styles.StatusBar.Background(colorSet.bgColor).Foreground(colorSet.helpColor)
	listModel.Styles.StatusEmpty = listModel.Styles.StatusEmpty.Background(colorSet.bgColor).Foreground(colorSet.helpColor)
	listModel.Styles.HelpStyle = listModel.Styles.HelpStyle.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.Help.Styles.FullDesc = listModel.Help.Styles.FullDesc.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.Help.Styles.FullKey = listModel.Help.Styles.FullKey.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.Help.Styles.ShortKey = listModel.Help.Styles.ShortKey.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.Help.Styles.ShortDesc = listModel.Help.Styles.ShortDesc.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.Help.Styles.ShortSeparator = listModel.Help.Styles.ShortSeparator.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.Help.Styles.FullSeparator = listModel.Help.Styles.FullSeparator.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.Help.Styles.Ellipsis = listModel.Help.Styles.Ellipsis.BorderForeground(colorSet.borderColor).Foreground(colorSet.helpColor).Background(colorSet.bgColor).BorderBackground(colorSet.bgColor)
	listModel.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up"), key.WithHelp("[↑]", "Up"))
	listModel.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down"), key.WithHelp("[↓]", "Down"))
	listModel.KeyMap.NextPage = key.NewBinding(key.WithKeys("right"), key.WithHelp("[→]", "Next Page"))
	listModel.KeyMap.PrevPage = key.NewBinding(key.WithKeys("left"), key.WithHelp("[←]", "Prev Page"))
	//have to disable because setting width on selected item makes it draw new lines
	listModel.SetFilteringEnabled(false)
	/*
		listModel.KeyMap.Filter = key.NewBinding(key.WithKeys("/"), key.WithHelp("[/]", "Search"))
		listModel.KeyMap.ClearFilter = key.NewBinding(key.WithKeys("esc"), key.WithHelp("[esc]", "Clear Search"))
		listModel.KeyMap.ClearFilter.SetEnabled(false)
		listModel.KeyMap.CancelWhileFiltering = key.NewBinding(key.WithKeys("esc"), key.WithHelp("[esc]", "Cancel Search"))
		listModel.KeyMap.CancelWhileFiltering.SetEnabled(false)
		listModel.KeyMap.AcceptWhileFiltering = key.NewBinding(key.WithKeys("enter"), key.WithHelp("[↵ Enter]", "Accept Search"))
		listModel.KeyMap.AcceptWhileFiltering.SetEnabled(false)
	*/
	listModel.KeyMap.GoToStart.SetEnabled(false)
	listModel.KeyMap.GoToEnd.SetEnabled(false)

	return myList{
		list:          listModel,
		delegateStyle: delegate,
		style: lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), false, true, true, true).Padding(0).
			Background(colorSet.bgColor).
			BorderBackground(colorSet.bgColor).
			Foreground(colorSet.borderColor).
			BorderForeground(colorSet.borderColor),
	}
}

func (myList *myList) Update(msg tea.Msg) tea.Cmd {
	newlist, cmd := myList.list.Update(msg)
	myList.list = newlist
	return cmd
}

func (myList *myList) View() string {
	return myList.style.Render(myList.list.View())
}
