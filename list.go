package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item string

// FilterValue is the value we use when filtering against this item when
// we're filtering the list.
func (i item) FilterValue() string {
	return ""
}

type itemDelegate struct{}

var itemStyle = lipgloss.NewStyle().PaddingLeft(5)
var selectedStyle = lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("170"))

// Render renders the item's view.
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	render := itemStyle.Render(str)
	if index == m.Index() {
		render = selectedStyle.Render("> " + str)

	}
	fmt.Fprintf(w, render)
}

// Height is the height of the list item.
func (d itemDelegate) Height() int {
	return 1
}

// Spacing is the size of the horizontal gap between list items in cells.
func (d itemDelegate) Spacing() int {
	return 0
}

// Update is the update loop for items. All messages in the list's update
// loop will pass through here except when the user is setting a filter.
// Use this method to perform item-level updates appropriate to this
// delegate.
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
