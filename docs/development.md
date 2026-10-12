# Development

Use Go 1.26+. The project keeps mutable command state inside a newly constructed
root, making repeated invocations and tests independent.

`cmd/dlut-pan` owns process exit and interrupt cancellation. `internal/cli`
validates input, resolves ephemeral credentials and renders output. `internal/api`
implements query encoding, request bodies, HTTP errors and streaming downloads.
`internal/config` contains known endpoint/space constants and environment lookup.

```bash
gofmt -w cmd internal scripts
go vet ./...
go test ./...
go test -race ./...
go mod verify
go run ./scripts/docs
go build -o bin/dlut-pan ./cmd/dlut-pan
go run ./scripts/package -version dev
```

Tests use local HTTP fixtures and temporary files. Do not require a real school
account in CI. Add meaningful tests for endpoint/query/body encoding, input
validation, redaction or filesystem semantics when extending those boundaries.

Generated help lives in `docs/commands`; edit the command tree and regenerate it.
Do not edit generated pages directly. CI regenerates pages and fails on drift.
Keep README claims aligned with implemented commands and separate live account
validation from fixture validation.

The packaging tool emits ZIP archives for linux/darwin/windows on amd64/arm64,
including README, license and authentication notes. It also writes checksums.
To publish a release after reviewing the change, push a semantic-version tag;
the release workflow runs shared checks and publishes their build artifacts.
