# Quick Flow Summary

A concise overview of how taco-docs works, from process start to documentation output.

## Entry point

```text
main.go → cmd.Execute() → NewCommand().Execute()
```

`main.go` does exactly one thing: call `cmd.Execute()`. On error, it exits with code 1. All logic lives in the `cmd` and `internal/cli` packages.

## Primary flow

```text
┌─────────────────────────────────────────────────────────────────────┐
│ 1. BUILD CLI TREE                                                   │
│    NewCommand() creates root cobra.Command with:                    │
│    * Shared Config (print.DefaultConfig())                          │
│    * Shared Runtime (cli.NewRuntime(config))                        │
│    * Persistent flags bound to Config fields                        │
│    * All formatter subcommands (markdown, json, yaml, etc.)         │
└─────────────────────────────────────────┬───────────────────────────┘
                                          │
┌─────────────────────────────────────────▼───────────────────────────┐
│ 2. PreRunE — RESOLVE CONFIGURATION                                  │
│    * Read formatter name from command annotations                   │
│    * Search for .terraform-docs.yml (module → CWD → HOME)           │
│    * Unmarshal YAML via viper into Config struct                    │
│    * Overlay CLI flags (flags win over config file)                 │
│    * Enforce version constraint from config                         │
└─────────────────────────────────────────┬───────────────────────────┘
                                          │
┌─────────────────────────────────────────▼───────────────────────────┐
│ 3. RunE — GENERATE DOCUMENTATION                                    │
│    * Discover modules (main + recursive submodules)                 │
│    * For each module:                                               │
│      a. terraform.LoadWithOptions(config)                           │
│         Parse .tf → enrich (comments, lock, providers) → sort       │
│      b. format.New(config)                                          │
│         Look up formatter in registry → return concrete type        │
│      c. formatter.Generate(module)                                  │
│         Render each section via Go templates                        │
│      d. formatter.Render(contentTemplate)                           │
│         Apply user's content template (or use default order)        │
│      e. writeContent(config, output)                                │
│         → stdout (no --output-file)                                 │
│         → file inject (between <!-- markers -->)                    │
│         → file replace (overwrite entire file)                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Module roles

| Package             | Responsibility                                            |
|---------------------|-----------------------------------------------------------|
| `cmd/`              | CLI structure, flag definitions, subcommand registration  |
| `internal/cli/`     | Runtime state, config resolution, file writing            |
| `print/`            | Config data model, validation, section visibility logic   |
| `terraform/`        | Parse Terraform modules, extract all documentable items   |
| `format/`           | Formatter registry, section rendering, output composition |
| `template/`         | Go template engine, sanitizers, built-in functions        |
| `plugin/`           | SDK for external formatter plugins (client + server)      |
| `internal/plugin/`  | Plugin binary discovery and RPC connection setup          |
| `internal/reader/`  | Line-by-line file extraction (comments, headers)          |
| `internal/types/`   | Rich value types for Terraform primitives                 |
| `internal/version/` | Version string construction                               |

## Design decisions

### Single shared config pointer

One `*print.Config` instance is created at root level and shared across all subcommands via pointer. This eliminates flag-to-config plumbing for each format and ensures persistent flags propagate automatically.

### Annotations over string parsing

The formatter name is carried in cobra command annotations (`"command": "markdown table"`) rather than parsed from the `Use` field. This decouples command naming from formatter lookup and supports aliases without fragile string manipulation.

### Registry pattern for formatters

Each formatter self-registers via `init()` into a package-level map. Adding a new format requires only creating the file with its `init()` function — no central import list or switch statement to modify.

### Plugin fallback

When `format.New()` can't find a built-in formatter, the CLI transparently falls through to plugin discovery. Users specify a formatter name and it works identically whether built-in or plugin-provided.

### Template-based rendering with content templates

Formatters render each section independently (header, inputs, outputs, etc.) via Go templates. A user-supplied content template can then recompose those sections in any order, supporting customized documentation layouts without modifying formatter code.

### Inject mode for README co-existence

The writer's inject mode splices generated content between comment markers (`<!-- BEGIN_TF_DOCS -->` / `<!-- END_TF_DOCS -->`), preserving hand-written content above and below. This enables a single README with both manual prose and auto-generated docs.

## Risks and unknowns

* **`terraform-config-inspect` is unmaintained.** The upstream library doesn't support newer OpenTofu constructs (e.g., `for_each` on providers). Workarounds exist in `fixOpenTofuProviders()` but may need expanding as OpenTofu evolves.

* **Provider function discovery is best-effort.** The HCL AST walk in `loadProviderFunctions` handles common expression types but may miss deeply nested or dynamically generated function calls.

* **Remote header/footer fetching has no caching.** HTTP requests for headers/footers happen on every invocation with a 5-second timeout. Repeated CI runs against unreliable URLs could cause intermittent failures.

* **Plugin system uses net/rpc (not gRPC).** The current hashicorp/go-plugin integration uses the older RPC protocol. This works but limits future extensibility compared to the gRPC transport that go-plugin also supports.

* **Recursive mode requires --output-file.** When documenting submodules recursively, stdout output is disallowed (multiple modules would concatenate without separators). This is validated at runtime but could surprise users.

* **The project is mid-refactor.** The binary is still named `terraform-docs` in the `Use` field and Dockerfile. Renaming to `taco-docs` is pending before the first release.
