package directorypicker

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages for the picker.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case dirReadResultMsg:
		// Discard stale results from a previous navigation
		if msg.gen != m.readGen {
			return m, nil
		}
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.entries = nil
			m.cursor = 0
			return m, nil
		}
		m.err = nil
		m.entries = msg.entries
		m.cursor = 0
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			// Only allow cancel while loading
			if keyMatches(msg.String(), m.keyMap.Cancel) {
				return m, func() tea.Msg { return CancelledMsg{} }
			}
			return m, nil
		}
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch {
	case keyMatches(key, m.keyMap.Cancel):
		return m, func() tea.Msg { return CancelledMsg{} }

	case keyMatches(key, m.keyMap.Enter):
		return m.handleEnter()

	case keyMatches(key, m.keyMap.ToggleHidden):
		m.showHidden = !m.showHidden
		m.loading = true
		m.readGen++
		return m, m.readDirCmd()

	case keyMatches(key, m.keyMap.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case keyMatches(key, m.keyMap.Down):
		if len(m.entries) > 0 && m.cursor < len(m.entries)-1 {
			m.cursor++
		}

	case keyMatches(key, m.keyMap.Right):
		if len(m.entries) > 0 && m.entries[m.cursor].IsDir {
			return m.navigateTo(m.entries[m.cursor].Path)
		}

	case keyMatches(key, m.keyMap.Left):
		parent := parentPath(m.currentPath)
		if parent != m.currentPath {
			return m.navigateTo(parent)
		}

	case keyMatches(key, m.keyMap.Home):
		m.cursor = 0

	case keyMatches(key, m.keyMap.End):
		if len(m.entries) > 0 {
			m.cursor = len(m.entries) - 1
		}

	case keyMatches(key, m.keyMap.PageUp):
		m.cursor -= 10
		if m.cursor < 0 {
			m.cursor = 0
		}

	case keyMatches(key, m.keyMap.PageDown):
		if len(m.entries) == 0 {
			m.cursor = 0
		} else {
			m.cursor += 10
			if m.cursor >= len(m.entries) {
				m.cursor = len(m.entries) - 1
			}
		}
	}
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	if m.mode == SelectFile {
		if len(m.entries) == 0 {
			// File mode + empty listing → nothing to select, do nothing
			return m, nil
		}
		selected := m.entries[m.cursor]
		if !selected.IsDir {
			return m, func() tea.Msg { return SelectedMsg{Path: selected.Path} }
		}
		// Directory in file mode → navigate into it
		return m.navigateTo(selected.Path)
	}

	// Directory mode → select current directory
	return m, func() tea.Msg { return SelectedMsg{Path: m.currentPath} }
}

func (m Model) navigateTo(path string) (Model, tea.Cmd) {
	m.currentPath = path
	m.loading = true
	m.readGen++
	return m, m.readDirCmd()
}
