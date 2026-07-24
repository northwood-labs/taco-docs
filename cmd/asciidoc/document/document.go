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

	"github.com/northwood-labs/taco-docs/internal/cli"
	"github.com/northwood-labs/taco-docs/print"
)

// NewCommand registers the "asciidoc document" subcommand. The document format
// renders each input/output as a dedicated subsection rather than a table row.
// This is better suited for modules with verbose descriptions or complex nested
// default values that would overflow a table cell in AsciiDoc.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "document [PATH]",
		Aliases:     []string{"doc"},
		Short:       "Generate AsciiDoc document of inputs and outputs",
		Annotations: cli.Annotations("asciidoc document"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}
	return cmd
}
