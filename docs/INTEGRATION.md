# Integrating with a Bubble Tea application

> Status: proposed integration contract. Snippets are Go-like pseudocode; final names depend on the implemented API.

## Responsibilities

### The picker

- Maintains the current path, entries, cursor, mode, styles, keymap, and loading/error state.
- Reads directories through Bubble Tea commands without blocking the application's `Update` loop.
- Renders the dialog or embedded content.
- Emits selection or cancellation through public messages.

### The consuming application

- Decides when to open/close the picker and which options to pass.
- Chooses the initial path and, if desired, remembers the last path outside the component.
- Forwards messages and window size while the picker is active.
- Decides what to do with the selected path: upload, download, save as a preference, etc.
- Revalidates the path before sensitive operations. Picker visibility/filtering is not a security boundary.

## Lifecycle

### Opening

When opening the picker, construct the model and return its initialization command:

```go
m.picker = directorypicker.New(directorypicker.Options{
    Title:       "Select file",
    InitialPath: m.lastUsedPath,
    Mode:        directorypicker.SelectFile,
    Layout:      directorypicker.Centered,
})
m.pickerOpen = true
return m, m.picker.Init()
```

Directory reads will be asynchronous. When a model is created dynamically during a Bubble Tea session, the parent must return/run `Init()`; do not assume Bubble Tea will call it automatically just because the model was assigned.

### While open

Delegate relevant messages to the picker and return its command. The parent must still update its own layout for `tea.WindowSizeMsg` and forward that size to the picker.

Pseudocode:

```go
func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Update the parent layout first when msg is WindowSizeMsg.
    // The visible picker must also receive it.

    if m.pickerOpen {
        switch msg := msg.(type) {
        case directorypicker.SelectedMsg:
            m.pickerOpen = false
            m.lastUsedPath = msg.Path // optional; persistence belongs to the parent
            return m, m.handleSelection(msg.Path)
        case directorypicker.CancelledMsg:
            m.pickerOpen = false
            return m, nil
        }

        var cmd tea.Cmd
        m.picker, cmd = m.picker.Update(msg)
        return m, cmd
    }

    // Normal application update.
    return m.updateNormally(msg)
}
```

Real code must preserve parent commands as well as picker commands when both need to run. Use `tea.Batch` to combine independent commands; do not accidentally discard an existing command when delegating.

### Rendering

In modal mode, show the picker's `View()` over the normal view while active. In embedded mode, compose the picker's returned content into the parent layout.

```go
func (m appModel) View() string {
    if m.pickerOpen && m.pickerUsesCenteredLayout {
        return m.picker.View()
    }
    return m.renderApplicationWithPickerIfEmbedded()
}
```

The snippet is illustrative; an application may keep banners or borders outside the modal area. If composing content itself, reserve space and avoid drawing the same dialog twice.

## Keys and focus

- While the picker is open, navigation keys go to it and must not trigger commands in the background view.
- The application may allow specific global keys (for example, resize); document the exception.
- The picker ignores messages it does not use.
- On selection or cancellation, close the picker and return focus to the screen that opened it.
- Do not add mouse support in the first release unless its behavior is defined and tested.

## Message contract

Proposed to match the current pattern:

```go
type SelectedMsg struct {
    Path string
}

type CancelledMsg struct{}
```

Rules:

- `SelectedMsg.Path` is non-empty and matches the configured mode.
- Directory selection returns the chosen directory according to the semantics defined in `CUSTOMIZATION.md`.
- Cancellation is distinct from selecting an empty path.
- Do not put a business-logic callback in `Options`; the parent handles the message.

The consumer must check that the path is still available before using it. The filesystem can change between displaying a listing and using an entry.

## Asynchronous reads and stale results

Reads will run through `tea.Cmd`. Each navigation starts a new read and the view may show a loading state. If commands finish out of order, an older result must not replace the latest directory.

The implementation should attach a generation/request ID to each result and discard stale messages. Keep this internal to the package; do not expose it as API.

Minimum states to test:

- successful read with entries;
- empty directory;
- permission error or invalid path;
- rapid navigation where results arrive out of order;
- cancellation while a read is pending;
- closing and reopening the picker.

## Terminal size

The parent model receives `tea.WindowSizeMsg`. Forward it to the picker for as long as the picker is active, including after it is opened. Passing width/height only to the constructor is not enough: later resizes must update the view.

Embedded mode uses the space assigned by the host layout. Centered mode uses the full terminal size. Test narrow/short terminals and resizing during loading and navigation.

## Multiple pickers / concurrency

The first release does not need to manage an arbitrary collection of concurrent pickers. An application can keep one active instance and create it when opening the modal. If simultaneous instances are supported, verify that styles, keymaps, and asynchronous results do not share mutable global state or cross between models.

Do not store the last path in a package-level global; it belongs to the consumer's model.

## Errors

Directory-read errors must be represented in the UI without panics and must not result in a false selection path. The parent can close the picker or let the user navigate away from the failed directory.

For later operations (upload/download), the parent displays business errors. The picker must not know about HTTP clients, credentials, workspaces, or Teams structures.

## Recommended integration tests

- A selection message closes the modal and passes the path to the parent handler.
- Cancellation closes the modal without running the file operation.
- `WindowSizeMsg` reaches the picker after it opens.
- Keys are not sent to the background view while the picker has focus.
- The initialization command is returned when opening the picker.
- The read command and its result messages are handled by the parent's `Update`.
- The parent preserves its own commands when also delegating to the picker.

Component unit tests should cover navigation and selection. Consumer integration tests should verify only the model contract; do not duplicate the picker's internal logic.
