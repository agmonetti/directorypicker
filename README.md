# directorypicker

A customizable, keyboard-driven file and directory picker component for Go TUI applications built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Supports modal dialog or embedded layouts, vim navigation, custom icon sets, sorting, file extension filtering, and built-in themes.

---

## Features

- **Keyboard-driven navigation**: Standard arrow keys + Vim shortcuts (`h`/`j`/`k`/`l`), `Home`/`End`, `PgUp`/`PgDown`.
- **Typed selection modes**: `SelectFile` or `SelectDirectory`.
- **Layout styles**: `Centered` floating modal or `Embedded` inside your existing layout.
- **Deep visual customization**: Custom borders (Rounded, Double, Thick, Normal, Hidden), colors, and custom icon sets (`DirIcon`, `FileIcon`, `Cursor`).
- **File extension filtering**: Filter listings by extensions (e.g. `[]string{".go", ".md"}`).
- **Configurable sorting**: Directories first, Name A→Z / Z→A, Files first, Size descending.
- **Hidden files**: Toggle visibility at runtime with `.` or configure default on/off.
- **Non-blocking I/O**: Directory reads run asynchronously via Bubble Tea commands (`tea.Cmd`).


---

## Interactive Setup

The fastest way to configure and scaffold `directorypicker` into your Go project is via the interactive CLI setup:

```bash
go run github.com/agmonetti/directorypicker/cmd/setup@latest
```

The wizard guides you step-by-step through:
1. **Selection mode**: Files or Directories
2. **Layout**: Centered modal popup or embedded inline
3. **Visual design**: Color palette (Tokyo Night, Catppuccin, Dracula, Cyberpunk, Nord, Matrix, Azure), border style, icon sets (Emojis, Nerd Font, Modern, Classic), and active row highlight (Banner, Text, Underline)
4. **Filtering & Sorting**: File extension filtering and configurable sort order
5. **Keybindings**: Vim-style shortcuts (`h`/`j`/`k`/`l`)
6. **Live Preview**: Interactive terminal test before saving
7. **Code Generation**: Generates a ready-to-run `main.go` or a drop-in helper function

---

## Installation

```bash
go get github.com/agmonetti/directorypicker
```

---

## Basic Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/agmonetti/directorypicker"
	tea "github.com/charmbracelet/bubbletea"
)

type appModel struct {
	picker   directorypicker.Model
	selected string
}

func (m appModel) Init() tea.Cmd {
	return m.picker.Init()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case directorypicker.SelectedMsg:
		m.selected = msg.Path
		return m, tea.Quit
	case directorypicker.CancelledMsg:
		m.selected = "(cancelled)"
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

func (m appModel) View() string {
	return m.picker.View()
}

func main() {
	picker := directorypicker.New(directorypicker.Options{
		Title:       "Select a file",
		InitialPath: ".",
		Mode:        directorypicker.SelectFile,
		Layout:      directorypicker.Centered,
		Styles:      directorypicker.TokyoNightStyles(),
	})

	p := tea.NewProgram(appModel{picker: picker}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
```

---


## Examples

| **Tokyo Night** (Modern Icons & Banner) | **Cyberpunk Neon** (Nerd Font & Thick Border) |
|:---:|:---:|
| <img src="docs/assets/tokyonight.png" width="420" alt="Tokyo Night Theme"> | <img src="docs/assets/cyberpunk.png" width="420" alt="Cyberpunk Neon Theme"> |
| **Catppuccin Mocha** (Folder Mode & Minimal) | **Matrix Green** (Retro Phosphor Terminal) |
| <img src="docs/assets/catppuccin.png" width="420" alt="Catppuccin Mocha Theme"> | <img src="docs/assets/matrix.png" width="420" alt="Matrix Green Theme"> |

---

## Guides & Documentation

- [Customization Guide](docs/CUSTOMIZATION.md): Detailed styling, custom keymaps, and options.
- [Integration Guide](docs/INTEGRATION.md): Bubble Tea lifecycle, messages, resizing, and embedding.
- [Installation Guide](docs/INSTALLATION.md): Requirements and release checklists.
- [Design Plan](PLAN.md): Architecture decisions and specification.

---

## License

[MIT](LICENSE) © [Agustin Monetti](https://github.com/agmonetti)
