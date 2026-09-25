package directorypicker

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the Bubble Tea model for the directory/file picker.
type Model struct {
	Title             string
	currentPath       string
	entries           []DirEntry
	cursor            int
	err               error
	loading           bool
	width             int
	height            int
	mode              Mode
	layout            Layout
	styles            Styles
	keyMap            KeyMap
	showHidden        bool
	maxListHeight     int
	allowedExtensions []string
	sortOrder         SortOrder
	readGen           int // generation counter to discard stale async reads
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

	styles := opts.Styles
	if isZeroStyles(styles) {
		styles = DefaultStyles()
	}

	keyMap := opts.KeyMap
	if len(keyMap.Up) == 0 {
		keyMap = DefaultKeyMap()
	}

	return Model{
		Title:             title,
		currentPath:       initial,
		mode:              opts.Mode,
		layout:            opts.Layout,
		width:             opts.Width,
		height:            opts.Height,
		styles:            styles,
		keyMap:            keyMap,
		showHidden:        opts.ShowHidden,
		maxListHeight:     opts.MaxListHeight,
		allowedExtensions: opts.AllowedExtensions,
		sortOrder:         opts.SortOrder,
		loading:           true,
	}
}


func isZeroStyles(s Styles) bool {
	b, _, _, _, _ := s.Popup.GetBorder()
	return b.Top == "" && !s.Title.GetBold() && s.Selected.GetForeground() == (lipgloss.NoColor{})
}

// Init returns a command that performs the initial directory read asynchronously.
func (m Model) Init() tea.Cmd {
	return m.readDirCmd()
}

// Styles returns the current styles configuration.
func (m Model) Styles() Styles {
	return m.styles
}

// KeyMap returns the current keymap configuration.
func (m Model) KeyMap() KeyMap {
	return m.keyMap
}

// ShowHidden returns whether hidden files are currently shown.
func (m Model) ShowHidden() bool {
	return m.showHidden
}

// CurrentPath returns the current directory path.
func (m Model) CurrentPath() string {
	return m.currentPath
}

// readDirCmd returns a tea.Cmd that reads the current directory in a goroutine.
func (m Model) readDirCmd() tea.Cmd {
	path := m.currentPath
	showFiles := m.mode == SelectFile
	showHidden := m.showHidden
	allowedExts := m.allowedExtensions
	sortOrder := m.sortOrder
	gen := m.readGen
	return func() tea.Msg {
		entries, err := readDir(path, showFiles, showHidden, allowedExts, sortOrder)
		return dirReadResultMsg{path: path, entries: entries, err: err, gen: gen}
	}
}
