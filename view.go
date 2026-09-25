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
	if m.maxListHeight > 0 {
		listH = m.maxListHeight
	}
	if listH > 20 && m.maxListHeight == 0 {
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

	// Cursor and icon setup
	cursorSym := m.styles.Cursor
	if cursorSym == "" {
		cursorSym = symCursor
	}
	if !strings.HasSuffix(cursorSym, " ") {
		cursorSym += " "
	}
	cursorPad := strings.Repeat(" ", lipgloss.Width(cursorSym))

	dirIcon := m.styles.DirIcon
	if dirIcon == "" {
		dirIcon = "▸"
	}
	if !strings.HasSuffix(dirIcon, " ") {
		dirIcon += " "
	}

	fileIcon := m.styles.FileIcon
	if fileIcon == "" {
		fileIcon = "·"
	}
	if !strings.HasSuffix(fileIcon, " ") {
		fileIcon += " "
	}

	// Build entry list
	var lines []string
	if m.loading {
		lines = append(lines, m.styles.Dim.Render("  Loading..."))
	} else if m.err != nil {
		lines = append(lines, m.styles.Error.Render("Error: "+m.err.Error()))
	} else if len(m.entries) == 0 {
		lines = append(lines, m.styles.Dim.Render("  (empty)"))
	} else {
		hasBg := m.styles.Selected.GetBackground() != (lipgloss.NoColor{})
		for i, entry := range m.entries {
			icon := dirIcon
			if !entry.IsDir {
				icon = fileIcon
			}

			if i == m.cursor {
				if hasBg {
					lineText := fmt.Sprintf("%s%s%s", cursorSym, icon, entry.Name)
					lines = append(lines, m.styles.Selected.Render(lineText))
				} else {
					lines = append(lines, fmt.Sprintf("%s%s%s", cursorSym, icon, m.styles.Selected.Render(entry.Name)))
				}
			} else {
				lines = append(lines, fmt.Sprintf("%s%s%s", cursorPad, icon, m.styles.Normal.Render(entry.Name)))
			}
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
	listPanel := m.styles.DirList.Width(innerW).Render(listContent)

	// Breadcrumb
	breadcrumb := m.styles.Breadcrumb.Render(m.buildBreadcrumb())

	// Current path
	currentPath := m.styles.Path.Render("Current: " + m.currentPath)

	// Help
	var helpText string
	if m.mode == SelectFile {
		helpText = "[↑/↓] Navigate  [→/←] Open/Parent  [Enter] Select file  [.] Hidden  [Esc] Cancel"
	} else {
		helpText = "[↑/↓] Navigate  [→/←] Open/Parent  [Enter] Select here  [.] Hidden  [Esc] Cancel"
	}
	help := m.styles.Help.Render(helpText)

	content := lipgloss.JoinVertical(lipgloss.Left,
		m.styles.Title.Render(m.Title),
		listPanel,
		"",
		breadcrumb,
		currentPath,
		"",
		help,
	)

	popup := m.styles.Popup.Width(popupW).Render(content)

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
