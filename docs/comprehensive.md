# Deep Architecture Audit

A comprehensive walkthrough of taco-docs internals: every entry point, initialization path, module boundary, decision point, and known risk.

## Entry points

| Entry point      | File                             | Purpose                                            |
|------------------|----------------------------------|----------------------------------------------------|
| Process start    | `main.go`                        | Calls `cmd.Execute()`, translates error to exit(1) |
| CLI execution    | `cmd/root.go`                    | Builds command tree, dispatches to cobra           |
| Config-only read | `print/config.go` `ReadConfig()` | Standalone config loading for tests/plugins        |
| Plugin server    | `plugin/server.go` `Serve()`     | Entry point for plugin binaries                    |

## CLI startup and initialization flow

### Phase 1: Command tree construction (`NewCommand`)

```text
cmd.Execute()
  └─ NewCommand()
       ├─ print.DefaultConfig()           → shared *Config with safe defaults
       ├─ cli.NewRuntime(config)          → shared *Runtime holding state
       ├─ cobra.Command{                  → root command
       │    Annotations: {"command":"root", "kind":"formatter"}
       │    PreRunE: runtime.PreRunEFunc
       │    RunE:    runtime.RunEFunc
       │  }
       ├─ PersistentFlags()               → bound directly to Config fields
       │    --config, --recursive, --recursive-path, --recursive-include-main
       │    --recursive-exclude, --show, --hide, --output-file, --output-mode
       │    --output-template, --output-check, --sort, --sort-by
       │    --header-from, --footer-from, --lockfile, --output-values
       │    --output-values-from, --read-comments
       └─ AddCommand()
            ├─ asciidoc.NewCommand(runtime, config)
            ├─ json.NewCommand(runtime, config)
            ├─ markdown.NewCommand(runtime, config)
            │    ├─ markdown/document.NewCommand(runtime, config)
            │    └─ markdown/table.NewCommand(runtime, config)
            ├─ pretty.NewCommand(runtime, config)
            ├─ tfvars.NewCommand(runtime, config)
            ├─ toml.NewCommand(runtime, config)
            ├─ xml.NewCommand(runtime, config)
            ├─ yaml.NewCommand(runtime, config)
            ├─ completion.NewCommand()
            └─ version.NewCommand()
```

**Key insight:** Every formatter subcommand receives the same `runtime` and `config` pointers. Cobra's persistent flag mechanism ensures flag values set at the root propagate to children automatically.

### Phase 2: Configuration resolution (`PreRunEFunc`)

```text
PreRunEFunc(cmd, args)
  │
  ├─ Read formatter name from cmd.Annotations["command"]
  │   (avoids parsing Use string; supports aliases)
  │
  ├─ If formatter == "root" && no args → show help, return nil
  │
  ├─ Store: isFlagChanged, rootDir, cmd on Runtime
  │
  ├─ Validate --config is non-empty
  │
  ├─ readConfig(viper, file, "")
  │   ├─ If --config flag was explicitly changed:
  │   │     Resolve to absolute path, set as viper config file
  │   ├─ Otherwise:
  │   │     Set config name ".terraform-docs", type "yml"
  │   │     Add search paths:
  │   │       1. rootDir/
  │   │       2. rootDir/.config/
  │   │       3. ./
  │   │       4. .config/
  │   │       5. $HOME/.tfdocs.d/
  │   └─ v.ReadInConfig()
  │       * PathError + explicit file → ErrConfigFileNotFound
  │       * ConfigFileNotFoundError → OK (no config is valid)
  │       * Other error → propagate
  │
  ├─ unmarshalConfig(viper, config)
  │   ├─ bindFlags(viper)
  │   │     For each changed flag:
  │   │       --show/--hide → clear both lists first, then set
  │   │       --sort-by-required/--sort-by-type → legacy flag mapping
  │   │       default → map flag name to viper key
  │   ├─ v.Unmarshal(config) → decode into Config struct
  │   ├─ If formatter != "root": force config.Formatter = formatter
  │   │   (explicit subcommand wins over config file)
  │   └─ config.Parse() → resolve show/hide into boolean section flags
  │
  └─ checkConstraint(config.Version, version.Core())
      Enforces "version:" field from config file against running binary
```

