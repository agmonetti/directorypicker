# directorypicker

> Status: initial planning. The code has not been extracted or published yet.

`directorypicker` is a proposed file and directory picker for Go TUI applications built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). The goal is a keyboard-driven, per-instance customizable component that can be used as a dialog or inside an existing layout.

## What we plan to build

- Select files or directories, depending on the configured mode.
- Navigate directories with the keyboard, including arrow keys and Vim-style bindings.
- Show the current path, breadcrumbs, help, and errors.
- Return results through Bubble Tea messages without coupling the component to application logic.
- Let each consuming application adapt the layout, styles, and key bindings.

The prototype comes from a component used internally by `msTTui`. It already separates selection and cancellation through messages and does not depend on private packages from that application.

## Planned customization

The first version prioritizes:

- Per-instance Lip Gloss styles.
- A configurable keymap, with help text consistent with the active bindings.
- Centered or embedded presentation.
- Hidden-file visibility.
- Configurable text and symbols, including an ASCII mode.
- Typed file and directory selection modes.

Local filtering, direct path entry, and cursor restoration may follow. Multi-selection, previews, file operations, and remote providers are out of the initial scope.

## Planned integration

The API will follow the Bubble Tea component pattern (`Init`, `Update`, `View`). The consumer handles selection/cancellation messages and decides what to do with the path—for example, upload a file or save a preferred directory. Directory reads will use Bubble Tea commands so they do not block the UI.

The API and installation example will be published after the module import path, minimum Go/Bubble Tea versions, and first working release are confirmed.

## Design and integration guides

- [PLAN.md](PLAN.md): scope, findings, phases, open decisions, and acceptance criteria.
- [Installation and publishing](docs/INSTALLATION.md): Go module, `go get`, versions, and release checklist.
- [Customization](docs/CUSTOMIZATION.md): proposed API, styles, keymap, modes, and behavior.
- [Integration](docs/INTEGRATION.md): Bubble Tea lifecycle, messages, resizing, commands, and tests.
- [RUNBOOK.md](RUNBOOK.md): session kickoff prompt, current status, and the next task.

The planned distribution is **as a Go module only**. There will be no npm package: npm does not make a Bubble Tea library importable from Go. An npm binary wrapper or a JavaScript port would be a separate product and is out of scope.

## Status and next steps

1. Confirm the module name/import path and compare the scope with `bubbles/filepicker`.
2. Extract the component and fix its selection, empty-list, and small-terminal edge cases.
3. Add tests and a Bubble Tea example covering messages and resizing.
4. Implement essential customization and prepare a tagged release.

## License

MIT. The code will be extracted from another repository owned by the same author, so the source repository's license does not prevent publishing it under MIT. See [LICENSE](LICENSE).
