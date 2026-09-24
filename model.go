package directorypicker

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the Bubble Tea model for the directory/file picker.
type Model struct {
	Title       string
	currentPath string
	entries     []DirEntry
	cursor      int
	err         error
	loading     bool
	width       int
	height      int
	mode        Mode
	layout      Layout
	readGen     int // generation counter to discard stale async reads
}

// New creates a picker model. Call Init() to start the first async directory read.
func New(opts Options) Model {
	home := homeDir()
	initial := home
	if opts.InitialPath != "" && pathExists(opts.InitialPath) {
		initial = opts.InitialPath
	}

	title := opts.Title
	if title == "" {
		if opts.Mode == SelectFile {
			title = "Select file"
		} else {
			title = "Select folder"
		}
	}

	return Model{
		Title:       title,
		currentPath: initial,
		mode:        opts.Mode,
		layout:      opts.Layout,
		width:       opts.Width,
		height:      opts.Height,
		loading:     true,
	}
}

// Init returns a command that performs the initial directory read asynchronously.
func (m Model) Init() tea.Cmd {
	return m.readDirCmd()
}

// readDirCmd returns a tea.Cmd that reads the current directory in a goroutine.
func (m Model) readDirCmd() tea.Cmd {
	path := m.currentPath
	showFiles := m.mode == SelectFile
	gen := m.readGen
	return func() tea.Msg {
		entries, err := readDir(path, showFiles)
		return dirReadResultMsg{path: path, entries: entries, err: err, gen: gen}
	}
}