**Decision points in PreRunE:**

* Config file location: explicit flag > module root > CWD > HOME
* Flag override semantics: CLI flags always win over config file values
* Show/hide: first changed flag clears both lists (replace, not merge)
* Formatter name: explicit subcommand > config file > error

### Phase 3: Documentation generation (`RunEFunc`)

```text
RunEFunc(cmd, args)
  │
  ├─ Build module list:
  │   ├─ If !recursive OR includeMain: append {config, rootDir}
  │   └─ If recursive && path != "":
  │       findSubmodules()
  │         Walk recursive.Path, skip hidden dirs and excluded names
  │         For each dir with .tf/.tf.json files:
  │           loadSubModule(path)
  │             loadModuleConfig(path)
  │               Try reading path/.terraform-docs.yml
  │               If found: mergeConfig() (submodule overrides root)
  │               If not: return nil (root config applies)
  │
  ├─ For each module:
  │   ├─ Resolve effective config (submodule config or root config)
  │   ├─ config.Validate()
  │   │     Check formatter non-empty, headerFrom valid, footerFrom valid,
  │   │     recursive path, sections, output template markers, sort type
  │   ├─ If recursive && output.File == "": error (can't mix to stdout)
  │   └─ generateContent(config) → see "Content generation pipeline"
  │
  └─ Return first error encountered (fail-fast)
```

## Command-specific flows

### Formatter subcommands (`markdown`, `json`, `yaml`, etc.)

Each formatter subcommand follows an identical pattern:

1. Create `*cobra.Command` with format-specific flags
2. Set `Annotations: cli.Annotations("formatter-name")`
3. Wire `PreRunE: runtime.PreRunEFunc` and `RunE: runtime.RunEFunc`
4. Return command (added to root via `AddCommand`)

Format-specific flags (e.g., `--anchor`, `--indent` for Markdown) are defined as persistent flags on the subcommand, making them available to sub-sub-commands (like `markdown document` and `markdown table`).

### Content generation pipeline (`generateContent`)

```text
generateContent(config)
  │
  ├─ terraform.LoadWithOptions(config)
  │   ├─ loadModule(path)
  │   │     tfconfig.LoadModule(path) → parse all .tf files
  │   │     Filter benign diagnostics (OpenTofu for_each provider errors)
  │   │     fixOpenTofuProviders() → HCL AST re-read for missing providers
  │   │
  │   ├─ loadModuleItems(tfmodule, config)
  │   │     loadHeader(config)    → from .tf comment block or .md/.adoc file
  │   │     loadFooter(config)    → same, if configured
  │   │     loadInputs()          → variables + comment fallback + CRLF norm
  │   │     loadModulecalls()     → module blocks + source/version extraction
  │   │     loadOutputs()         → outputs + optional terraform output values
  │   │     loadProviders()       → from resource usage + lock file versions
  │   │     loadProviderFunctions() → HCL AST walk for provider::name::func
  │   │     loadRequirements()    → required_core + required_providers
  │   │     loadResources()       → managed + data resources, deduplicated
  │   │
  │   └─ sortItems(module, config)
  │         Sort each collection by name/required/type/position
  │
  ├─ format.New(config)
  │   │ Look up config.Formatter in global initializers map
  │   │ If found → return built-in formatter
  │   │ If not found → fall through to plugin discovery
  │   │
  │   └─ [plugin fallback path]
  │       plugin.Discover()
  │         Scan: TFDOCS_PLUGIN_DIR → ./.plugins/ → ~/.tfdocs.d/plugins/
  │         For each tfdocs-format-* binary:
  │           Spawn subprocess, establish RPC handshake
  │         Return plugin List
  │       plugins.Get(formatterName)
  │       client.Execute(&ExecuteArgs{Module, Config})
  │       writeContent(config, pluginOutput)
  │       return
  │
  ├─ formatter.Generate(module)
  │   │ Template-based formatters (markdown, asciidoc):
  │   │   forEach(render) → render each section template independently
  │   │     "all", "header", "footer", "inputs", "modules",
  │   │     "outputs", "providers", "requirements", "resources"
  │   │   Each section → template.Render(name, module) → string
  │   │
  │   │ Data formatters (json, yaml, toml, xml):
  │   │   copySections(config, module) → filter by show/hide
  │   │   Serialize filtered module with stdlib encoder
  │   │   Store result via withContent()
  │   │
  │   └─ Section strings stored on generator fields
  │
  ├─ formatter.Render(config.Content)
  │   │ If canRender==false (json, yaml, xml, toml):
  │   │   Return content as-is (custom templates incompatible)
  │   │ If content template is empty:
  │   │   Return content as-is (default section ordering)
  │   │ Otherwise:
  │   │   Parse user's content template
  │   │   Expose generator sections as {{ .Header }}, {{ .Inputs }}, etc.
  │   │   Expose {{ include "file" }} and {{ include_optional "file" "fallback" }}
  │   │   Execute template → final output string
  │   │
  │   └─ Return rendered string
  │
  └─ writeContent(config, content)
        If config.Output.File == "":
          stdoutWriter → print to stdout + newline
        Else:
          fileWriter with mode:
            "replace" → overwrite entire file
            "inject"  → splice between begin/end markers
          If --output-check:
            Compare against existing file, error if different
            (for CI staleness detection)
```

