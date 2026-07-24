// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package json

import (
	"github.com/spf13/cobra"

	"github.com/terraform-docs/terraform-docs/internal/cli"
	"github.com/terraform-docs/terraform-docs/print"
)

// NewCommand registers the "json" subcommand. JSON output is useful for
// programmatic consumption — CI pipelines, linters, or custom rendering
// tools can parse the structured output without depending on fragile text
// parsing of Markdown or AsciiDoc. The "escape" flag controls whether HTML
// entities are escaped in the JSON encoder, which matters when the JSON
// will later be embedded inside HTML pages.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "json [PATH]",
		Short:       "Generate JSON of inputs and outputs",
		Annotations: cli.Annotations("json"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&config.Settings.Escape, "escape", true, "escape special characters")

	return cmd
}
