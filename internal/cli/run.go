/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/terraform-docs/terraform-docs/format"
	"github.com/terraform-docs/terraform-docs/internal/plugin"
	"github.com/terraform-docs/terraform-docs/internal/version"
	pluginsdk "github.com/terraform-docs/terraform-docs/plugin"
	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
)

// Runtime holds the state accumulated during a single CLI invocation. It bridges
// the gap between cobra's stateless command handlers and the multi-step process
// of reading config, discovering modules, and generating output. A single Runtime
// instance is shared across PreRunE and RunE so that config parsing only happens
// once and the resolved state is available to the generation step.
type Runtime struct {
	rootDir string

	formatter string
	config    *print.Config

	cmd           *cobra.Command
	isFlagChanged func(string) bool
}

// NewRuntime returns a new instance of Runtime. A nil config is tolerated so
// callers in tests or plugin contexts don't need to construct a full Config;
// the default will be used instead.
func NewRuntime(config *print.Config) *Runtime {
	if config == nil {
		config = print.DefaultConfig()
	}
	return &Runtime{config: config}
}

// PreRunEFunc is cobra's pre-run hook for all formatter commands. Its purpose is
// to resolve configuration before the actual generation runs. The order matters:
//  1. Extract the formatter name from annotations (avoids fragile string parsing of Use).
//  2. Show help and exit for the root command with no args (UX: don't error, just help).
//  3. Read the config file (with fallback search paths for discoverability).
//  4. Overlay CLI flags onto the config (flags take precedence over file).
//  5. Validate version constraints so mismatched configs fail early with clear messages.
func (r *Runtime) PreRunEFunc(cmd *cobra.Command, args []string) error {
	r.formatter = cmd.Annotations["command"]

	// Root command without arguments means the user ran bare `terraform-docs` —
	// rather than erroring, show help since that's the conventional UX.
	if r.formatter == "root" && len(args) == 0 {
		cmd.Help() //nolint:errcheck,gosec
		os.Exit(0)
	}

	r.isFlagChanged = cmd.Flags().Changed
	r.rootDir = args[0]
	r.cmd = cmd

	// An empty --config value is an explicit user error (they passed -c ""),
	// not a missing-file situation. Fail early with a clear message.
	if r.config.File == "" {
		return fmt.Errorf("value of '--config' can't be empty")
	}

	// Viper handles the layered config file search. A new instance is created
	// per invocation to avoid stale state from prior runs in test scenarios.
	v := viper.New()

	if err := r.readConfig(v, r.config.File, ""); err != nil {
		return err
	}

	// CLI flags override config-file values — this is the standard UX expectation
	// where explicit invocation-time choices win over persistent settings.
	if err := r.unmarshalConfig(v, r.config); err != nil {
		return err
	}

	// Version constraint checking prevents confusing failures when a config file
	// requires features from a newer terraform-docs version than is installed.
	return checkConstraint(r.config.Version, version.Core())
}

// module pairs a filesystem path with its resolved configuration. Each module
// (root or sub) may have its own .terraform-docs.yml that overrides the root config.
type module struct {
	rootDir string
	config  *print.Config
}

// RunEFunc is cobra's main execution hook. It orchestrates the end-to-end flow:
// discover modules (potentially recursively), validate each module's config, and
// generate + write output for each one. The recursive mode exists because many
// Terraform repos contain a top-level module plus submodules under a "modules/"
// directory, and users want all of them documented in a single command.
func (r *Runtime) RunEFunc(cmd *cobra.Command, args []string) error { //nolint:gocyclo
	modules := []module{}

	// Include the main module unless the user explicitly excluded it via config.
	// This gives users control over documenting only submodules when needed.
	if !r.config.Recursive.Enabled || r.config.Recursive.IncludeMain {
		modules = append(modules, module{r.rootDir, r.config})
	}

	// Recursive discovery scans the configured path for Terraform submodules.
	// This only makes sense when output goes to files — stdout output of multiple
	// modules would be impossible to separate.
	if r.config.Recursive.Enabled && r.config.Recursive.Path != "" {
		items, err := r.findSubmodules()
		if err != nil {
			return err
		}

		modules = append(modules, items...)
	}

	for _, module := range modules {
		cfg := r.config

		// Submodules may override root config with their own .terraform-docs.yml,
		// allowing per-module customization (e.g., different sections or formatters).
		if module.config != nil {
			cfg = module.config
		}

		cfg.ModuleRoot = module.rootDir

		if err := cfg.Validate(); err != nil {
			return err
		}

		// Recursive mode without an output file would produce concatenated stdout
		// with no separator between modules — an unusable result. Fail explicitly.
		if r.config.Recursive.Enabled && cfg.Output.File == "" {
			return fmt.Errorf("value of '--output-file' cannot be empty with '--recursive'")
		}

		if err := generateContent(cfg); err != nil {
			return err
		}
	}

	return nil
}

