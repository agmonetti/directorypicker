package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/agmonetti/directorypicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ── styles ──────────────────────────────────────────────────────────────────

var (
	colorAccent = lipgloss.Color("#58a6ff")
	colorOk     = lipgloss.Color("#3fb950")
	colorDim    = lipgloss.Color("#6e7681")
	colorWarn   = lipgloss.Color("#d29922")
	colorErr    = lipgloss.Color("#f85149")

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	styleSubHeader = lipgloss.NewStyle().
			Foreground(colorDim)

	styleConfirmed = lipgloss.NewStyle().
			Foreground(colorOk).
			Bold(true)

	styleDivider = lipgloss.NewStyle().
			Foreground(colorDim)

	styleCode = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 2).
			MarginTop(1)

	styleSuccess = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorOk)

	styleCmd = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14"))
)

// ── header ───────────────────────────────────────────────────────────────────

func printHeader() {
	fmt.Println(styleHeader.Render("⚡ directorypicker — Interactive Setup Wizard"))
	fmt.Println(styleSubHeader.Render("Configure your file/folder picker and generate ready-to-use Go code."))
	fmt.Println()
}

// ── step helpers ─────────────────────────────────────────────────────────────

func printStepDivider(n, total int, title string) {
	bar := strings.Repeat("─", 44-len(title))
	fmt.Println(styleDivider.Render(fmt.Sprintf("  Step %d / %d  ·  %s  %s", n, total, title, bar)))
}

func printConfirmed(label, value string) {
	fmt.Printf("  %s  %s\n\n",
		styleConfirmed.Render("✔ "+label+":"),
		lipgloss.NewStyle().Foreground(colorAccent).Render(value),
	)
}

func runStep(n, total int, title string, field huh.Field) {
	printStepDivider(n, total, title)
	form := huh.NewForm(huh.NewGroup(field)).WithTheme(huh.ThemeCharm())
	if err := form.Run(); err != nil {
		if err == huh.ErrUserAborted {
			fmt.Println(styleSubHeader.Render("\n  Cancelled by user."))
			os.Exit(0)
		}
		fmt.Printf("  Error: %v\n", err)
		os.Exit(1)
	}
}

// ── main ─────────────────────────────────────────────────────────────────────