## File and module responsibilities

### `main.go`

Single responsibility: translate `cmd.Execute()` error into non-zero exit code. No logic, no imports beyond `os` and `cmd`.

### `cmd/root.go`

Constructs the complete CLI tree. Creates the shared Config and Runtime. Defines all persistent flags. Registers all subcommands.

### `cmd/<format>/` subpackages

Each subpackage exposes `NewCommand(runtime, config)` returning a configured `*cobra.Command`. Format-specific flags (anchor, indent, escape, etc.) are defined here. They share `PreRunE`/`RunE` with the root.

### `internal/cli/run.go`

The orchestration brain. `Runtime` struct bridges `PreRunE` and `RunE` state. Responsibilities:

* Config file discovery and reading (multi-path viper search)
* Config unmarshalling with flag overlay
* Version constraint enforcement
* Module discovery (recursive walk)
* Submodule config merging
* Core pipeline: load → format → write

### `internal/cli/writer.go`

Output destination abstraction. Two implementations:

* `stdoutWriter` — appends newline, writes to `os.Stdout`
* `fileWriter` — supports inject/replace modes, check mode for CI, template wrapping around content

### `print/config.go`

Central configuration model. Defines `Config` struct with all user preferences. Provides `DefaultConfig()` (safe defaults), `Validate()` (fail-fast checks), `Parse()` (resolve show/hide to booleans). Section visibility logic (`show` wins over `hide`, `all` keyword supported).

### `terraform/load.go`

Module parsing orchestrator. Uses `terraform-config-inspect` as the primary parser, then enriches with:

* Comment extraction (fallback descriptions)
* Lock file version merging
* Provider function discovery via HCL AST walk
* Remote content fetching (HTTP headers/footers)
* Output value injection from `terraform output -json`
* Resource deduplication and registry URL construction

### `terraform/module.go`

Domain model. The `Module` struct is the universal data contract between loading and formatting. Includes `Has*()` helpers for template guards.

### `format/type.go`

Formatter registry and factory. Global `initializers` map populated by `init()` calls. `New(config)` factory resolves formatter name to concrete type. Defines the `Type` interface (Generate, Render, Content, section accessors).

### `format/generator.go`

Shared base for all formatters. The `generator` struct stores rendered section strings. `forEach` iterates section names, calls the formatter's render function, and stores results. `Render` applies user content templates with `include`/`include_optional` support.

### `format/<name>.go` (one per format)

Each file: defines concrete struct (embeds `*generator`), registers via `init()`, implements `Generate(module)`. Template-based formats call `forEach`; data formats use stdlib encoders.

### `template/template.go`

Go `text/template` wrapper. Provides `Render(name, module)` and `RenderContent(name, data)`. Registers built-in functions (sanitize, anchor, indent, trim, ternary, tostring) plus all sprig functions. `normalize()` strips leading whitespace from template source for readable Go code.

### `template/sanitizer.go`

