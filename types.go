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

// SelectedMsg is emitted when the user selects a path.
type SelectedMsg struct {
	Path string
}

// CancelledMsg is emitted when the user cancels the picker.
type CancelledMsg struct{}

// Options configures a new picker instance.
type Options struct {
	Title       string
	InitialPath string
	Mode        Mode
	Layout      Layout
	Width       int
	Height      int
}

// Internal messages for async directory reads.
type dirReadResultMsg struct {
	path    string // which directory was read
	entries []DirEntry
	err     error
	gen     int // generation counter to discard stale results
}
