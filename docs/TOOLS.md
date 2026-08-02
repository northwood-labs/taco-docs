# Tools Reference

Detailed reference for the tools used in this project. This supplements `AGENTS.md` with specifics an agent may need when debugging or running tools.

## Build and task runner

### Taskfile (task)

* **Config:** `Taskfile.dist.yml` (root), with includes from `tools/Taskfile.*.yml`
* **Override:** Create `Taskfile.yml` at root to override tasks without touching the managed config.
* **Includes:** `Taskfile.x.yml` (optional, project-specific non-synced tasks)
* **Shell:** Uses `/opt/homebrew/opt/bash/bin/bash` (GNU bash via Homebrew)
* **Utilities:** GNU find, grep, sed, rm, xargs (Homebrew-installed at `/opt/homebrew/opt/*/bin/g*`)

### Makefile

Thin shim only — forwards all targets to `task`. Exists for muscle memory compatibility.

## Go toolchain

| Tool             | Purpose                             | Invocation                                                         |
|------------------|-------------------------------------|--------------------------------------------------------------------|
| `go build`       | Compile                             | `task build:go`                                                    |
| `go test`        | Run tests                           | `go test ./...`                                                    |
| `go vet`         | Static analysis                     | `go vet ./...`                                                     |
| `go mod tidy`    | Clean module deps                   | `task tidy:go`                                                     |
| `golangci-lint`  | Meta-linter (~60 linters)           | `golangci-lint run --fix ./...`                                    |
| `gofumpt`        | Strict gofmt                        | Via `hk` or `golangci-lint`                                        |
| `goimports`      | Import organization                 | Via `hk` (local prefix: `github.com/northwood-labs,go.nwlabs.dev`) |
| `golines`        | Line length enforcement (120 chars) | Via `hk`                                                           |
| `fieldalignment` | Struct field alignment optimization | Via `hk`                                                           |
| `modernize`      | Apply Go idiom modernizations       | Via `hk`                                                           |
| `smrcptr`        | Same receiver pointer consistency   | Via `hk`                                                           |
| `govulncheck`    | Known vulnerability scanning        | Via `hk` (slow profile)                                            |

### Golangci-lint

* **Config:** `.golangci.yml` (v2 format, `default: none`)
* **Never edit this file** to resolve diagnostics.
* Uses `source:` pattern matching for lint suppression comments (not `// nolint`).
* Depends on: fieldalignment, go-fix, go-fmt, go-modernize, go-same-receiver-pointer, go-vet, golines (must pass first).

## Linting orchestration (hk)

* **Config:** `hk.pkl` (Pkl language)
* **Version:** hk v1.48.0
* **Profiles:** default (fast), `slow` (includes `golangci-lint`, `govulncheck`, `trivy`, `golines`)

### Linter groups defined in hk.pkl

| Group        | Notable steps                                                                                 |
|--------------|-----------------------------------------------------------------------------------------------|
| `standard`   | `editorconfig-checker`, `rumdl`, `lychee`, `pinact`, `pkl`, `yamlfmt`, `zizmor`               |
| `containers` | `hadolint`                                                                                    |
| `golang`     | `go-fmt`, `goimports`, `golines`, `golangci-lint`, `go-vet`, `go-modernize`, `fieldalignment` |
| `security`   | `govulncheck`, `osv-scanner`, `trivy` (licenses + vulns), `trufflehog`                        |
| `shell`      | `shellcheck`, `shfmt`                                                                         |
| `github`     | `ghalint` (action + workflow linting)                                                         |
| `tofu`       | OpenTofu-format, `tflint`                                                                     |
| `python`     | `ruff`, `ruff-format`, `zuban` (type checking)                                                |
| `frontend`   | `oxfmt`, `oxlint`, `sort-package-json`, `tsc`                                                 |

## Markdown tooling

| Tool                   | Purpose                  | Config                       |
|------------------------|--------------------------|------------------------------|
| `rumdl`                | Markdown linting         | `.rumdl.toml`                |
| `editorconfig-checker` | EditorConfig enforcement | `.editorconfig-checker.json` |
| `lychee`               | Link checking            | `lychee.toml`                |

### `rumdl`

Rust-based Markdown linter. Key settings:

* GFM flavor
* Sentence case for H2-H6 headings
* Asterisk unordered lists, ordered numbering for ordered lists
* Fenced code blocks with language identifiers required
* Leading and trailing pipes on tables
* Proper names enforcement (see `.rumdl.toml` `[MD044].names`)
* Cache enabled at `.rumdl_cache/`

## Security scanning

| Tool          | Config                                                       | Purpose                          |
|---------------|--------------------------------------------------------------|----------------------------------|
| `trivy`       | `trivy-vuln.yaml`, `trivy-license.yaml`, `.trivyignore.yaml` | Vulnerability + license scanning |
| `osv-scanner` | `osv-scanner.toml`                                           | OSV database vulnerability check |
| `trufflehog`  | (none)                                                       | Secret detection in git history  |
| `govulncheck` | (none)                                                       | Go-specific vuln scanning        |

## Configuration management

### Config-manager

Syncs shared configuration from an upstream repository. Managed sections are delimited by:

```text
@config-manager:start <section-name>
...managed content...
@config-manager:end <section-name>
```

* **Config:** `.config-manager.d/config.toml`
* **Lock file:** `.repo-config.lock`
* **Source repo:** `file://../.github` (local sibling)

Files with managed sections include: `Taskfile.dist.yml`, `hk.pkl`, `.golangci.yml`, `.rumdl.toml`, `.editorconfig`, GitHub Actions workflows, and more.

## Editor and formatting

### EditorConfig

* **Config:** `.editorconfig`
* UTF-8, LF line endings, trim trailing whitespace
* Go: tab indent
* Markdown: 2-space indent
* Shell: 4-space indent, 120 char line length
* YAML/JSON/TOML/Pkl: 2-space indent

### YAML formatting

* **Tool:** `yamlfmt`
* **Config:** `.yamlfmt.yml`

### TOML formatting

* **Tool:** `taplo`
* **Config:** `.taplo.toml`

## CHANGELOG and releases

### `git-cliff`

* **Config:** `cliff.toml`
* Generates `CHANGELOG.md` from conventional commits
* Runs automatically via GitHub Actions on push to `main`
* Groups: Features, Bug Fixes, Performance, Documentation, Refactor, etc.

### Commit validation

* **Tool:** `gommit`
* **Config:** `.gommit.toml`
* Validates conventional commit message format on pre-commit

## CI/CD (GitHub actions)

### Workflows

| Workflow                    | Trigger    | Purpose                                 |
|-----------------------------|------------|-----------------------------------------|
| `update-on-push.yml`        | push main  | Generate CHANGELOG, update CONTRIBUTORS |
| `dependabot-auto-merge.yml` | Dependabot | Auto-merge Dependabot PRs               |

### Security hardening

* Uses `step-security/harden-runner` with egress policy blocking
* Explicit allowed-endpoints list for network access
* `permissions: read-all` default with scoped overrides
* Concurrency groups to prevent parallel workflow runs

## Homebrew dependencies

Defined in `tools/Brewfile.*`:

* `tools/Brewfile.standard` — Cross-language tools (linters, formatters, security)
* `tools/Brewfile.golang` — Go-specific tools

Install all with: `task install:tools`
