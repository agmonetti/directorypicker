package directorypicker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const sampleValidYAML = `
name: "Nord Frost"
colors:
  accent: "#88c0d0"
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
  border: "#434c5e"
  border_focus: "#88c0d0"
  error: "#bf616a"
borders:
  popup: "rounded"
  dir_list: "double"
icons:
  cursor: "❯ "
  dir_icon: "📁 "
  file_icon: "📄 "
selection:
  foreground: "#2e3440"
  background: "#88c0d0"
  bold: true
  underline: false
`

func TestParseThemeYAML_Valid(t *testing.T) {
	styles, err := ParseThemeYAML([]byte(sampleValidYAML))
	if err != nil {
		t.Fatalf("expected valid YAML to parse, got: %v", err)
	}

	if styles.Title.GetForeground() != lipgloss.Color("#88c0d0") {
		t.Errorf("expected Title foreground #88c0d0, got %v", styles.Title.GetForeground())
	}
	if styles.Normal.GetForeground() != lipgloss.Color("#eceff4") {
		t.Errorf("expected Normal foreground #eceff4, got %v", styles.Normal.GetForeground())
	}
	if styles.Selected.GetForeground() != lipgloss.Color("#2e3440") {
		t.Errorf("expected Selected foreground #2e3440, got %v", styles.Selected.GetForeground())
	}
	if styles.Selected.GetBackground() != lipgloss.Color("#88c0d0") {
		t.Errorf("expected Selected background #88c0d0, got %v", styles.Selected.GetBackground())
	}
	if styles.Cursor != "❯ " {
		t.Errorf("expected Cursor '❯ ', got %q", styles.Cursor)
	}
	if styles.DirIcon != "📁 " {
		t.Errorf("expected DirIcon '📁 ', got %q", styles.DirIcon)
	}
}

func TestParseThemeYAML_JSONCompatibility(t *testing.T) {
	jsonTheme := `{
		"name": "JSON Dracula",
		"colors": {
			"accent": "#bd93f9",
			"text": "#f8f8f2",
			"muted": "#6272a4",
			"dim": "#44475a",
			"error": "#ff5555"
		},
		"borders": {
			"popup": "rounded",
			"dir_list": "normal"
		},
		"icons": {
			"cursor": "> "
		},
		"selection": {
			"foreground": "#282a36",
			"background": "#bd93f9"
		}
	}`

	styles, err := ParseThemeYAML([]byte(jsonTheme))
	if err != nil {
		t.Fatalf("expected JSON theme to parse via YAML loader, got error: %v", err)
	}

	if styles.Title.GetForeground() != lipgloss.Color("#bd93f9") {
		t.Errorf("expected Title #bd93f9, got %v", styles.Title.GetForeground())
	}
}

func TestParseThemeYAML_StrictValidation_MissingRequiredColors(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		expectedErr string
	}{
		{
			name: "missing accent",
			yaml: `
colors:
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
  error: "#bf616a"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "colors.accent is required",
		},
		{
			name: "missing text",
			yaml: `
colors:
  accent: "#88c0d0"
  muted: "#4c566a"
  dim: "#3b4252"
  error: "#bf616a"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "colors.text is required",
		},
		{
			name: "missing muted",
			yaml: `
colors:
  accent: "#88c0d0"
  text: "#eceff4"
  dim: "#3b4252"
  error: "#bf616a"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "colors.muted is required",
		},
		{
			name: "missing dim",
			yaml: `
colors:
  accent: "#88c0d0"
  text: "#eceff4"
  muted: "#4c566a"
  error: "#bf616a"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "colors.dim is required",
		},
		{
			name: "missing error",
			yaml: `
colors:
  accent: "#88c0d0"
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "colors.error is required",
		},
		{
			name: "missing selection foreground",
			yaml: `
colors:
  accent: "#88c0d0"
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
  error: "#bf616a"
selection:
  background: "#88c0d0"
`,
			expectedErr: "selection.foreground is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseThemeYAML([]byte(tc.yaml))
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.expectedErr)
			}
			if !strings.Contains(err.Error(), tc.expectedErr) {
				t.Errorf("expected error to contain %q, got %q", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestParseThemeYAML_StrictValidation_InvalidColors(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		expectedErr string
	}{
		{
			name: "invalid hex color",
			yaml: `
colors:
  accent: "#nothex"
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
  error: "#bf616a"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "invalid hex color \"#nothex\"",
		},
		{
			name: "invalid ANSI color number",
			yaml: `
colors:
  accent: "300"
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
  error: "#bf616a"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "invalid ANSI color number 300",
		},
		{
			name: "invalid color string",
			yaml: `
colors:
  accent: "my-custom-blue"
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
  error: "#bf616a"
selection:
  foreground: "#2e3440"
`,
			expectedErr: "invalid color value \"my-custom-blue\"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseThemeYAML([]byte(tc.yaml))
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.expectedErr)
			}
			if !strings.Contains(err.Error(), tc.expectedErr) {
				t.Errorf("expected error to contain %q, got %q", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestParseThemeYAML_StrictValidation_InvalidBorder(t *testing.T) {
	yamlWithBadBorder := `
colors:
  accent: "#88c0d0"
  text: "#eceff4"
  muted: "#4c566a"
  dim: "#3b4252"
  error: "#bf616a"
borders:
  popup: "wavy-lines"
selection:
  foreground: "#2e3440"
`

	_, err := ParseThemeYAML([]byte(yamlWithBadBorder))
	if err == nil {
		t.Fatal("expected error for invalid border, got nil")
	}
	if !strings.Contains(err.Error(), "unknown border style \"wavy-lines\"") {
		t.Errorf("expected error to mention wavy-lines, got: %v", err)
	}
}

func TestLoadThemeFile(t *testing.T) {
	// Test file not found
	_, err := LoadThemeFile("/non/existent/path/theme.yaml")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}

	// Test valid file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_theme.yaml")
	if err := os.WriteFile(filePath, []byte(sampleValidYAML), 0644); err != nil {
		t.Fatalf("failed to write temp theme file: %v", err)
	}

	styles, err := LoadThemeFile(filePath)
	if err != nil {
		t.Fatalf("expected LoadThemeFile to succeed, got: %v", err)
	}

	if styles.Title.GetForeground() != lipgloss.Color("#88c0d0") {
		t.Errorf("expected title foreground #88c0d0, got %v", styles.Title.GetForeground())
	}
}
