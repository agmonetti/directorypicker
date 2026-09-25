package directorypicker

// Mode determines whether the picker selects files or directories.
type Mode uint8

const (
	SelectDirectory Mode = iota
	SelectFile
)

// Layout determines how the picker renders itself.
type Layout uint8

const (
	Centered Layout = iota
	Embedded
)

// SortOrder controls how entries are sorted in the picker listing.
type SortOrder uint8

const (
	SortDirsFirst  SortOrder = iota // Directories first, then files A→Z (default)
	SortByNameAsc                   // All entries A→Z
	SortByNameDesc                  // All entries Z→A
	SortFilesFirst                  // Files first, then directories A→Z
	SortBySizeDesc                  // Largest files first (directories grouped at top)
)

// SelectedMsg is emitted when the user selects a path.
type SelectedMsg struct {
	Path string
}

// CancelledMsg is emitted when the user cancels the picker.
type CancelledMsg struct{}

// Options configures a new picker instance.
type Options struct {
	Title             string
	InitialPath       string
	Mode              Mode
	Layout            Layout
	Width             int
	Height            int
	Styles            Styles
	KeyMap            KeyMap
	ShowHidden        bool
	MaxListHeight     int      // 0 = auto (uses 50% of terminal height, capped at 20)
	AllowedExtensions []string // nil or empty = all files; e.g. []string{".go", ".md"}
	SortOrder         SortOrder
}

// Internal messages for async directory reads.
type dirReadResultMsg struct {
	path    string // which directory was read
	entries []DirEntry
	err     error
	gen     int // generation counter to discard stale results
}
