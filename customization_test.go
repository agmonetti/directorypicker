package directorypicker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestCustomStyles(t *testing.T) {
	custom := Styles{
		Title: lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		Popup: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
	}

	m := New(Options{
		Styles: custom,
		Width:  80,
		Height: 24,
	})

	if m.Styles().Title.GetForeground() != lipgloss.Color("196") {
		t.Errorf("expected Title foreground color 196, got %v", m.Styles().Title.GetForeground())
	}
}

func TestPresets(t *testing.T) {
	min := MinimalStyles()
	if min.Title.GetForeground() != lipgloss.Color("250") {
		t.Errorf("expected minimal title foreground 250, got %v", min.Title.GetForeground())
	}

	rnd := RoundedStyles()
	if rnd.Title.GetForeground() != lipgloss.Color("#7D56F4") {
		t.Errorf("expected rounded title foreground #7D56F4, got %v", rnd.Title.GetForeground())
	}

	ascii := ASCIIStyles()
	if ascii.Title.GetBold() != true {
		t.Errorf("expected ascii title to be bold")
	}

	cat := CatppuccinStyles()
	if cat.Title.GetForeground() != lipgloss.Color("#cba6f7") {
		t.Errorf("expected Catppuccin title mauve #cba6f7, got %v", cat.Title.GetForeground())
	}

	tn := TokyoNightStyles()
	if tn.Title.GetForeground() != lipgloss.Color("#bb9af7") {
		t.Errorf("expected Tokyo Night title purple #bb9af7, got %v", tn.Title.GetForeground())
	}

	drac := DraculaStyles()
	if drac.Title.GetForeground() != lipgloss.Color("#ff79c6") {
		t.Errorf("expected Dracula title pink #ff79c6, got %v", drac.Title.GetForeground())
	}

	cyber := CyberpunkStyles()
	if cyber.Title.GetForeground() != lipgloss.Color("#ffe600") {
		t.Errorf("expected Cyberpunk title yellow #ffe600, got %v", cyber.Title.GetForeground())
	}

	nord := NordStyles()
	if nord.Title.GetForeground() != lipgloss.Color("#88c0d0") {
		t.Errorf("expected Nord title frost #88c0d0, got %v", nord.Title.GetForeground())
	}

	matrix := MatrixStyles()
	if matrix.Title.GetForeground() != lipgloss.Color("#00ff66") {
		t.Errorf("expected Matrix title green #00ff66, got %v", matrix.Title.GetForeground())
	}
}

func TestCustomKeyMap(t *testing.T) {
	customKeys := DefaultKeyMap()
	customKeys.Cancel = []string{"q", "ctrl+c"}

	m := New(Options{
		KeyMap: customKeys,
	})
	m.loading = false

	// Esc should no longer cancel
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(CancelledMsg); ok {
			t.Error("expected Esc not to cancel with custom KeyMap")
		}
	}

	// 'q' should cancel
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected 'q' to emit command")
	}
	msg := cmd()
	if _, ok := msg.(CancelledMsg); !ok {
		t.Errorf("expected CancelledMsg from 'q', got %T", msg)
	}
}

func TestShowHiddenOption(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, ".hidden"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(tmp, "visible.txt"), []byte("data"), 0644)

	m := New(Options{
		InitialPath: tmp,
		Mode:        SelectFile,
		ShowHidden:  true,
	})

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected Init() to return read command")
	}
	msg := cmd()
	readResult, ok := msg.(dirReadResultMsg)
	if !ok {
		t.Fatalf("expected dirReadResultMsg, got %T", msg)
	}

	var foundHidden bool
	for _, entry := range readResult.entries {
		if entry.Name == ".hidden" {
			foundHidden = true
			break
		}
	}

	if !foundHidden {
		t.Errorf("expected .hidden to be included when ShowHidden is true")
	}
}

func TestToggleHiddenKey(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, ".hidden"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(tmp, "visible.txt"), []byte("data"), 0644)

	m := New(Options{
		InitialPath: tmp,
		Mode:        SelectFile,
		ShowHidden:  false,
	})
	m.loading = false

	if m.ShowHidden() {
		t.Fatal("expected ShowHidden to start false")
	}

	// Press '.' to toggle hidden
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'.'}})
	if !m.ShowHidden() {
		t.Error("expected ShowHidden to be true after pressing '.'")
	}
	if cmd == nil {
		t.Fatal("expected toggle hidden to trigger a read command")
	}

	msg := cmd()
	readResult, ok := msg.(dirReadResultMsg)
	if !ok {
		t.Fatalf("expected dirReadResultMsg, got %T", msg)
	}

	var foundHidden bool
	for _, entry := range readResult.entries {
		if entry.Name == ".hidden" {
			foundHidden = true
			break
		}
	}
	if !foundHidden {
		t.Errorf("expected .hidden to be loaded after toggling hidden")
	}

	// Update with read result
	m, _ = m.Update(readResult)
	if m.loading {
		t.Error("expected loading to be false after read result")
	}
}

func TestViewRendersBreadcrumbAndHelp(t *testing.T) {
	m := New(Options{
		Width:  100,
		Height: 30,
		Mode:   SelectFile,
	})
	m.loading = false
	m.entries = []DirEntry{
		{Name: "test.txt", Path: "/test.txt", IsDir: false},
	}

	output := m.View()
	if !strings.Contains(output, "Hidden") {
		t.Errorf("expected view help text to mention Hidden, got:\n%s", output)
	}
}
