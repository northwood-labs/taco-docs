// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/northwood-labs/taco-docs/cmd/asciidoc"
	"github.com/northwood-labs/taco-docs/cmd/completion"
	"github.com/northwood-labs/taco-docs/cmd/json"
	"github.com/northwood-labs/taco-docs/cmd/markdown"
	"github.com/northwood-labs/taco-docs/cmd/pretty"
	"github.com/northwood-labs/taco-docs/cmd/tfvars"
	"github.com/northwood-labs/taco-docs/cmd/toml"
	versioncmd "github.com/northwood-labs/taco-docs/cmd/version"
	"github.com/northwood-labs/taco-docs/cmd/xml"
	"github.com/northwood-labs/taco-docs/cmd/yaml"
	"github.com/northwood-labs/taco-docs/internal/cli"
	"github.com/northwood-labs/taco-docs/internal/version"
	"github.com/northwood-labs/taco-docs/print"
)

// Execute is the top-level entry point for the CLI, called by main.main().
// It constructs the full command tree and delegates execution to cobra. Errors
// are written to stderr so they remain visible even when stdout is redirected
// to a file (a common use case for documentation generation).
func Execute() error {
	if err := NewCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return err
	}
	return nil
}

// NewCommand constructs the root cobra.Command with all persistent flags and
// subcommands attached. The root command doubles as the default formatter
// (resolved via config file) when invoked without an explicit subcommand. A
// shared Config and Runtime are created here so that flag values flow through
// a single source of truth — every subcommand shares the same config pointer,
// which means persistent flags set at the root level propagate to all children
// without manual wiring.
func NewCommand() *cobra.Command {
	config := print.DefaultConfig()
	runtime := cli.NewRuntime(config)
	cmd := &cobra.Command{
		Args:          cobra.MaximumNArgs(1),
		Use:           "terraform-docs [PATH]",
		Short:         "A utility to generate documentation from Terraform modules in various output formats",
		Long:          "A utility to generate documentation from Terraform modules in various output formats",
		Version:       version.Full(),
		SilenceUsage:  true, // Prevent cobra from printing usage on every error — users only need usage on --help.
		SilenceErrors: true, // We handle error formatting ourselves in Execute() for consistent stderr output.
		Annotations:   cli.Annotations("root"),
		PreRunE:       runtime.PreRunEFunc,
		RunE:          runtime.RunEFunc,
	}

	// Persistent flags are defined here because they apply universally to all
	// output formatters — they control which Terraform module to read, how to
	// discover submodules, which sections to include, and how to write output.
	cmd.PersistentFlags().StringVarP(&config.File, "config", "c", ".terraform-docs.yml", "config file name")
	cmd.PersistentFlags().
		BoolVar(&config.Recursive.Enabled, "recursive", false, "update submodules recursively (default false)")
	cmd.PersistentFlags().
		StringVar(&config.Recursive.Path, "recursive-path", "modules", "submodules path to recursively update")
	cmd.PersistentFlags().
		BoolVar(&config.Recursive.IncludeMain, "recursive-include-main", true, "include the main module")
	cmd.PersistentFlags().
		StringSliceVar(&config.Recursive.Exclude, "recursive-exclude", []string{}, "exclude directories from recursive update")

	// Show/hide flags let users cherry-pick documentation sections without
	// editing the config file — useful for CI pipelines that need different
	// outputs from the same module.
	cmd.PersistentFlags().
		StringSliceVar(&config.Sections.Show, "show", []string{}, "show section ["+print.AllSections+"]")
	cmd.PersistentFlags().
		StringSliceVar(&config.Sections.Hide, "hide", []string{}, "hide section ["+print.AllSections+"]")

	// Output flags control file-writing behavior, enabling in-place README
	// updates (inject mode) or full file replacement (replace mode).
	cmd.PersistentFlags().
		StringVar(&config.Output.File, "output-file", "", "file path to insert output into (default \"\")")
	cmd.PersistentFlags().
		StringVar(&config.Output.Mode, "output-mode", "inject", "output to file method ["+print.OutputModes+"]")
	cmd.PersistentFlags().StringVar(&config.Output.Template, "output-template", print.OutputTemplate, "output template")
	cmd.PersistentFlags().
		BoolVar(&config.Output.Check, "output-check", false, "check if content of output file is up to date (default false)")

	cmd.PersistentFlags().BoolVar(&config.Sort.Enabled, "sort", true, "sort items")
	cmd.PersistentFlags().StringVar(&config.Sort.By, "sort-by", "name", "sort items by criteria ["+print.SortTypes+"]")

	cmd.PersistentFlags().
		StringVar(&config.HeaderFrom, "header-from", "main.tf", "relative path of a file to read header from")
	cmd.PersistentFlags().
		StringVar(&config.FooterFrom, "footer-from", "", "relative path of a file to read footer from (default \"\")")

	cmd.PersistentFlags().BoolVar(&config.Settings.LockFile, "lockfile", true, "read .terraform.lock.hcl if exist")

	// Output-values flags enable injecting actual Terraform output values
	// (from `terraform output -json`) into the documentation, giving readers
	// visibility into current state alongside the schema.
	cmd.PersistentFlags().
		BoolVar(&config.OutputValues.Enabled, "output-values", false, "inject output values into outputs (default false)")
	cmd.PersistentFlags().
		StringVar(&config.OutputValues.From, "output-values-from", "", "inject output values from file into outputs (default \"\")")

	cmd.PersistentFlags().
		BoolVar(&config.Settings.ReadComments, "read-comments", true, "use comments as description when description is empty")

	// Each formatter subcommand represents a distinct output format. They share
	// the same runtime and config so that flag values and config-file settings
	// are consistent regardless of which format is selected.
	cmd.AddCommand(asciidoc.NewCommand(runtime, config))
	cmd.AddCommand(json.NewCommand(runtime, config))
	cmd.AddCommand(markdown.NewCommand(runtime, config))
	cmd.AddCommand(pretty.NewCommand(runtime, config))
	cmd.AddCommand(tfvars.NewCommand(runtime, config))
	cmd.AddCommand(toml.NewCommand(runtime, config))
	cmd.AddCommand(xml.NewCommand(runtime, config))
	cmd.AddCommand(yaml.NewCommand(runtime, config))

	// Utility subcommands that don't produce documentation output.
	cmd.AddCommand(completion.NewCommand())
	cmd.AddCommand(versioncmd.NewCommand())

	return cmd
}
