package main

import (
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//if has args, use them as query and open browser
	if len(os.Args) > 1 {
		query := strings.Join(os.Args[1:], "+")
		openInBroswer(query)
		return
	}

	m := NewModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
