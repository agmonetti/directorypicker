package directorypicker

// KeyMap defines the keybindings for navigating and interacting with the picker.
type KeyMap struct {
	Up           []string
	Down         []string
	Left         []string
	Right        []string
	Enter        []string
	Cancel       []string
	ToggleHidden []string
	Home         []string
	End          []string
	PageUp       []string
	PageDown     []string
}

// DefaultKeyMap returns the default keybindings for directorypicker.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:           []string{"up", "k"},
		Down:         []string{"down", "j"},
		Left:         []string{"left", "h"},
		Right:        []string{"right", "l"},
		Enter:        []string{"enter"},
		Cancel:       []string{"esc"},
		ToggleHidden: []string{".", "H", "ctrl+h"},
		Home:         []string{"home", "g"},
		End:          []string{"end", "G"},
		PageUp:       []string{"pgup"},
		PageDown:     []string{"pgdown"},
	}
}

// keyMatches checks if a key string matches any of the configured bindings.
func keyMatches(key string, bindings []string) bool {
	for _, b := range bindings {
		if b == key {
			return true
		}
	}
	return false
}
