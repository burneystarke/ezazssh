package main

import (
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

type (
	titleSet struct {
		smalltitle  string
		mediumtitle string
		largetitle  string
		title       string
	}

	title struct {
		titleSet titleSet
		style    lipgloss.Style
		colorSet colorSet
	}
)

func (title *title) Resize(width int, height int) {
	title.style = title.style.Width(width - title.style.GetHorizontalBorderSize())

	tbwidth, tbheight, _ := getStringsDims(title.titleSet.largetitle)
	title.titleSet.title = title.titleSet.largetitle
	if title.style.GetWidth()-title.style.GetHorizontalPadding() < tbwidth {
		tbwidth, tbheight, _ = getStringsDims(title.titleSet.mediumtitle)
		title.titleSet.title = title.titleSet.mediumtitle

	}
	if title.style.GetWidth()-title.style.GetHorizontalPadding() < tbwidth {
		tbwidth, tbheight, _ = getStringsDims(title.titleSet.smalltitle)
		title.titleSet.title = title.titleSet.smalltitle

	}
	if title.style.GetWidth()-title.style.GetHorizontalPadding() < tbwidth {
		title.style.MaxWidth(tbwidth)
	}
	title.style = title.style.Height(tbheight)
}

func (title *title) View() string {
	return title.style.Render(rainbows(lipgloss.NewStyle().Background(title.colorSet.bgColor), title.titleSet.title, title.colorSet.gradientStartColor, title.colorSet.gradientEndColor))
}

func NewTitle(colorSet colorSet, titleSet titleSet) title {
	return title{
		style: lipgloss.NewStyle().Padding(1).Align(lipgloss.Center).Border(
			lipgloss.Border{
				Top:         "─",
				Bottom:      "─",
				Left:        "│",
				Right:       "│",
				TopLeft:     "╭",
				TopRight:    "╮",
				BottomLeft:  "├",
				BottomRight: "┤"},
			true).
			Height(6).
			Background(colorSet.bgColor).
			BorderBackground(colorSet.bgColor).
			Foreground(colorSet.borderColor).
			BorderForeground(colorSet.borderColor),
		colorSet: colorSet,
		titleSet: titleSet,
	}

}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range []rune(s) {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str

}

func getStringsDims(s string) (int, int, []string) {
	rows := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	var max = -1
	for _, r := range rows {
		r = strings.TrimSpace(r)
		if utf8.RuneCountInString(r) > max {
			max = utf8.RuneCountInString(r)
		}
	}
	return max, len(rows), rows
}

func rainbows(base lipgloss.Style, s string, startcolor lipgloss.Color, endcolor lipgloss.Color) string {
	var ret []string

	max, _, rows := getStringsDims(s)
	var tilt int = 1
	//var tilt = int(math.Abs(float64(max/len(rows) - len(rows))))

	fullcolors := gamut.Blends(startcolor, endcolor, (len(rows)*tilt+max)*10)

	for i, r := range rows {
		sc, _ := colorful.MakeColor(fullcolors[(len(rows)*tilt-1-i*tilt)*10])
		ec, _ := colorful.MakeColor(fullcolors[(len(rows)*tilt-1-i*tilt+max)*10])
		colors := gamut.Blends(sc, ec, utf8.RuneCountInString(r))
		ret = append(ret, rainbow(base, r, colors))
	}

	return strings.Join(ret, "\n")

}