func main() {
	printHeader()

	// Initial step count (16 steps in file mode, 15 in dir mode)
	step := 1
	totalSteps := 16

	// ── Step 1: Selection mode ───────────────────────────────────────────────
	var modeChoice string = "file"
	runStep(step, totalSteps, "Selection mode",
		huh.NewSelect[string]().
			Title("What should the user be able to pick?").
			Description("Determines what type of filesystem entry the user can select.").
			Options(
				huh.NewOption("Files  (e.g. open or load a file)", "file"),
				huh.NewOption("Directories  (e.g. save destination or output folder)", "dir"),
			).
			Value(&modeChoice),
	)
	step++

	modeLabel := "Files"
	if modeChoice == "dir" {
		modeLabel = "Directories"
		totalSteps = 15 // extension filter is skipped in dir mode
	}
	printConfirmed("Selection mode", modeLabel)

	// ── Step 2: UI layout ────────────────────────────────────────────────────
	var layoutChoice string = "centered"
	runStep(step, totalSteps, "UI layout",
		huh.NewSelect[string]().
			Title("How should the picker appear on screen?").
			Description("Controls how the component is positioned in the terminal.").
			Options(
				huh.NewOption("Centered floating modal  (dialog box in the middle)", "centered"),
				huh.NewOption("Embedded  (inline in your own layout)", "embedded"),
			).
			Value(&layoutChoice),
	)
	step++

	layoutLabel := "Centered floating modal"
	if layoutChoice == "embedded" {
		layoutLabel = "Embedded"
	}
	printConfirmed("UI layout", layoutLabel)

	// ── Step 3: Color palette ────────────────────────────────────────────────
	var paletteChoice string = "tokyonight"
	runStep(step, totalSteps, "Color palette",
		huh.NewSelect[string]().
			Title("Which color palette do you prefer?").
			Description("Sets the base colors for borders, text, breadcrumbs, and highlights.").
			Options(
				huh.NewOption("Tokyo Night     — deep night blue, purple & cyan highlights", "tokyonight"),
				huh.NewOption("Catppuccin Mocha — soft pastel dark aesthetic (mauve & sky)", "catppuccin"),
				huh.NewOption("Dracula          — gothic dark with vivid purple, pink & green", "dracula"),
				huh.NewOption("Cyberpunk Neon   — electric neon yellow, hot pink & cyan", "cyberpunk"),
				huh.NewOption("Nord Frost       — arctic cool icy cyan & polar blues", "nord"),
				huh.NewOption("Matrix Green     — retro phosphor green terminal look", "matrix"),
				huh.NewOption("Azure Classic    — Microsoft Azure blue (#0078D4)", "azure"),
				huh.NewOption("Monochrome       — low-contrast subtle grayscale", "minimal"),
			).
			Value(&paletteChoice),
	)
	step++

	paletteNames := map[string]string{
		"tokyonight": "Tokyo Night (night blue, purple & cyan)",
		"catppuccin": "Catppuccin Mocha (mauve & sky)",
		"dracula":    "Dracula (purple, pink & green)",
		"cyberpunk":  "Cyberpunk Neon (electric yellow, pink & cyan)",
		"nord":       "Nord Frost (icy cyan & polar blues)",
		"matrix":     "Matrix Green (retro phosphor terminal)",
		"azure":      "Azure Classic (Microsoft blue)",
		"minimal":    "Monochrome (subtle grayscale)",
	}
	printConfirmed("Color palette", paletteNames[paletteChoice])

	// ── Step 4: Border style ─────────────────────────────────────────────────
	var borderChoice string = "rounded"
	runStep(step, totalSteps, "Border style",
		huh.NewSelect[string]().
			Title("What style of borders should frame the picker?").
			Options(
				huh.NewOption("Rounded borders      ╭───╮  (modern & sleek)", "rounded"),
				huh.NewOption("Double-line borders  ╔═══╗  (classic retro dialog)", "double"),
				huh.NewOption("Thick / Heavy        ┏━━━┓  (bold & striking)", "thick"),
				huh.NewOption("Single thin line     ┌───┐  (clean & minimal)", "normal"),
				huh.NewOption("Borderless / Hidden         (frameless clean spacing)", "hidden"),
			).
			Value(&borderChoice),
	)
	step++

	borderNames := map[string]string{
		"rounded": "Rounded (╭───╮)",
		"double":  "Double-line (╔═══╗)",
		"thick":   "Thick / Heavy (┏━━━┓)",
		"normal":  "Single thin line (┌───┐)",
		"hidden":  "Borderless (frameless)",
	}
	printConfirmed("Border style", borderNames[borderChoice])

	// ── Step 5: Icon set ─────────────────────────────────────────────────────
	var iconChoice string = "emoji"
	runStep(step, totalSteps, "Icon set",
		huh.NewSelect[string]().
			Title("Which icon set should represent entries?").
			Options(
				huh.NewOption("Emojis          📁 Folders  /  📄 Files  (expressive & friendly)", "emoji"),
				huh.NewOption("Nerd Font        Folders  /  󰈔 Files  (for developer terminals)", "nerd"),
				huh.NewOption("Modern Minimal  ❯ Folders  /  · Files  (clean & sharp)", "modern"),
				huh.NewOption("Classic         ▸ Folders  /  · Files  (standard Unicode)", "classic"),
			).
			Value(&iconChoice),
	)
	step++

	iconNames := map[string]string{
		"emoji":   "Emojis (📁 / 📄)",
		"nerd":    "Nerd Font ( / 󰈔)",
		"modern":  "Modern Minimal (❯ / ·)",
		"classic": "Classic (▸ / ·)",
	}
	printConfirmed("Icon set", iconNames[iconChoice])

	// ── Step 6: Active row highlight ─────────────────────────────────────────
	var highlightChoice string = "banner"
	runStep(step, totalSteps, "Active row highlight",
		huh.NewSelect[string]().
			Title("How should the currently selected entry be highlighted?").
			Options(
				huh.NewOption("Solid colored banner  (high-contrast filled row, like fzf / lazygit)", "banner"),
				huh.NewOption("Colored text & arrow  (clean text highlight, no background fill)", "text"),
				huh.NewOption("Bold & Underline      (minimalist terminal convention)", "underline"),
			).
			Value(&highlightChoice),
	)
	step++

	highlightNames := map[string]string{
		"banner":    "Solid colored banner (filled background row)",
		"text":      "Colored text with cursor arrow",
		"underline": "Bold & Underlined",
	}
	printConfirmed("Row highlight", highlightNames[highlightChoice])

	// ── Step 7: File extension filter (file mode only) ───────────────────────
	var allowedExtensions []string
	extFilterLabel := "All files"
	if modeChoice == "file" {
		var extChoices []string
		runStep(step, totalSteps, "File extension filter",
			huh.NewMultiSelect[string]().
				Title("Filter by file type").
				Description("Leave empty or choose All (*) to allow every file. Space toggles, Enter continues.").
				Options(
					huh.NewOption("All files  (*)", "*"),
					huh.NewOption("Go source files  (.go)", ".go"),
					huh.NewOption("Markdown docs    (.md)", ".md"),
					huh.NewOption("JSON files       (.json)", ".json"),
					huh.NewOption("YAML configs     (.yaml / .yml)", ".yaml"),
					huh.NewOption("Images           (.png .jpg .gif .svg .webp)", ".png"),
					huh.NewOption("Plain text       (.txt)", ".txt"),
					huh.NewOption("Shell scripts    (.sh .bash)", ".sh"),
				).
				Value(&extChoices),
		)
		step++

		hasAll := false
		for _, e := range extChoices {
			if e == "*" {
				hasAll = true
				break
			}
		}
		if !hasAll && len(extChoices) > 0 {
			for _, e := range extChoices {
				switch e {
				case ".yaml":
					allowedExtensions = append(allowedExtensions, ".yaml", ".yml")
				case ".png":
					allowedExtensions = append(allowedExtensions, ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp")
				case ".sh":
					allowedExtensions = append(allowedExtensions, ".sh", ".bash")
				default:
					allowedExtensions = append(allowedExtensions, e)
				}
			}
			extFilterLabel = strings.Join(extChoices, ", ")
		}
		printConfirmed("File extension filter", extFilterLabel)
	}

	// ── Step 8: Max visible items ─────────────────────────────────────────────
	var maxItemsChoice string = "0"
	runStep(step, totalSteps, "Max visible items",
		huh.NewSelect[string]().
			Title("How many items should the list show at once?").
			Description("Controls the visible height of the file/folder list.").
			Options(
				huh.NewOption("Auto  (50% of terminal height, max 20)", "0"),
				huh.NewOption("10 rows", "10"),
				huh.NewOption("15 rows", "15"),
				huh.NewOption("20 rows", "20"),
				huh.NewOption("30 rows", "30"),
			).
			Value(&maxItemsChoice),
	)
	step++

	maxItemsLabel := "Auto"
	if maxItemsChoice != "0" {
		maxItemsLabel = maxItemsChoice + " rows"
	}
	printConfirmed("Max visible items", maxItemsLabel)

	// ── Step 9: Window title ──────────────────────────────────────────────────
	var titleInput string
	runStep(step, totalSteps, "Window title",
		huh.NewInput().
			Title("Window title  (optional)").
			Description("Leave empty to use the automatic title based on the selection mode.").
			Placeholder("e.g. \"Select configuration file\"").
			Value(&titleInput),
	)
	step++

	titleLabel := titleInput
	if titleLabel == "" {
		if modeChoice == "file" {
			titleLabel = "(auto) Select file"
		} else {
			titleLabel = "(auto) Select folder"
		}
	}
	printConfirmed("Window title", titleLabel)

	// ── Step 10: Starting directory ───────────────────────────────────────────
	var initialPath string = "."
	runStep(step, totalSteps, "Starting directory",
		huh.NewInput().
			Title("Starting directory").
			Description("The directory where navigation begins. Use '.' for the current directory.").
			Value(&initialPath),
	)
	step++

	printConfirmed("Starting directory", initialPath)

	// ── Step 11: Show hidden files ────────────────────────────────────────────
	var showHidden bool = false
	runStep(step, totalSteps, "Show hidden files",
		huh.NewConfirm().
			Title("Show hidden files by default?").
			Description("Files and folders starting with '.' — can always be toggled at runtime with '.'.").
			Value(&showHidden),
	)
	step++

	hiddenLabel := "No"
	if showHidden {
		hiddenLabel = "Yes"
	}
	printConfirmed("Show hidden files", hiddenLabel)

	// ── Step 12: Vim-style keybindings ────────────────────────────────────────
	var vimKeys bool = true
	runStep(step, totalSteps, "Vim keybindings",
		huh.NewConfirm().
			Title("Enable Vim-style navigation shortcuts?").
			Description("Adds h/j/k/l as aliases for arrow keys in addition to the standard keys.").
			Value(&vimKeys),
	)
	step++

	vimLabel := "Yes (h/j/k/l + arrows)"
	if !vimKeys {
		vimLabel = "No (arrows only)"
	}
	printConfirmed("Vim keybindings", vimLabel)

	// ── Step 13: Sort order ───────────────────────────────────────────────────
	var sortOrderChoice string = "dirs_first"
	runStep(step, totalSteps, "Sort order",
		huh.NewSelect[string]().
			Title("How should entries be sorted in the list?").
			Options(
				huh.NewOption("Directories first, then files A→Z  (default)", "dirs_first"),
				huh.NewOption("Name A→Z  (all entries mixed)", "name_asc"),
				huh.NewOption("Name Z→A  (all entries mixed)", "name_desc"),
				huh.NewOption("Files first, then directories A→Z", "files_first"),
				huh.NewOption("Size — largest first  (directories at top)", "size_desc"),
			).
			Value(&sortOrderChoice),
	)
	step++

	sortLabels := map[string]string{
		"dirs_first":  "Directories first",
		"name_asc":    "Name A→Z",
		"name_desc":   "Name Z→A",
		"files_first": "Files first",
		"size_desc":   "Size (largest first)",
	}
	printConfirmed("Sort order", sortLabels[sortOrderChoice])

	// ── Assemble styles from choices ─────────────────────────────────────────
	pickerStyles, stylesCodeBlock := assembleStyles(paletteChoice, borderChoice, iconChoice, highlightChoice)

	// ── Step 14: Live preview ─────────────────────────────────────────────────
	var wantPreview bool = true
	runStep(step, totalSteps, "Live preview",
		huh.NewConfirm().
			Title("Preview the picker interactively before saving?").
			Description("Opens an interactive test in the terminal with your exact visual settings.").
			Value(&wantPreview),
	)
	step++

	previewLabel := "Yes"
	if !wantPreview {
		previewLabel = "No"
	}
	printConfirmed("Live preview", previewLabel)

	// ── Build picker options ──────────────────────────────────────────────────
	pickerMode := directorypicker.SelectFile
	modeStr := "directorypicker.SelectFile"
	if modeChoice == "dir" {
		pickerMode = directorypicker.SelectDirectory
		modeStr = "directorypicker.SelectDirectory"
	}

	pickerLayout := directorypicker.Centered
	layoutStr := "directorypicker.Centered"
	if layoutChoice == "embedded" {
		pickerLayout = directorypicker.Embedded
		layoutStr = "directorypicker.Embedded"
	}

	keyMap := directorypicker.DefaultKeyMap()
	if !vimKeys {
		keyMap.Up = []string{"up"}
		keyMap.Down = []string{"down"}
		keyMap.Left = []string{"left"}
		keyMap.Right = []string{"right"}
		keyMap.Home = []string{"home"}
		keyMap.End = []string{"end"}
	}

	finalTitle := titleInput
	if finalTitle == "" {
		if pickerMode == directorypicker.SelectFile {
			finalTitle = "Select file"
		} else {
			finalTitle = "Select folder"
		}
	}

	maxListHeight := 0
	if maxItemsChoice != "0" {
		fmt.Sscanf(maxItemsChoice, "%d", &maxListHeight)
	}

	sortOrderMap := map[string]directorypicker.SortOrder{
		"dirs_first":  directorypicker.SortDirsFirst,
		"name_asc":    directorypicker.SortByNameAsc,
		"name_desc":   directorypicker.SortByNameDesc,
		"files_first": directorypicker.SortFilesFirst,
		"size_desc":   directorypicker.SortBySizeDesc,
	}
	sortOrderStrMap := map[string]string{
		"dirs_first":  "directorypicker.SortDirsFirst",
		"name_asc":    "directorypicker.SortByNameAsc",
		"name_desc":   "directorypicker.SortByNameDesc",
		"files_first": "directorypicker.SortFilesFirst",
		"size_desc":   "directorypicker.SortBySizeDesc",
	}

	opts := directorypicker.Options{
		Title:             finalTitle,
		InitialPath:       initialPath,
		Mode:              pickerMode,
		Layout:            pickerLayout,
		Styles:            pickerStyles,
		KeyMap:            keyMap,
		ShowHidden:        showHidden,
		MaxListHeight:     maxListHeight,
		AllowedExtensions: allowedExtensions,
		SortOrder:         sortOrderMap[sortOrderChoice],
	}

	// ── Run live preview if requested ─────────────────────────────────────────
	if wantPreview {
		fmt.Println(styleSubHeader.Render("  Launching interactive preview... (Enter to select · Esc to exit)"))
		previewModel := previewAppModel{picker: directorypicker.New(opts)}
		p := tea.NewProgram(previewModel, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("  Warning: could not run preview: %v\n", err)
		}
		fmt.Println(styleSuccess.Render("  ✔ Preview complete."))
		fmt.Println()
	}

	// ── Step 15: Output action ────────────────────────────────────────────────
	var actionChoice string = "both"
	runStep(step, totalSteps, "Output action",
		huh.NewSelect[string]().
			Title("How do you want to receive this configuration?").
			Options(
				huh.NewOption("Both — create file and print snippet", "both"),
				huh.NewOption("Create 'picker_component.go' ready to use in this directory", "file"),
				huh.NewOption("Print code snippet only to copy & paste", "snippet"),
			).
			Value(&actionChoice),
	)
	step++

	actionLabels := map[string]string{
		"both":    "Both (file + snippet)",
		"file":    "Create file (picker_component.go)",
		"snippet": "Print snippet only",
	}
	printConfirmed("Output action", actionLabels[actionChoice])

	// ── Step 16: Code format (1b: runnable app vs helper function) ───────────
	var formatChoice string = "runnable"
	runStep(step, totalSteps, "Code format",
		huh.NewSelect[string]().
			Title("What format of Go code do you want?").
			Description("Choose whether you want a full standalone application or just the helper function.").
			Options(
				huh.NewOption("Full runnable application  (main.go with tea.NewProgram, runs with 'go run')", "runnable"),
				huh.NewOption("Helper function only       (NewConfiguredPicker() to insert into your code)", "helper"),
			).
			Value(&formatChoice),
	)

	formatLabels := map[string]string{
		"runnable": "Full runnable application (main.go)",
		"helper":   "Helper function only (NewConfiguredPicker)",
	}
	printConfirmed("Code format", formatLabels[formatChoice])

	// ── Generate code ─────────────────────────────────────────────────────────
	detectedPkg := detectPackageName(".")
	if detectedPkg == "directorypicker" {
		// Prevent creating an internal circular import when run inside the directorypicker repo
		detectedPkg = "main"
	}

	generatedCode := generateIntegrationCode(
		detectedPkg,
		formatChoice == "runnable",
		modeStr, layoutStr, stylesCodeBlock,
		sortOrderStrMap[sortOrderChoice],
		finalTitle, initialPath,
		showHidden, vimKeys,
		maxListHeight, allowedExtensions,
	)

	if actionChoice == "file" || actionChoice == "both" {
		targetFile := "picker_component.go"
		err := os.WriteFile(targetFile, []byte(generatedCode), 0644)
		if err != nil {
			fmt.Printf("  ❌ Error writing %s: %v\n", targetFile, err)
		} else {
			fmt.Println(styleSuccess.Render(fmt.Sprintf("\n  ✔ File created: %s", targetFile)))
			if formatChoice == "runnable" {
				fmt.Println("  Test it immediately with:")
				fmt.Println(styleCmd.Render("    go run picker_component.go"))
			}
		}
	}

	if actionChoice == "snippet" || actionChoice == "both" {
		fmt.Println("\n  Integration snippet ready to use:")
		fmt.Println(styleCode.Render(generatedCode))
	}

	fmt.Println("\n  If you haven't added the dependency to your project yet, run:")
	fmt.Println(styleCmd.Render("    go get github.com/agmonetti/directorypicker\n"))
}

// ── style assembler ───────────────────────────────────────────────────────────

func assembleStyles(palette, borderChoice, iconChoice, highlightChoice string) (directorypicker.Styles, string) {
	var s directorypicker.Styles
	var baseCall string

	switch palette {
	case "tokyonight":
		s = directorypicker.TokyoNightStyles()
		baseCall = "directorypicker.TokyoNightStyles()"
	case "catppuccin":
		s = directorypicker.CatppuccinStyles()
		baseCall = "directorypicker.CatppuccinStyles()"
	case "dracula":
		s = directorypicker.DraculaStyles()
		baseCall = "directorypicker.DraculaStyles()"
	case "cyberpunk":
		s = directorypicker.CyberpunkStyles()
		baseCall = "directorypicker.CyberpunkStyles()"
	case "nord":
		s = directorypicker.NordStyles()
		baseCall = "directorypicker.NordStyles()"
	case "matrix":
		s = directorypicker.MatrixStyles()
		baseCall = "directorypicker.MatrixStyles()"
	case "minimal":
		s = directorypicker.MinimalStyles()
		baseCall = "directorypicker.MinimalStyles()"
	default:
		s = directorypicker.DefaultStyles()
		baseCall = "directorypicker.DefaultStyles()"
	}

	// Border
	var borderFunc string
	var b lipgloss.Border
	switch borderChoice {
	case "double":
		b = lipgloss.DoubleBorder()
		borderFunc = "lipgloss.DoubleBorder()"
	case "thick":
		b = lipgloss.ThickBorder()
		borderFunc = "lipgloss.ThickBorder()"
	case "normal":
		b = lipgloss.NormalBorder()
		borderFunc = "lipgloss.NormalBorder()"
	case "hidden":
		b = lipgloss.HiddenBorder()
		borderFunc = "lipgloss.HiddenBorder()"
	default:
		b = lipgloss.RoundedBorder()
		borderFunc = "lipgloss.RoundedBorder()"
	}
	s.Popup = s.Popup.Border(b)
	s.DirList = s.DirList.Border(b)

	// Icons
	var dirIcon, fileIcon, cursor string
	switch iconChoice {
	case "emoji":
		dirIcon = "📁 "
		fileIcon = "📄 "
		cursor = "❯ "
	case "nerd":
		dirIcon = " "
		fileIcon = "󰈔 "
		cursor = "❯ "
	case "modern":
		dirIcon = "❯ "
		fileIcon = "· "
		cursor = "▶ "
	default:
		dirIcon = "▸ "
		fileIcon = "· "
		cursor = "▶ "
	}
	s.DirIcon = dirIcon
	s.FileIcon = fileIcon
	s.Cursor = cursor

	// Highlight style
	var accentHex, bgHex string
	switch palette {
	case "tokyonight":
		accentHex = "#bb9af7"
		bgHex = "#7aa2f7"
	case "catppuccin":
		accentHex = "#cba6f7"
		bgHex = "#cba6f7"
	case "dracula":
		accentHex = "#ff79c6"
		bgHex = "#bd93f9"
	case "cyberpunk":
		accentHex = "#ffe600"
		bgHex = "#00ff9f"
	case "nord":
		accentHex = "#88c0d0"
		bgHex = "#88c0d0"
	case "matrix":
		accentHex = "#00ff66"
		bgHex = "#00ff66"
	case "minimal":
		accentHex = "255"
		bgHex = "240"
	default: // azure
		accentHex = "#0078D4"
		bgHex = "#0078D4"
	}

	var highlightCode string
	switch highlightChoice {
	case "text":
		s.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color(accentHex)).Bold(true)
		highlightCode = fmt.Sprintf("\tstyles.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color(%q)).Bold(true)", accentHex)
	case "underline":
		s.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color(accentHex)).Bold(true).Underline(true)
		highlightCode = fmt.Sprintf("\tstyles.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color(%q)).Bold(true).Underline(true)", accentHex)
	default: // banner
		s.Selected = lipgloss.NewStyle().Background(lipgloss.Color(bgHex)).Foreground(lipgloss.Color("#000000")).Bold(true)
		highlightCode = fmt.Sprintf("\tstyles.Selected = lipgloss.NewStyle().Background(lipgloss.Color(%q)).Foreground(lipgloss.Color(\"#000000\")).Bold(true)", bgHex)
	}

	codeBlock := fmt.Sprintf(`func buildCustomStyles() directorypicker.Styles {
	styles := %s
	styles.Popup = styles.Popup.Border(%s)
	styles.DirList = styles.DirList.Border(%s)
	styles.DirIcon = %q
	styles.FileIcon = %q
	styles.Cursor = %q
%s
	return styles
}()`, baseCall, borderFunc, borderFunc, dirIcon, fileIcon, cursor, highlightCode)

	return s, codeBlock
}