// readConfig searches for and reads a configuration file using viper's multi-path
// lookup. The search order is intentional:
//   - If --config was explicitly provided, use that exact file (fail if missing).
//   - Otherwise, search module root, then CWD, then user home — this supports both
//     per-module configs and global user defaults without requiring any flags.
func (r *Runtime) readConfig(v *viper.Viper, file string, submoduleDir string) error {
	if r.isFlagChanged("config") {
		// User explicitly specified a config file — resolve it absolutely so that
		// relative paths work regardless of CWD vs. module path.
		if absFile, err := filepath.Abs(file); err == nil {
			file = absFile
		}
		v.SetConfigFile(file)
	} else {
		v.SetConfigName(".terraform-docs")
		v.SetConfigType("yml")
	}

	// Search paths are ordered by proximity/specificity: submodule > module root > CWD > HOME.
	// This lets more specific configs shadow broader ones.
	if submoduleDir != "" {
		v.AddConfigPath(submoduleDir)
		v.AddConfigPath(submoduleDir + "/.config")
	}

	v.AddConfigPath(r.rootDir)
	v.AddConfigPath(r.rootDir + "/.config")
	v.AddConfigPath(".")
	v.AddConfigPath(".config")
	v.AddConfigPath("$HOME/.tfdocs.d")

	if err := v.ReadInConfig(); err != nil {
		var perr *os.PathError
		if errors.As(err, &perr) {
			return fmt.Errorf("config file %s not found", file)
		}

		var cerr viper.ConfigFileNotFoundError
		if !errors.As(err, &cerr) {
			return err
		}

		// When no config file exists and the user invoked the root command,
		// there's nothing to do — show help rather than an obscure error.
		if r.formatter == "root" {
			r.cmd.Help() //nolint:errcheck,gosec
			os.Exit(0)
		}
	}

	return nil
}

// unmarshalConfig decodes the viper config into the Config struct and applies
// flag overrides. For explicit subcommands (non-root), the formatter is forced
// to the subcommand name so that a config file's "formatter" field can't
// silently redirect a user's explicit `terraform-docs json` to a different format.
func (r *Runtime) unmarshalConfig(v *viper.Viper, config *print.Config) error {
	r.bindFlags(v)

	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode config, %w", err)
	}

	// When a user explicitly runs e.g. `terraform-docs markdown table ./`, the
	// subcommand name must win over any "formatter:" value in the config file.
	// Only the root command defers to the config file's formatter choice.
	if r.formatter != "root" {
		config.Formatter = r.formatter
	}

	config.Parse()

	return nil
}

// bindFlags synchronizes CLI flag values into the viper config layer. Only
// flags explicitly changed by the user are bound — this ensures that default
// flag values don't accidentally override config-file settings. The special
// handling of show/hide ensures that CLI flags fully replace (not merge with)
// config-file section lists, which matches user expectations of "I said --show
// inputs, so only show inputs."
func (r *Runtime) bindFlags(v *viper.Viper) {
	sectionsCleared := false
	fs := r.cmd.Flags()
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			return
		}

		switch f.Name {
		case "show", "hide":
			// Clear both show and hide from config on first encounter of either
			// flag. This implements "CLI flag replaces file config" semantics —
			// without this, config-file items would persist and merge in unexpected ways.
			if !sectionsCleared {
				v.Set("sections.show", []string{})
				v.Set("sections.hide", []string{})
				sectionsCleared = true
			}

			items, err := fs.GetStringSlice(f.Name)
			if err != nil {
				return
			}
			v.Set(flagMappings[f.Name], items)
		case "sort-by-required", "sort-by-type":
			// Legacy flags that set the sort criteria by their mere presence.
			v.Set("sort.by", flagMappings[f.Name])
		default:
			if _, ok := flagMappings[f.Name]; !ok {
				return
			}
			v.Set(flagMappings[f.Name], f.Value)
		}
	})
}

// mergeConfig creates a copy of the root config and overlays a submodule's
// config onto it. This preserves the root settings as defaults while allowing
// submodules to selectively override specific values.
func (r *Runtime) mergeConfig(v *viper.Viper) (*print.Config, error) {
	copy := *r.config
	merged := &copy

	// If the submodule config specifies section visibility, reset both lists
	// so that the submodule's preferences fully replace (not merge with) the root's.
	if v.IsSet("sections.show") || v.IsSet("sections.hide") {
		merged.Sections.Show = []string{}
		merged.Sections.Hide = []string{}
	}

	if err := r.unmarshalConfig(v, merged); err != nil {
		return nil, err
	}

	return merged, nil
}

