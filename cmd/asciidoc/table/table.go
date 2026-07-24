// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package table

import (
	"github.com/spf13/cobra"

	"github.com/northwood-labs/taco-docs/internal/cli"
	"github.com/northwood-labs/taco-docs/print"
)

// NewCommand registers the "asciidoc table" subcommand. The table format
// arranges inputs and outputs in AsciiDoc table syntax for compact presentation.
// This is the preferred AsciiDoc format when brevity matters and descriptions
// are short enough to fit in table cells without degrading readability.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "table [PATH]",
		Aliases:     []string{"tbl"},
		Short:       "Generate AsciiDoc tables of inputs and outputs",
		Annotations: cli.Annotations("asciidoc table"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}
	return cmd
}
