package directorypicker

import "github.com/charmbracelet/lipgloss"

// Styles defines the Lip Gloss styles for each element in the picker component.
type Styles struct {
	Popup      lipgloss.Style
	DirList    lipgloss.Style
	Title      lipgloss.Style
	Path       lipgloss.Style
	Breadcrumb lipgloss.Style
	Selected   lipgloss.Style
	Normal     lipgloss.Style
	Dim        lipgloss.Style
	Help       lipgloss.Style
	Error      lipgloss.Style
	Cursor     string
	DirIcon    string
	FileIcon   string
}

// DefaultStyles returns the standard styles for directorypicker.
func DefaultStyles() Styles {
	colorAccent := lipgloss.Color("#0078D4")
	colorText := lipgloss.Color("252")
	colorMuted := lipgloss.Color("240")
	colorDim := lipgloss.Color("238")
	colorBorder := lipgloss.Color("245")
	colorFocus := lipgloss.Color("#2899F5")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorFocus).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(colorMuted),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(colorAccent),
		Selected: lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(colorText),
		Dim: lipgloss.NewStyle().
			Foreground(colorMuted),
		Help: lipgloss.NewStyle().
			Foreground(colorDim),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")),
	}
}

// MinimalStyles returns a subtle, low-contrast monochrome style.
func MinimalStyles() Styles {
	fg := lipgloss.Color("250")
	muted := lipgloss.Color("242")
	border := lipgloss.Color("239")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(fg).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(muted),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(fg),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Underline(true),
		Normal: lipgloss.NewStyle().
			Foreground(fg),
		Dim: lipgloss.NewStyle().
			Foreground(muted),
		Help: lipgloss.NewStyle().
			Foreground(muted),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")),
	}
}

// RoundedStyles returns a modern rounded border style with vibrant highlights.
func RoundedStyles() Styles {
	accent := lipgloss.Color("#7D56F4")
	cyan := lipgloss.Color("#04B575")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("243")).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(cyan),
		Selected: lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("253")),
		Dim: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")),
	}
}

// ASCIIStyles returns a basic style using standard ASCII borders.
func ASCIIStyles() Styles {
	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1),
		Path:       lipgloss.NewStyle(),
		Breadcrumb: lipgloss.NewStyle(),
		Selected: lipgloss.NewStyle().
			Bold(true),
		Normal: lipgloss.NewStyle(),
		Dim:    lipgloss.NewStyle(),
		Help:   lipgloss.NewStyle(),
		Error:  lipgloss.NewStyle(),
	}
}

// CatppuccinStyles returns a rich pastel theme based on Catppuccin Mocha.
func CatppuccinStyles() Styles {
	mauve := lipgloss.Color("#cba6f7")
	sky := lipgloss.Color("#89dceb")
	text := lipgloss.Color("#cdd6f4")
	subtext := lipgloss.Color("#a6adc8")
	surface := lipgloss.Color("#45475a")
	border := lipgloss.Color("#585b70")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mauve).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(mauve).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(subtext),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(sky),
		Selected: lipgloss.NewStyle().
			Background(mauve).
			Foreground(lipgloss.Color("#11111b")).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(text),
		Dim: lipgloss.NewStyle().
			Foreground(surface),
		Help: lipgloss.NewStyle().
			Foreground(subtext),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")),
	}
}

// TokyoNightStyles returns a sleek Japanese dark theme with deep blues and purples.
func TokyoNightStyles() Styles {
	blue := lipgloss.Color("#7aa2f7")
	cyan := lipgloss.Color("#7dcfff")
	purple := lipgloss.Color("#bb9af7")
	text := lipgloss.Color("#c0caf5")
	muted := lipgloss.Color("#565f89")
	border := lipgloss.Color("#3b4261")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(purple).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(muted),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(cyan),
		Selected: lipgloss.NewStyle().
			Background(blue).
			Foreground(lipgloss.Color("#1a1b26")).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(text),
		Dim: lipgloss.NewStyle().
			Foreground(muted),
		Help: lipgloss.NewStyle().
			Foreground(muted),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")),
	}
}

// DraculaStyles returns a vibrant gothic palette with purples, pinks, and greens.
func DraculaStyles() Styles {
	purple := lipgloss.Color("#bd93f9")
	pink := lipgloss.Color("#ff79c6")
	green := lipgloss.Color("#50fa7b")
	text := lipgloss.Color("#f8f8f2")
	comment := lipgloss.Color("#6272a4")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(pink).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(comment).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(pink).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(comment),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(green),
		Selected: lipgloss.NewStyle().
			Background(purple).
			Foreground(lipgloss.Color("#282a36")).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(text),
		Dim: lipgloss.NewStyle().
			Foreground(comment),
		Help: lipgloss.NewStyle().
			Foreground(comment),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")),
	}
}

// CyberpunkStyles returns a high-contrast electric neon yellow/cyan/pink palette.
func CyberpunkStyles() Styles {
	neonYellow := lipgloss.Color("#ffe600")
	neonPink := lipgloss.Color("#ff007f")
	neonCyan := lipgloss.Color("#00f0ff")
	neonGreen := lipgloss.Color("#00ff9f")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(neonPink).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(neonCyan).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(neonYellow).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(neonGreen),
		Selected: lipgloss.NewStyle().
			Background(neonGreen).
			Foreground(lipgloss.Color("#000000")).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")),
		Dim: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")),
		Error: lipgloss.NewStyle().
			Foreground(neonPink),
	}
}

// NordStyles returns a cool Arctic palette of muted blues, cyan, and frost whites.
func NordStyles() Styles {
	frost := lipgloss.Color("#88c0d0")
	polarBlue := lipgloss.Color("#81a1c1")
	snowWhite := lipgloss.Color("#eceff4")
	nightDark := lipgloss.Color("#4c566a")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frost).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(nightDark).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(frost).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(polarBlue),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(snowWhite),
		Selected: lipgloss.NewStyle().
			Background(frost).
			Foreground(lipgloss.Color("#2e3440")).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(snowWhite),
		Dim: lipgloss.NewStyle().
			Foreground(nightDark),
		Help: lipgloss.NewStyle().
			Foreground(nightDark),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bf616a")),
	}
}

// MatrixStyles returns a retro terminal green-on-black phosphorescent theme.
func MatrixStyles() Styles {
	phosphor := lipgloss.Color("#00ff66")
	darkGreen := lipgloss.Color("#007733")
	border := lipgloss.Color("#00441b")

	return Styles{
		Popup: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(phosphor).
			Padding(1, 2),
		DirList: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(border).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(phosphor).
			MarginBottom(1),
		Path: lipgloss.NewStyle().
			Foreground(darkGreen),
		Breadcrumb: lipgloss.NewStyle().
			Foreground(phosphor),
		Selected: lipgloss.NewStyle().
			Background(phosphor).
			Foreground(lipgloss.Color("#001100")).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#99ff99")),
		Dim: lipgloss.NewStyle().
			Foreground(darkGreen),
		Help: lipgloss.NewStyle().
			Foreground(darkGreen),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff3333")),
	}
}

