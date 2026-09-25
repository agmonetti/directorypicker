package directorypicker

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

var hexColorRegex = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

var validNamedColors = map[string]bool{
	"black":         true,
	"red":           true,
	"green":         true,
	"yellow":        true,
	"blue":          true,
	"magenta":       true,
	"cyan":          true,
	"white":         true,
	"brightblack":   true,
	"brightred":     true,
	"brightgreen":   true,
	"brightyellow":  true,
	"brightblue":    true,
	"brightmagenta": true,
	"brightcyan":    true,
	"brightwhite":   true,
}

var validBorderStyles = map[string]bool{
	"rounded": true,
	"double":  true,
	"thick":   true,
	"normal":  true,
	"hidden":  true,
	"none":    true,
}

// ThemeConfig represents the serializable structure of a directorypicker theme.
type ThemeConfig struct {
	Name      string         `yaml:"name"`
	Colors    ThemeColors    `yaml:"colors"`
	Borders   ThemeBorders   `yaml:"borders"`
	Icons     ThemeIcons     `yaml:"icons"`
	Selection ThemeSelection `yaml:"selection"`
}

// ThemeColors defines the color palette used by the picker components.
type ThemeColors struct {
	Accent      string `yaml:"accent"`       // Main accent color (title, breadcrumbs, highlights)
	Text        string `yaml:"text"`         // Normal body text color
	Muted       string `yaml:"muted"`        // Secondary / path text color
	Dim         string `yaml:"dim"`          // Dimmed text (help, subtle details)
	Border      string `yaml:"border"`       // Directory list border color (defaults to muted if empty)
	BorderFocus string `yaml:"border_focus"` // Popup outer border color (defaults to accent if empty)
	Error       string `yaml:"error"`        // Error text color
}

// ThemeBorders defines the border styles for container frames.
type ThemeBorders struct {
	Popup   string `yaml:"popup"`    // Border style for outer popup: rounded, double, thick, normal, hidden, none
	DirList string `yaml:"dir_list"` // Border style for directory list box
}

// ThemeIcons defines the icons and cursor used in entry listings.
type ThemeIcons struct {
	Cursor   string `yaml:"cursor"`    // Cursor prefix symbol (e.g. "> " or "❯ ")
	DirIcon  string `yaml:"dir_icon"`  // Folder icon prefix (e.g. "📁 " or " ")
	FileIcon string `yaml:"file_icon"` // File icon prefix (e.g. "📄 " or "󰈔 ")
}

// ThemeSelection defines visual highlighting for the currently active row.
type ThemeSelection struct {
	Foreground string `yaml:"foreground"` // Selected item text color
	Background string `yaml:"background"` // Optional background fill color for banner highlight
	Bold       *bool  `yaml:"bold"`       // Whether selected item text is bold (default true)
	Underline  *bool  `yaml:"underline"`  // Whether selected item text is underlined (default false)
}

// Validate performs strict validation on the ThemeConfig structure according to option 3b.
// It ensures all required keys are present and values (colors, borders) conform to accepted formats.
func (c *ThemeConfig) Validate() error {
	// Validate required colors
	if err := validateRequiredColor("colors.accent", c.Colors.Accent); err != nil {
		return err
	}
	if err := validateRequiredColor("colors.text", c.Colors.Text); err != nil {
		return err
	}
	if err := validateRequiredColor("colors.muted", c.Colors.Muted); err != nil {
		return err
	}
	if err := validateRequiredColor("colors.dim", c.Colors.Dim); err != nil {
		return err
	}
	if err := validateRequiredColor("colors.error", c.Colors.Error); err != nil {
		return err
	}

	// Validate optional colors if provided
	if c.Colors.Border != "" {
		if err := validateColorFormat("colors.border", c.Colors.Border); err != nil {
			return err
		}
	}
	if c.Colors.BorderFocus != "" {
		if err := validateColorFormat("colors.border_focus", c.Colors.BorderFocus); err != nil {
			return err
		}
	}

	// Validate selection
	if err := validateRequiredColor("selection.foreground", c.Selection.Foreground); err != nil {
		return err
	}
	if c.Selection.Background != "" {
		if err := validateColorFormat("selection.background", c.Selection.Background); err != nil {
			return err
		}
	}

	// Validate borders
	if c.Borders.Popup != "" {
		if !validBorderStyles[strings.ToLower(strings.TrimSpace(c.Borders.Popup))] {
			return fmt.Errorf("unknown border style %q for borders.popup: must be one of rounded, double, thick, normal, hidden, none", c.Borders.Popup)
		}
	}
	if c.Borders.DirList != "" {
		if !validBorderStyles[strings.ToLower(strings.TrimSpace(c.Borders.DirList))] {
			return fmt.Errorf("unknown border style %q for borders.dir_list: must be one of rounded, double, thick, normal, hidden, none", c.Borders.DirList)
		}
	}

	return nil
}