// findSubmodules walks the recursive path looking for directories that contain
// Terraform files. It respects exclusion rules and hidden directories (prefixed
// with ".") to avoid scanning vendor/, .terraform/, or other non-module dirs.
func (r *Runtime) findSubmodules() ([]module, error) {
	dir := filepath.Join(r.rootDir, r.config.Recursive.Path)

	if _, err := os.Stat(dir); err != nil {
		// A missing recursive path is not an error — the user may have configured
		// a default "modules" path even though this particular repo doesn't have one.
		if os.IsNotExist(err) {
			return []module{}, nil
		}

		return nil, err
	}

	modules := []module{}

	err := filepath.WalkDir(dir, func(path string, file os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip non-directories and the root of the walk itself.
		if !file.IsDir() || path == dir {
			return nil
		}

		// Skip hidden dirs (like .terraform/) and explicitly excluded directories.
		if strings.HasPrefix(file.Name(), ".") || slices.Contains(r.config.Recursive.Exclude, file.Name()) {
			return filepath.SkipDir
		}

		module, err := r.loadSubModule(path)
		if err != nil {
			return err
		}

		if module != nil {
			modules = append(modules, *module)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return modules, nil
}

// loadSubModule checks whether a directory qualifies as a Terraform module
// (contains .tf files) and, if so, loads any module-specific config. Directories
// without Terraform files are silently skipped to avoid noisy errors for utility
// directories (like shared scripts) that happen to live under the recursive path.
func (r *Runtime) loadSubModule(path string) (*module, error) {
	hasTerraformFiles, err := containsTerraformFiles(path)
	if err != nil {
		return nil, err
	}

	if !hasTerraformFiles {
		return nil, nil
	}

	cfg, err := r.loadModuleConfig(path)
	if err != nil {
		return nil, err
	}

	return &module{rootDir: path, config: cfg}, nil
}

// containsTerraformFiles determines whether a directory has any .tf or .tf.json
// files. This is the minimum signal that a directory is a Terraform module —
// without it, attempting to parse the directory would produce confusing errors.
func containsTerraformFiles(path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasSuffix(name, ".tf") || strings.HasSuffix(name, ".tf.json") {
			return true, nil
		}
	}

	return false, nil
}

// loadModuleConfig attempts to read a submodule-specific config file. If it
// exists, the submodule config is merged with the root config (submodule wins);
// if no config file is found, nil is returned and the root config applies as-is.
func (r *Runtime) loadModuleConfig(path string) (*print.Config, error) {
	var cfg *print.Config

	cfgfile := filepath.Join(path, r.config.File)
	if _, err := os.Stat(cfgfile); !os.IsNotExist(err) {
		v := viper.New()

		if err = r.readConfig(v, cfgfile, path); err != nil {
			return nil, err
		}

		if cfg, err = r.mergeConfig(v); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// checkConstraint enforces the "version" field from .terraform-docs.yml. This
// prevents silent misbehavior when a config file uses features from a newer
// terraform-docs version — instead the user gets an explicit error with the
// constraint and current version, guiding them to upgrade.
func checkConstraint(versionRange string, currentVersion string) error {
	if versionRange == "" {
		return nil
	}

	semver, err := goversion.NewSemver(currentVersion)
	if err != nil {
		return err
	}

	constraint, err := goversion.NewConstraint(versionRange)
	if err != nil || !constraint.Check(semver) {
		return fmt.Errorf("current version: %s, constraints: '%s'", semver, constraint)
	}

	return nil
}

// generateContent is the core pipeline: load Terraform module → format → write.
// It first attempts built-in formatters, then falls through to plugin discovery.
// This two-stage lookup allows the plugin system to extend terraform-docs with
// custom formatters without modifying the core binary.
func generateContent(config *print.Config) error {
	module, err := terraform.LoadWithOptions(config)
	if err != nil {
		return err
	}

	formatter, err := format.New(config)
	// If the formatter name isn't recognized by built-in formatters, try plugins.
	// This fallback design means plugins are transparent to users — they just
	// specify a formatter name and it works whether built-in or plugin-provided.
	if err != nil {
		plugins, perr := plugin.Discover()
		if perr != nil {
			return fmt.Errorf("formatter '%s' not found", config.Formatter)
		}

		client, found := plugins.Get(config.Formatter)
		if !found {
			return fmt.Errorf("formatter '%s' not found", config.Formatter)
		}

		content, cerr := client.Execute(&pluginsdk.ExecuteArgs{
			Module: module,
			Config: config,
		})
		if cerr != nil {
			return cerr
		}

		return writeContent(config, content)
	}

	err = formatter.Generate(module)
	if err != nil {
		return err
	}

	// Render applies the user's custom content template (if provided) to the
	// generated sections. An empty template means "use default section order."
	content, err := formatter.Render(config.Content)
	if err != nil {
		return err
	}

	return writeContent(config, content)
}

// writeContent directs generated documentation to either stdout or a file,
// depending on whether --output-file was specified. This abstraction allows
// the generation pipeline to be output-destination-agnostic.
func writeContent(config *print.Config, content string) error {
	var w io.Writer

	if config.Output.File != "" {
		// File writing supports both inject (between markers) and replace modes,
		// enabling users to maintain hand-written content around the generated sections.
		w = &fileWriter{
			file: config.Output.File,
			dir:  config.ModuleRoot,

			mode: config.Output.Mode,

			check: config.Output.Check,

			template: config.Output.Template,
			begin:    config.Output.BeginComment,
			end:      config.Output.EndComment,
		}
	} else {
		w = &stdoutWriter{}
	}

	_, err := io.WriteString(w, content)

	return err
}
