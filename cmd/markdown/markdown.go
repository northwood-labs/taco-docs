/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package markdown

import (
	"github.com/spf13/cobra"

	"github.com/terraform-docs/terraform-docs/cmd/markdown/document"
	"github.com/terraform-docs/terraform-docs/cmd/markdown/table"
	"github.com/terraform-docs/terraform-docs/internal/cli"
	"github.com/terraform-docs/terraform-docs/print"
)

// NewCommand registers the "markdown" subcommand and its children ("document"
// and "table"). Markdown is the most commonly used output format because it
// renders natively on GitHub, GitLab, Bitbucket, and most documentation sites.
// The parent command defaults to the "markdown table" format when invoked
// directly (via annotations), providing backward compatibility.
//
// Flags defined here are specific to Markdown rendering concerns — anchor
// generation, ATX header style, column visibility, HTML usage, and indentation.
// They are persistent so that both the "document" and "table" sub-formats
// inherit them without redeclaration.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "markdown [PATH]",
		Aliases:     []string{"md"},
		Short:       "Generate Markdown of inputs and outputs",
		Annotations: cli.Annotations("markdown"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}

	// Anchor links allow cross-referencing within the generated docs — users
	// can link directly to a specific input or output from elsewhere in a README.
	cmd.PersistentFlags().BoolVar(&config.Settings.Anchor, "anchor", true, "create anchor links")

	// ATX-closed headers (e.g., "## Title ##") satisfy certain Markdown linters
	// that require the closing hashes for consistency.
	cmd.PersistentFlags().BoolVar(&config.Settings.AtxClosed, "atx-closed", false, "close ATX style headers")

	// Column visibility flags let users suppress columns that add noise for
	// their specific use case (e.g., hiding "Default" when all inputs are required).
	cmd.PersistentFlags().BoolVar(&config.Settings.Default, "default", true, "show Default column or section")
	cmd.PersistentFlags().BoolVar(&config.Settings.Required, "required", true, "show Required column or section")
	cmd.PersistentFlags().BoolVar(&config.Settings.Sensitive, "sensitive", true, "show Sensitive column or section")
	cmd.PersistentFlags().BoolVar(&config.Settings.Type, "type", true, "show Type column or section")

	// Escape controls whether special Markdown characters (like underscores in
	// variable names) are backslash-escaped to prevent unintended formatting.
	cmd.PersistentFlags().BoolVar(&config.Settings.Escape, "escape", true, "escape special characters")

	// HTML flag controls whether HTML tags (like <br/>) are used for multi-line
	// content in tables. Disabling this is necessary for renderers that strip HTML.
	cmd.PersistentFlags().BoolVar(&config.Settings.HTML, "html", true, "use HTML tags in generated output")

	// Hide-empty suppresses sections with no items, keeping the output concise
	// for modules that don't use all Terraform features (e.g., no providers).
	cmd.PersistentFlags().
		BoolVar(&config.Settings.HideEmpty, "hide-empty", false, "hide empty sections (default false)")

	// Indent controls the heading depth — allows embedding generated docs at
	// any nesting level within a larger document structure.
	cmd.PersistentFlags().
		IntVar(&config.Settings.Indent, "indent", 2, "indentation level of Markdown sections [1, 2, 3, 4, 5]")

	// Sub-formats: "document" renders each section as prose with headers,
	// while "table" renders inputs/outputs as Markdown tables for compact viewing.
	cmd.AddCommand(document.NewCommand(runtime, config))
	cmd.AddCommand(table.NewCommand(runtime, config))

	return cmd
}
