# Customization

> Status: design proposal, not an implemented API. Names and signatures below are candidates and must be reconciled with the code and tests before compatibility is promised.

## Principles

1. Visual and keyboard settings are configured **per instance**, not through global variables.
2. Defaults make the picker useful with only a few lines of setup.
3. The API uses Bubble Tea and Lip Gloss types already present in the ecosystem; it adds no configuration framework.
4. Each option has defined behavior. Avoid redundant flags, plugins, and callbacks when a value or public message is enough.
5. Business logic stays in the parent model. The picker returns a path; it does not upload, download, or persist files.

## Modes and selection semantics

Replace the current string (`"dir"`/`"file"`) with an exported type and constants. Names are not final; this is an illustrative proposal:

```go
type Mode uint8

const (
    SelectDirectory Mode = iota
    SelectFile
)
```

Required behavior:

- **Directory mode:** `Enter` selects the current directory. Directories in the list are opened using the configured navigation action.
- **File mode:** directories are navigable only; only a visible and enabled file can produce a selection message.
- **Cancel:** `Esc` emits cancellation, never a selection with an empty path.
- **Empty list/error:** file mode must not accidentally select the current directory; the cursor must never be negative and the model must not panic.
- A filter must not turn a directory into a selectable file or prevent navigating through it.

Decide whether directory mode needs a separate action to select a child directory from the list. Do not add one without a clear, documented interaction.

## Constructor options

Keep a simple constructor (`New(Options) Model`) with useful defaults. The initial option set should stay small:

| Option | Purpose | Proposed default |
|---|---|---|
| `Mode` | Choose file or directory selection | Directory, matching the original use |
| `InitialPath` | Start at a specific path | Home if omitted or invalid |
| `Title` | Heading shown above the list | Generic selection title |
| `Width`, `Height` | Available space | Zero until `tea.WindowSizeMsg`, or values supplied by the consumer |
| `ShowHidden` | Include hidden entries | `false` |
| `Layout` | Centered dialog or embedded content | Centered, matching the original experience |
| `Styles` | Override visual styles | Built-in defaults |
| `KeyMap` | Override actions/keys | Arrow keys and current Vim-style bindings |
| `Text` | Replace visible copy | English to preserve current behavior; overridable per instance |
| `Symbols` | Change cursor and entry markers | Built-in symbols with an ASCII fallback |
| `AllowedTypes` | Possible later extension filter | Exclude from v0.1 unless a consumer needs it; empty means all types |

Do not duplicate an option such as `LastUsedPath` if the parent can pass the value as `InitialPath`. Do not persist preferences inside the library.

## Styles

Expose a public styles type for existing and new visual elements:

- title;
- dialog border and list border;
- cursor and selected entry;
- normal directory and file entries;
- secondary text, path, breadcrumbs, and help;
- empty and error states.

Provide `DefaultStyles()` and allow a custom styles value in `Options`. Do not store themes in mutable package globals. Two pickers in one application must be able to render with different styles.

Lip Gloss visible width differs from `len(string)` for ANSI sequences and wide characters. Use terminal-aware measurement and reserve width for symbols, cursor, and padding; do not truncate Unicode by counting bytes.

### Conceptual example

```go
opts := directorypicker.Options{
    Mode:        directorypicker.SelectFile,
    InitialPath: "/tmp",
    ShowHidden:  false,
    Layout:      directorypicker.Embedded,
    Styles:      directorypicker.DefaultStyles(),
}
picker := directorypicker.New(opts)
```

This illustrates the intended configuration. Names and compilability remain subject to implementation and API review.

## Keymap

Configure actions rather than duplicating the logic for every key. Initial actions:

- move up/down;
- open directory;
- go to parent directory;
- select;
- cancel;
- go to first/last entry;
- previous/next page;
- toggle hidden files (when enabled).

Keep familiar defaults: arrow keys and `h/j/k/l`. The picker should safely ignore messages it does not use, allowing the parent to process messages while the picker is inactive. While active, the parent can delegate all `tea.KeyMsg` values to it.

Build visible help from the active key bindings, or document custom bindings clearly. Do not show one key while binding another.

## Layout and dimensions

Propose two composition modes:

- **Centered/modal:** the component places the popup within the full available terminal area.
- **Embedded:** the component returns content for the parent to lay out; it does not center itself in the whole terminal.

Rules:

- Forward every `tea.WindowSizeMsg` while the picker is visible.
- Respect the available space; do not enforce absolute minimums that exceed small terminals.
- Account for borders, padding, breadcrumbs, path, help, and content rows when calculating height and width.
- Keep the cursor visible with a sliding window.
- Test zero width/height, dimensions smaller than the preferred popup, and reduced terminal sizes.
- Do not impose visual bounds unrelated to the terminal without a documented reason.

## Hidden files, filters, and paths

- `ShowHidden` controls initial visibility; the toggle binding changes the model state.
- An extension filter may restrict selection, but directories must remain navigable.
- Define case handling and compound extensions before implementing filters. A later version could start with case-insensitive matching of configured suffixes.
- Do not implicitly expand `~` in every path without documenting it; decide this alongside direct path entry.
- Directory symlinks are followed and treated as normal navigable directories. An application's security must not depend on the picker's visual filter.

## Text and symbols

Allow visible copy and symbols to be customized with small structs, without adding an internationalization package. At minimum cover:

- title and help text;
- current-path and empty-directory labels;
- file/directory indicators;
- selected cursor.

Provide a selectable ASCII fallback for limited terminals. Do not base Unicode support detection solely on `runtime.GOOS`.

## Errors and validation

- Show directory-read errors in the component and expose the error for tests/consumers if this can be done without needless duplicate state.
- An invalid initial path must have a defined fallback or an observable error; do not hide failures without explanation.
- Where appropriate, validate at selection time that the entry still exists and matches the configured mode. The parent may validate again before uploading, deleting, or opening a file.
- Filters are an ergonomic feature, not a security boundary.

## Exclude from the initial API

- Plugins, renderer factories, callbacks for every key, or a generic theme system.
- Persistence, favorite management, or global configuration.
- Recursive filters and previews.
- Destructive file operations.
- Abstraction over other TUI frameworks.

Reconsider an additional option only when a real consuming application needs it and cannot achieve it simply in the parent.

## Customization acceptance criteria

- Two simultaneous instances can use different styles, titles, initial paths, and keymaps.
- Defaults work without configuring every field.
- Help text matches the active keys.
- Selection respects the mode and filters, including empty directories and errors.
- The layout fits small dimensions and preserves Unicode characters.
- Changing options does not require modifying global variables.
