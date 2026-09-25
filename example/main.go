package main

import (
	"fmt"
	"os"

	"github.com/agmonetti/directorypicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appModel struct {
	picker   directorypicker.Model
	active   bool
	selected string
	width    int
	height   int
}

func initialModel() appModel {
	opts := directorypicker.Options{
		Title:       "Select a file or folder",
		Mode:        directorypicker.SelectFile,
		Layout:      directorypicker.Centered,
		ShowHidden:  false,
		Styles:      directorypicker.DefaultStyles(),
		InitialPath: ".",
	}

	return appModel{
		picker: directorypicker.New(opts),
		active: true,
	}
}

func (m appModel) Init() tea.Cmd {
	return m.picker.Init()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		var cmd tea.Cmd
		m.picker, cmd = m.picker.Update(msg)
		return m, cmd

	case directorypicker.SelectedMsg:
		m.active = false
		m.selected = msg.Path
		return m, nil

	case directorypicker.CancelledMsg:
		m.active = false
		m.selected = "(cancelled)"
		return m, nil

	case tea.KeyMsg:
		if !m.active {
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			case "r", "enter", " ":
				m.active = true
				m.selected = ""
				return m, m.picker.Init()
			}
		}
	}

	if m.active {
		var cmd tea.Cmd
		m.picker, cmd = m.picker.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m appModel) View() string {
	if m.active {
		return m.picker.View()
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#0078D4")).
		Padding(1, 3).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255"))

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246"))

	pathStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#2899F5"))

	content := fmt.Sprintf(
		"%s\n\nResult:\n%s\n\n%s",
		titleStyle.Render("directorypicker Demo"),
		pathStyle.Render(m.selected),
		textStyle.Render("Press [Space] to pick again  •  [q] to quit"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		cardStyle.Render(content),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running example: %v\n", err)
		os.Exit(1)
	}
}
