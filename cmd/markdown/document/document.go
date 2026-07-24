// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package document

import (
	"github.com/spf13/cobra"

	"github.com/terraform-docs/terraform-docs/internal/cli"
	"github.com/terraform-docs/terraform-docs/print"
)

// NewCommand registers the "markdown document" subcommand. The document format
// renders each input/output as its own subsection with a heading, rather than
// compressing everything into a table. This is preferred when modules have
// lengthy descriptions, complex default values, or when the documentation will
// be read on platforms that poorly render wide Markdown tables.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "document [PATH]",
		Aliases:     []string{"doc"},
		Short:       "Generate Markdown document of inputs and outputs",
		Annotations: cli.Annotations("markdown document"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}
	return cmd
}