// ToStyles converts the ThemeConfig into a Lip Gloss Styles struct ready for use in directorypicker.
func (c *ThemeConfig) ToStyles() Styles {
	accent := lipgloss.Color(c.Colors.Accent)
	text := lipgloss.Color(c.Colors.Text)
	muted := lipgloss.Color(c.Colors.Muted)
	dim := lipgloss.Color(c.Colors.Dim)
	errColor := lipgloss.Color(c.Colors.Error)

	borderColor := muted
	if c.Colors.Border != "" {
		borderColor = lipgloss.Color(c.Colors.Border)
	}

	borderFocus := accent
	if c.Colors.BorderFocus != "" {
		borderFocus = lipgloss.Color(c.Colors.BorderFocus)
	}

	popupBorder := parseBorder(c.Borders.Popup, lipgloss.RoundedBorder())
	dirListBorder := parseBorder(c.Borders.DirList, lipgloss.RoundedBorder())

	// Build selected style
	selStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Selection.Foreground))
	if c.Selection.Background != "" {
		selStyle = selStyle.Background(lipgloss.Color(c.Selection.Background))
	}
	isBold := true
	if c.Selection.Bold != nil {
		isBold = *c.Selection.Bold
	}
	if isBold {
		selStyle = selStyle.Bold(true)
	}
	if c.Selection.Underline != nil && *c.Selection.Underline {
		selStyle = selStyle.Underline(true)
	}

	cursor := c.Icons.Cursor
	if cursor == "" {
		cursor = "> "
	}

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(popupBorder).
			BorderForeground(borderFocus).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(dirListBorder).
			BorderForeground(borderColor).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(muted),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(accent),
		Selected: selStyle,
		Normal: lipgloss.NewStyle().
			Foreground(text),
		Dim: lipgloss.NewStyle().
			Foreground(dim),
		Help: lipgloss.NewStyle().
			Foreground(dim),
		Error: lipgloss.NewStyle().
			Foreground(errColor),
		Cursor:   cursor,
		DirIcon:  c.Icons.DirIcon,
		FileIcon: c.Icons.FileIcon,
	}
}

// ParseThemeYAML deserializes and strictly validates a theme from YAML bytes, returning a ready-to-use Styles struct.
func ParseThemeYAML(data []byte) (Styles, error) {
	var cfg ThemeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Styles{}, fmt.Errorf("failed to parse theme YAML: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Styles{}, fmt.Errorf("theme validation error: %w", err)
	}

	return cfg.ToStyles(), nil
}

// LoadThemeConfig loads and strictly validates a ThemeConfig from a YAML file.
func LoadThemeConfig(path string) (*ThemeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file %q: %w", path, err)
	}

	var cfg ThemeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse theme YAML in %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("theme %q validation error: %w", path, err)
	}

	return &cfg, nil
}

// LoadThemeFile reads a YAML theme file from disk, strictly validates it, and returns the converted Styles.
func LoadThemeFile(path string) (Styles, error) {
	cfg, err := LoadThemeConfig(path)
	if err != nil {
		return Styles{}, err
	}
	return cfg.ToStyles(), nil
}

func validateRequiredColor(field, val string) error {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return fmt.Errorf("%s is required", field)
	}
	return validateColorFormat(field, clean)
}

func validateColorFormat(field, val string) error {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return nil
	}

	// Check hex color
	if strings.HasPrefix(clean, "#") {
		if !hexColorRegex.MatchString(clean) {
			return fmt.Errorf("invalid hex color %q for %s: must be #RGB or #RRGGBB format", val, field)
		}
		return nil
	}

	// Check ANSI number 0-255
	if n, err := strconv.Atoi(clean); err == nil {
		if n < 0 || n > 255 {
			return fmt.Errorf("invalid ANSI color number %d for %s: must be between 0 and 255", n, field)
		}
		return nil
	}

	// Check named color
	if validNamedColors[strings.ToLower(clean)] {
		return nil
	}

	return fmt.Errorf("invalid color value %q for %s: must be a hex code (e.g. #88c0d0), ANSI number (0-255), or standard ANSI color name", val, field)
}

func parseBorder(name string, defaultBorder lipgloss.Border) lipgloss.Border {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "rounded":
		return lipgloss.RoundedBorder()
	case "double":
		return lipgloss.DoubleBorder()
	case "thick":
		return lipgloss.ThickBorder()
	case "normal":
		return lipgloss.NormalBorder()
	case "hidden", "none":
		return lipgloss.HiddenBorder()
	default:
		return defaultBorder
	}
}