Character escaping for Markdown and AsciiDoc. Handles pipe characters in tables, special Markdown chars in names, newlines in descriptions. Format-aware (different rules for table cells vs. prose sections).

### `template/anchor.go`

Generates anchor links for inputs/outputs. Format-specific (Markdown uses `#`, AsciiDoc uses `<<>>`). Respects escape and anchor settings.

### `plugin/` (SDK package)

Public API for plugin authors:

* `Serve(opts)` — call from `main()` to start plugin RPC server
* `Client` — RPC client wrapping Name/Version/Execute calls
* `ExecuteArgs` — the data sent to plugins (Module + Config)
* Handshake config shared between client and server

### `internal/plugin/discovery.go`

Plugin binary scanner. Search priority: `TFDOCS_PLUGIN_DIR` env > `./.plugins/` > `~/.tfdocs.d/plugins/`. Naming convention: `tfdocs-format-<name>`. Each binary spawned as subprocess with [go-plugin] RPC handshake.

### `internal/reader/lines.go`

Line-by-line file reader with configurable condition (which lines to read) and parser (how to transform each line). Used by `loadComments` and `loadSection` in the `terraform` package.

### `internal/types/`

Rich value types that handle Terraform's type system complexity. `ValueOf(any)` converts raw values into typed wrappers (String, Number, Bool, List, Map, Nil, Empty). Custom JSON/XML/YAML marshalling for each type.

### `internal/version/version.go`

Version string construction. `Core()` returns bare semver, `Short()` adds `v` prefix, `Full()` appends commit hash and prerelease info. Used for `--version` flag and constraint checking.

## Decision points and side effects

### Config resolution precedence

```text
1. CLI flags (always win)
2. Config file values (from first file found in search path)
3. DefaultConfig() values (hardcoded safe defaults)
```

### Formatter selection

```text
1. Explicit subcommand name → forced into config.Formatter
2. Root command + config file "formatter:" field → used as-is
3. Root command + no config → error (ErrFormatterEmpty)
```

### Section visibility

```text
If --show specified:  only listed sections rendered (whitelist)
If --hide specified:  all except listed sections rendered (blacklist)
If neither:           all sections rendered
If both:              validation error (mutual exclusion)
```

### Output destination

```text
If --output-file empty:          stdout
If --output-file + mode inject:  splice between markers in existing file
If --output-file + mode replace: overwrite entire file
If --output-check:               compare only, error if stale (CI mode)
```

### Side effects during execution

| Action                     | Side effect                                     |
|----------------------------|-------------------------------------------------|
| `loadOutputValues`         | May exec `terraform output -json` subprocess    |
| `loadSection` (web source) | HTTP GET to remote URL, creates temp file       |
| `fileWriter.Write`         | Creates or modifies files on disk               |
| `fileWriter` check mode    | Reads file, reports mismatch (no writes)        |
| Plugin discovery           | Spawns plugin binaries as subprocesses          |
| `fmt.Printf` in writer     | Prints status messages ("updated successfully") |

## Risks, gaps, and follow-up inspections

### Architectural risks

* **`terraform-config-inspect` dependency.** Unmaintained upstream. The library doesn't handle OpenTofu-specific syntax. The `fixOpenTofuProviders` workaround is ad-hoc and may need expansion for new OpenTofu features (e.g., `moved` blocks, `import` blocks with complex expressions).

* **Global formatter registry via init().** The init-based registration pattern makes it impossible to test formatters in isolation without importing the entire format package. It also means the registry is process-global with no reset mechanism for parallel tests.

* **Plugin process lifecycle.** Plugin subprocesses are spawned during `generateContent` but cleanup (`List.Clean()`) is never called in the main flow. Long-running or abandoned plugin processes could leak.

* **Template panic on include failure.** The `include` function in content templates calls `panic(err)` on file read failure rather than returning an error. This produces an unrecoverable crash with a stack trace instead of a user-friendly error message.

### Data integrity gaps