// ── preview app model ─────────────────────────────────────────────────────────

type previewAppModel struct {
	picker directorypicker.Model
}

func (m previewAppModel) Init() tea.Cmd {
	return m.picker.Init()
}

func (m previewAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case directorypicker.SelectedMsg, directorypicker.CancelledMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

func (m previewAppModel) View() string {
	return m.picker.View()
}

// ── helpers ───────────────────────────────────────────────────────────────────

func detectPackageName(dir string) string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go") && fi.Name() != "picker_component.go"
	}, parser.PackageClauseOnly)

	if err == nil {
		for name := range pkgs {
			if name != "" {
				return name
			}
		}
	}
	return "main"
}

func generateIntegrationCode(
	pkgName string,
	isRunnable bool,
	modeStr, layoutStr, stylesCodeBlock, sortOrderStr,
	title, initialPath string,
	showHidden, vimKeys bool,
	maxListHeight int,
	allowedExts []string,
) string {
	vimKeySetup := ""
	if !vimKeys {
		vimKeySetup = `
	km := directorypicker.DefaultKeyMap()
	km.Up = []string{"up"}
	km.Down = []string{"down"}
	km.Left = []string{"left"}
	km.Right = []string{"right"}
`
	}

	keyMapAssign := ""
	if !vimKeys {
		keyMapAssign = "\n\t\tKeyMap:            km,"
	}

	maxHeightLine := ""
	if maxListHeight > 0 {
		maxHeightLine = fmt.Sprintf("\n\t\tMaxListHeight:     %d,", maxListHeight)
	}

	extLine := ""
	if len(allowedExts) > 0 {
		quoted := make([]string, len(allowedExts))
		for i, e := range allowedExts {
			quoted[i] = fmt.Sprintf("%q", e)
		}
		extLine = fmt.Sprintf("\n\t\tAllowedExtensions: []string{%s},", strings.Join(quoted, ", "))
	}

	var sb strings.Builder

	if isRunnable {
		sb.WriteString(`package main

import (
	"fmt"
	"os"

	"github.com/agmonetti/directorypicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
`)
	} else {
		sb.WriteString(fmt.Sprintf(`package %s

import (
	"github.com/agmonetti/directorypicker"
	"github.com/charmbracelet/lipgloss"
)
`, pkgName))
	}

	sb.WriteString(fmt.Sprintf(`
// NewConfiguredPicker returns a ready-to-use directorypicker model with your saved settings.
func NewConfiguredPicker() directorypicker.Model {%s
	return directorypicker.New(directorypicker.Options{
		Title:             "%s",
		InitialPath:       "%s",
		Mode:              %s,
		Layout:            %s,
		Styles:            %s,
		SortOrder:         %s,%s%s%s
		ShowHidden:        %t,
	})
}
`, vimKeySetup, title, initialPath, modeStr, layoutStr, stylesCodeBlock, sortOrderStr,
		keyMapAssign, maxHeightLine, extLine, showHidden))

	if isRunnable {
		sb.WriteString(`
type demoApp struct {
	picker   directorypicker.Model
	selected string
}

func (a demoApp) Init() tea.Cmd {
	return a.picker.Init()
}

func (a demoApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case directorypicker.SelectedMsg:
		a.selected = msg.Path
		return a, tea.Quit
	case directorypicker.CancelledMsg:
		a.selected = "(cancelled)"
		return a, tea.Quit
	}
	var cmd tea.Cmd
	a.picker, cmd = a.picker.Update(msg)
	return a, cmd
}

func (a demoApp) View() string {
	return a.picker.View()
}

func main() {
	app := demoApp{picker: NewConfiguredPicker()}
	p := tea.NewProgram(app, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(demoApp); ok && m.selected != "" {
		fmt.Printf("Selected: %s\n", m.selected)
	}
}
`)
	}

	return sb.String()
}
