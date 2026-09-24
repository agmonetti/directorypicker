package directorypicker

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func tmpTree(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "subdir"), 0o755)
	os.WriteFile(filepath.Join(dir, "file.txt"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte("x"), 0o644)
	return dir
}

func initModel(t *testing.T, mode Mode) Model {
	t.Helper()
	dir := tmpTree(t)
	m := New(Options{InitialPath: dir, Mode: mode, Width: 80, Height: 24})
	// Simulate Init: run the cmd synchronously
	cmd := m.Init()
	if cmd != nil {
		msg := cmd()
		m, _ = m.Update(msg)
	}
	return m
}

func TestDirModeSelectsCurrentDir(t *testing.T) {
	m := initModel(t, SelectDirectory)
	// Enter selects the current directory
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command from Enter")
	}
	msg := cmd()
	sel, ok := msg.(SelectedMsg)
	if !ok {
		t.Fatalf("expected SelectedMsg, got %T", msg)
	}
	if sel.Path == "" {
		t.Fatal("selected path is empty")
	}
}

func TestFileModeEmptyDirNoSelection(t *testing.T) {
	dir := t.TempDir() // empty dir
	m := New(Options{InitialPath: dir, Mode: SelectFile, Width: 80, Height: 24})
	cmd := m.Init()
	if cmd != nil {
		msg := cmd()
		m, _ = m.Update(msg)
	}
	// Enter on empty listing in file mode should do nothing
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(SelectedMsg); ok {
			t.Fatal("file mode should not select on empty listing")
		}
	}
}

func TestFileModeSelectsFile(t *testing.T) {
	m := initModel(t, SelectFile)
	// Find the file entry
	fileIdx := -1
	for i, e := range m.entries {
		if !e.IsDir {
			fileIdx = i
			break
		}
	}
	if fileIdx < 0 {
		t.Fatal("no file entry found")
	}
	// Move cursor to the file
	m.cursor = fileIdx
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	msg := cmd()
	sel, ok := msg.(SelectedMsg)
	if !ok {
		t.Fatalf("expected SelectedMsg, got %T", msg)
	}
	if !filepath.IsAbs(sel.Path) {
		t.Fatalf("expected absolute path, got %s", sel.Path)
	}
}

func TestFileModeNavigatesDir(t *testing.T) {
	m := initModel(t, SelectFile)
	// First entry should be the directory (dirs sorted first)
	if len(m.entries) == 0 || !m.entries[0].IsDir {
		t.Skip("no dir entry at index 0")
	}
	m.cursor = 0
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	// Should navigate, not select
	if cmd == nil {
		t.Fatal("expected a read command from navigation")
	}
	msg := cmd()
	if _, ok := msg.(SelectedMsg); ok {
		t.Fatal("file mode enter on dir should navigate, not select")
	}
	if _, ok := msg.(dirReadResultMsg); !ok {
		t.Fatalf("expected dirReadResultMsg, got %T", msg)
	}
}

func TestCancelEmitsCancelledMsg(t *testing.T) {
	m := initModel(t, SelectDirectory)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command from Esc")
	}
	msg := cmd()
	if _, ok := msg.(CancelledMsg); !ok {
		t.Fatalf("expected CancelledMsg, got %T", msg)
	}
}

func TestCursorBoundsEmptyDir(t *testing.T) {
	dir := t.TempDir()
	m := New(Options{InitialPath: dir, Mode: SelectDirectory, Width: 80, Height: 24})
	cmd := m.Init()
	if cmd != nil {
		msg := cmd()
		m, _ = m.Update(msg)
	}
	if m.cursor != 0 {
		t.Fatalf("cursor should be 0 on empty dir, got %d", m.cursor)
	}
	// PgDown on empty should stay at 0
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	if m.cursor != 0 {
		t.Fatalf("cursor should be 0 after pgdown on empty, got %d", m.cursor)
	}
}

func TestHiddenFilesExcluded(t *testing.T) {
	m := initModel(t, SelectFile)
	for _, e := range m.entries {
		if e.Name == ".hidden" {
			t.Fatal(".hidden should be excluded by default")
		}
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := initModel(t, SelectDirectory)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if m.width != 120 || m.height != 40 {
		t.Fatalf("expected 120x40, got %dx%d", m.width, m.height)
	}
}

func TestStaleReadDiscarded(t *testing.T) {
	m := initModel(t, SelectDirectory)
	// Simulate a stale result from a previous generation
	stale := dirReadResultMsg{
		path:    "/some/old/path",
		entries: []DirEntry{{Name: "stale", Path: "/stale", IsDir: true}},
		gen:     m.readGen - 1,
	}
	m, _ = m.Update(stale)
	for _, e := range m.entries {
		if e.Name == "stale" {
			t.Fatal("stale read result should have been discarded")
		}
	}
}