* **Comment extraction is heuristic.** `loadComments` reads contiguous `#`/`//` lines above a declaration. Multi-line `/* */` comments before a variable are not captured by this mechanism (they're only handled for headers/footers).

* **Provider function discovery misses edge cases.** The AST walker handles common expression types but doesn't cover all possible nesting (e.g., `dynamic` blocks with provider functions in content expressions are not explicitly handled).

* **Output value merging is best-effort.** Outputs not found in `terraform output -json` result are silently set to `"null"` rather than flagging a potential mismatch.

* **Resource version extraction is lossy.** `resourceVersion` takes only the last constraint and strips operators. A constraint like `>= 3.0, < 4.0` would yield `"latest"`, losing useful information.

### Operational gaps

* **No structured logging.** The codebase uses `fmt.Println` for status messages and error reporting. There is no log-level control, no structured output, and no way to suppress informational messages without redirecting `stderr`.

* **HTTP fetching has no retry or caching.** Remote header/footer content is fetched fresh every run with a 5-second timeout and no retry logic. CI environments with network issues will see intermittent failures.

* **Temp file cleanup relies on defer.** Remote content creates temp files that are cleaned via `defer os.Remove()`. If the process is killed (SIGKILL), temp files may accumulate.

* **File permissions inconsistency.** `writeToFile` uses `0o644` with a `lint:allow_possible_insecure` comment. The steering doc specifies `0o0666` for files to support Windows. This should be reconciled.

### Follow-up inspections recommended

* Audit all `init()` registrations to ensure no naming collisions between formats
* Review `copySections` filtering logic to confirm it handles all section types consistently
* Verify plugin handshake config is versioned correctly for backward compatibility
* Check that `fixOpenTofuProviders` handles all resource types (not just managed/data)
* Assess whether `loadProviderFunctions` should report partially-parsed files as warnings
* Confirm the `Dockerfile` build stage produces a working binary (currently commented out)

## Design rationale

### Why a shared Config pointer across all subcommands

Cobra persistent flags write directly to the pointed-to struct fields. By sharing a single `*Config`, flag values propagate from root to any subcommand without manual copying. This eliminates an entire class of "forgot to pass the flag value" bugs and keeps the flag definitions co-located with their effect.

### Why annotations instead of parsing command names

Command names can have aliases (`markdown` / `md`), sub-subcommands (`markdown table`), and the `Use` field contains usage syntax (not a clean identifier). Annotations provide a stable, non-fragile mapping from command to formatter name that survives renaming and aliasing.

### Why init() registration for formatters

The pattern achieves open/closed principle: new formats are added by creating a file, not modifying existing code. The Go linker ensures the init function runs when the package is imported (which happens via the `cmd/<format>` subpackage imports in `root.go`). No central switch statement or factory map to maintain.

### Why a two-phase `PreRunE`/`RunE` split

`PreRunE` handles configuration resolution — a complex multi-step process (file search, unmarshal, flag overlay, version check) that must complete before generation begins. `RunE` handles orchestration (module discovery, pipeline execution). Separating them allows `PreRunE` to bail early (show help, version mismatch) without polluting the generation logic with config concerns.

### Why templates for section rendering

Go templates provide a declarative way to express output formatting that's far more readable than imperative string building. They support user customization (content templates) without code changes, and the composition model (named sub-templates) maps cleanly to the section-based documentation structure.

### Why inject mode with comment markers

The dominant use case is a README that has both hand-written content (project description, examples, badges) and auto-generated sections (inputs, outputs). Inject mode lets the tool update only its section without touching anything else. The markers are format-agnostic comments (HTML for Markdown, `//` for AsciiDoc) so they're invisible in rendered output.

### Why plugin discovery falls through from built-in

This makes plugins transparent to users. A user adds `formatter: custom` to their config and installs a `tfdocs-format-custom` binary. The CLI tries built-in first (fast path, no subprocess), then falls through to plugins only when needed. No explicit plugin registration or configuration required.

### Why a generator base struct with forEach

All template-based formatters need the same mechanics: render N sections, store them, expose them for content templates. The `generator` base provides this infrastructure so each formatter only defines its template text and any format-specific rendering quirks. The functional options pattern (`withContent`, `withHeader`, etc.) keeps the section-to-field mapping explicit without a fragile positional API.

[go-plugin]: https://github.com/hashicorp/go-plugin
