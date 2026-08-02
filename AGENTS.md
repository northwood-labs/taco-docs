# AGENTS.md

Guidance for AI agents working in this repository. Read this before making changes.

## Project overview

**taco-docs** is a hard-fork of [`terraform-docs`](https://terraform-docs.io) maintained by Northwood Labs. It is a Go CLI that generates documentation from Terraform/OpenTofu modules in multiple output formats.

**Current status:** Major refactoring before cutting a first release under the new name.

## Architecture

```text
main.go               → Entry point; delegates to cmd.Execute()
cmd/                  → Cobra CLI commands (one subpackage per output format)
  root.go             → Root command, persistent flags, subcommand registration
  asciidoc/           → AsciiDoc formatter subcommand
  json/               → JSON formatter subcommand
  markdown/           → Markdown formatter subcommand
  pretty/             → Pretty-print (colorized) subcommand
  tfvars/             → tfvars (HCL + JSON) subcommand
  toml/               → TOML formatter subcommand
  xml/                → XML formatter subcommand
  yaml/               → YAML formatter subcommand
  completion/         → Shell completion subcommand
  version/            → Version subcommand
format/               → Output formatter implementations (one per format)
template/             → Go template engine, sanitizers, anchor generation
terraform/            → Terraform/OpenTofu module parsing (inputs, outputs,
                        providers, resources, module calls)
print/                → Print configuration, output settings, utilities
internal/
  cli/                → CLI runtime (PreRunE/RunE wiring, annotations)
  reader/             → File/directory reading utilities
  testutil/           → Shared test helpers
  types/              → Shared type definitions
  version/            → Version string construction
plugin/               → hashicorp/go-plugin based plugin system
                        (client, server, plugin interface)
examples/             → Example Terraform module used for testing/demos
docs/                 → User-facing documentation (developer guide, how-to,
                        reference, user guide)
```

## Technology stack

| Layer         | Technology                                              |
|---------------|---------------------------------------------------------|
| Language      | Go 1.26                                                 |
| CLI framework | `spf13/cobra` + `spf13/viper`                           |
| HCL parsing   | `hashicorp/hcl/v2`, `terraform-config-inspect`          |
| Templating    | Go `text/template` + `Masterminds/sprig`                |
| Plugin system | `hashicorp/go-plugin` (gRPC)                            |
| Config format | YAML (`.terraform-docs.yml`) + TOML + Viper             |
| Build tool    | Taskfile (task) with Makefile shim                      |
| Linting       | `golangci-lint` (~60 linters), `hk` (pre-commit runner) |
| Markdown lint | `rumdl` (Rust-based Markdown linter)                    |
| Link checking | `lychee`                                                |
| Security      | `trivy`, `osv-scanner`, `trufflehog`, `govulncheck`     |
| CHANGELOG     | `git-cliff` (conventional commits)                      |
| CI/CD         | GitHub Actions                                          |
| Config sync   | `config-manager` (syncs from upstream `.github` repo)   |

## Build and development commands

The project uses [Taskfile](https://taskfile.dev) as its task runner. The `Makefile` is a thin shim that forwards to `task`.

```bash
# List all available tasks
task --list

# Build and install locally
task build:go

# Run linters on changed files (fast, default)
task lint

# Run linters on ALL files (comprehensive)
task lint:deep

# Run a specific linter on changed files
task lint:golangci-lint

# Tidy Go modules
task tidy:go

# Update Go dependencies
task deps:go

# Install development tools (via Homebrew)
task install:tools

# Install git hooks
task install:hooks

# Clean various caches
task clean
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./format/...
go test ./terraform/...

# Run tests with verbose output
go test -v ./...
```

Test data lives in `testdata/` subdirectories within each package that has tests.

## Linting and code quality

### Zero diagnostics policy

This project enforces a **Zero Diagnostics Policy**. Code is not complete until all linters report zero issues. The full policy and fix patterns are documented in `.kiro/steering/go-code-conventions.md` (loaded automatically when editing `.go` files).

### Key rules

* **Never edit `.golangci.yml`** to resolve diagnostics.
* **Never use `// nolint:xxx`** directives. Use `// lint:*` comments only (see suppression table in steering doc).
* Suppression is a last resort — fix the root cause.
* After editing any `.go` file: run `golangci-lint run --fix ./...`, then `go vet ./...`, then `go test` for the affected package.
* Line length limit is **120 characters** for code, **80 characters** for comments.

### Pre-commit hooks

Managed by [hk](https://github.com/jdx/hk) via `hk.pkl` (Pkl language config). Hooks run on pre-commit (fast) and pre-push (slow profile including golangci-lint, govulncheck, trivy).

```bash
# Install hooks
hk install

# Run all hooks manually (fix mode)
hk fix --no-fail-fast

# Run a specific step
hk fix --step golangci-lint --profile slow
```

## Kiro configuration

### Steering documents

Located in `.kiro/steering/`:

| File                     | Inclusion | Purpose                                                                             |
|--------------------------|-----------|-------------------------------------------------------------------------------------|
| `go-code-conventions.md` | `**/*.go` | Zero Diagnostics Policy, lint fix patterns, code conventions, suppression reference |
| `markdown-style.md`      | `**/*.md` | Markdown formatting rules enforced by rumdl                                         |

These are automatically injected into agent context when files matching the glob are read. They contain extensive detail — treat them as authoritative.

### Hooks

No custom Kiro hooks are currently defined (`.kiro/hooks/` does not exist yet). The project uses `hk` for git hooks instead.

### Config-manager

Files marked with `@config-manager:start`/`@config-manager:end` comment blocks are managed by the [config-manager](https://github.com/northwood-labs/config-manager) tool. **Do not edit content inside these markers** unless you understand the sync implications — changes may be overwritten on next sync.

Configuration lives in `.config-manager.d/`:

* `config.toml` — Shared upstream repo reference and settings
* `golang.toml` — Go-specific managed config
* `standard.toml` — Standard project config (linters, CI, etc.)

## Conventions for agents

### Making go changes

1. Read the file before editing.
2. Follow declaration order: `const`, `var`, `type`, `func`.
3. Only one `var()` block per file.
4. Max 4 parameters per function (excluding `context.Context`) — use input structs.
5. Max 3 return values (excluding `error`) — use result structs.
6. Never combine function call + nil-check in a single `if` statement.
7. Use `filepath.Join` for paths, never string concatenation with `/`.
8. Use `strings.Builder` for output collection, not repeated `fmt.Print*` calls.
9. Sentinel errors live in `cmd/errors.go`; use `Err` prefix + PascalCase.
10. All `slog` calls must use `*Context` variants with constant keys in `snake_case`.

### Making Markdown changes

1. Use ATX headings (`#`), asterisk lists (`*`), underscore emphasis (`_`).
2. H1 in Title Case; all other headings in sentence case.
3. Always specify language on fenced code blocks.
4. Proper names have exact casing (see steering doc for the full list).
5. 2-space indentation, LF line endings, single trailing newline.

### Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org) syntax. The `gommit` tool validates commit messages on pre-commit. Common prefixes:

* `chore:` — Maintenance tasks
* `docs:` — Documentation changes
* `feat:` — New features
* `fix:` — Bug fixes
* `lint:` — Linting fixes
* `perf:` — Performance improvements
* `refactor:` — Code restructuring
* `test:` — Test additions/changes

### Branch strategy

* Never push directly to `main`.
* Create feature branches and open pull requests.
* CI generates CHANGELOG automatically on push to `main`.

## File layout conventions

* One Cobra command per subpackage in `cmd/`.
* One formatter per file in `format/` (e.g., `json.go`, `yaml.go`).
* Test files are colocated with the code they test (`*_test.go`).
* Test fixtures go in `testdata/` directories.
* Shell scripts go in `scripts/`.
* Taskfile includes go in `tools/Taskfile.*.yml`.
* Homebrew dependency lists go in `tools/Brewfile.*`.

## Tools reference

For detailed tool documentation, see `docs/TOOLS.md`.
