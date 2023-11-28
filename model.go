package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	input    string
	list     list.Model
	selected string
	quitted  bool
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m model) Init() tea.Cmd {
	return nil
}

const DefaultBroswer = "Microsoft Edge"
const GoogleSearchUrl = "https://www.google.com/search?q="

func openInBroswer(query string) {
	browser := os.Getenv("BROSWER")
	if browser == "" {
		browser = DefaultBroswer
	}
	// TODO implement linux/windows support
	cmd := exec.Command("open", "-a", browser, GoogleSearchUrl+query)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+n", "down":
			m.list.CursorDown()
		case "ctrl+p", "up":
			m.list.CursorUp()
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = string(i)
			} else {
				m.selected = m.input
			}
			openInBroswer(m.selected)
			return m, tea.Quit
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "ctrl+c":
			m.quitted = true
			return m, tea.Quit
		case "ctrl+w", "alt+backspace", "ctrl+backspace":
			splited := strings.Split(m.input, " ")
			m.input = strings.Join(splited[:len(splited)-1], " ")
		default:
			if strings.Contains(keypress, "ctrl") || strings.Contains(keypress, "alt") || strings.Contains(keypress, "shift") {
				return m, nil
			}

			// TODO improve it
			if strings.ContainsAny(keypress, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890.(,;:![{]}=+-_|/?<>'\" ") {
				m.input += msg.String()
			}
		}
	}

	items, err := GoogleSuggest(m.input)
	if err != nil {
		log.Println(err)
	}
	m.list.SetItems(items)
	return m, nil
}

var (
	appStyle      = lipgloss.NewStyle().Width(50).Padding(1, 2).Margin(1, 0, 2, 4).Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("#F4B400"))
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m model) View() string {
	if m.selected != "" {
		return quitTextStyle.Render("Opening " + m.selected + " in broswer...")
	}
	if m.quitted {
		return quitTextStyle.Render("Byeeee")
	}
	return appStyle.Render(m.list.View() + "\n󰖟 Web Search: " + m.input)
}

func NewModel() tea.Model {
	var m model
	list := list.New(nil, itemDelegate{}, 35, 15)
	list.Title = "󱍢 Google Suggestions"
	list.Styles.Title.Background(lipgloss.Color("#F4B400"))
	list.SetShowHelp(false)
	list.SetShowStatusBar(false)
	m.list = list
	return m
}
