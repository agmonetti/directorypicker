package directorypicker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the picker.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	popupW := m.width * 75 / 100
	if popupW > 90 {
		popupW = 90
	}
	if popupW < 30 {
		popupW = 30
	}
	if popupW > m.width-2 {
		popupW = m.width - 2
	}

	listH := m.height * 50 / 100
	if listH > 20 {
		listH = 20
	}
	if listH < 3 {
		listH = 3
	}
	if listH > m.height-8 {
		listH = m.height - 8
		if listH < 1 {
			listH = 1
		}
	}

	innerW := popupW - 6 // popup padding + border

	// Build entry list
	var lines []string
	if m.loading {
		lines = append(lines, dimStyle.Render("  Loading..."))
	} else if m.err != nil {
		lines = append(lines, errorStyle.Render("Error: "+m.err.Error()))
	} else if len(m.entries) == 0 {
		lines = append(lines, dimStyle.Render("  (empty)"))
	} else {
		for i, entry := range m.entries {
			cursor := "  "
			style := normalStyle
			if i == m.cursor {
				cursor = symCursor
				style = selectedStyle
			}
			icon := "▸" // dir
			if !entry.IsDir {
				icon = "·" // file
			}
			lines = append(lines, fmt.Sprintf("%s%s %s", cursor, icon, style.Render(entry.Name)))
		}
	}

	// Sliding window
	if len(lines) > listH {
		start := m.cursor - listH/2
		if start < 0 {
			start = 0
		}
		end := start + listH
		if end > len(lines) {
			end = len(lines)
			start = end - listH
			if start < 0 {
				start = 0
			}
		}
		lines = lines[start:end]
	}

	listContent := strings.Join(lines, "\n")
	listPanel := dirListStyle.Width(innerW).Render(listContent)

	// Breadcrumb
	breadcrumb := breadcrumbStyle.Render(m.buildBreadcrumb())

	// Current path
	currentPath := pathStyle.Render("Current: " + m.currentPath)

	// Help
	var helpText string
	if m.mode == SelectFile {
		helpText = "[↑/↓] Navigate  [→/←] Open/Parent  [Enter] Select file  [Esc] Cancel"
	} else {
		helpText = "[↑/↓] Navigate  [→/←] Open/Parent  [Enter] Select here  [Esc] Cancel"
	}
	help := helpStyle.Render(helpText)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(m.Title),
		listPanel,
		"",
		breadcrumb,
		currentPath,
		"",
		help,
	)

	popup := popupStyle.Width(popupW).Render(content)

	if m.layout == Embedded {
		return popup
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
}

func (m Model) buildBreadcrumb() string {
	home := homeDir()
	path := m.currentPath

	if path == home {
		return "Home"
	}
	if path == "/" {
		return "/"
	}

	rel, err := filepath.Rel(home, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return path
	}
	if rel == "." {
		return "Home"
	}

	parts := strings.Split(rel, string(os.PathSeparator))
	result := "Home"
	for _, part := range parts {
		result += " > " + part
	}

	if len(result) > 60 {
		result = "Home > ... > " + parts[len(parts)-1]
	}
	return result
}
