# Installation and publishing

> Status: planned installation guide. The package has not been implemented or published yet.

## Goal

`directorypicker` will be a **Go library**, not an executable application. Consumers will add it to an existing Bubble Tea program and manage the model lifecycle.

## Confirmed decisions

- Language: Go.
- UI framework: Bubble Tea.
- Presentation: Lip Gloss.
- License: MIT; see `../LICENSE`.
- Minimum versions: Go 1.21+ / Bubble Tea v0.24+.
- Module path: `github.com/agmonetti/directorypicker`.
- Planned publication: a versioned Go module on GitHub, using semantic version tags.
- Distribution: Go module only. There will be no npm package because it would not be importable from a Bubble Tea application written in Go.

## Planned installation for consumers

Once a tag exists, installation will look like this. No npm package will be published for this Go library; an npm binary wrapper would be a separate product and is not planned.

```sh
go get github.com/agmonetti/directorypicker@v0.1.0
```

For a Go library, `go get` adds or updates the module in the consumer project; it does not install a global binary.

Planned import:

```go
import "github.com/agmonetti/directorypicker"
```

The `v0.1.0` version is illustrative. Use the stable version documented on GitHub; avoid `@latest` in reproducible instructions.

## Planned dependencies

- Bubble Tea: model, messages, and commands.
- Lip Gloss: styling and terminal-aware rendering/measurement.
- Go standard library: filesystem access.

Declare and test minimum versions in `go.mod` during extraction. Avoid additional dependencies until there is a concrete need.

## Using it in a consumer project

Integration will follow the standard Bubble Tea lifecycle:

1. Create the model with `directorypicker.New(options)`.
2. Return its `Init()` command from the application's `Init()`.
3. Forward relevant messages to the picker while it is active and run/return its `tea.Cmd`.
4. Display `picker.View()` or compose it into the layout, depending on the selected mode.
5. Handle public selection and cancellation messages in the parent model.
6. Forward `tea.WindowSizeMsg` for as long as the picker is visible.

The exact example will be added once the API is settled. Until then, see [`INTEGRATION.md`](INTEGRATION.md) for the proposed contract and flow.

## Planned local checks

From the extracted module root:

```sh
gofmt -w .
go test ./...
go vet ./...
```

At least one example must compile. Add CI that runs `go test ./...` and `go vet ./...` using the declared minimum Go version. Do not publish tags until these checks pass.

## Versioning

- `v0.x`: the API may still change; document breaking changes in the changelog/release notes.
- `v1.0.0`: reserve for a tested and stabilized API.
- Every tag must point to code and installation documentation that work with that tag.

`v0.1.0` is a proposed first release, not a date commitment.

## Release checklist

- [ ] Add a minimal `go.mod` with Go 1.21+ and verify with `go list -m`.
- [ ] Run tests, vet, and the example against Bubble Tea v0.24+.
- [ ] Verify the MIT `LICENSE` is included in the tag.
- [ ] Ensure all references point to `github.com/agmonetti/directorypicker`.
- [ ] Create a semantic version tag and test `go get github.com/agmonetti/directorypicker@TAG` from a temporary module.
