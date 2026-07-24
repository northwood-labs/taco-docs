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

// NewCommand registers the "markdown table" subcommand. The table format is
// the most popular Markdown output because it presents inputs and outputs in
// a compact, scannable grid. It renders well on GitHub/GitLab and allows
// reviewers to quickly compare names, types, defaults, and descriptions
// side by side — which is exactly what teams need during code review.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "table [PATH]",
		Aliases:     []string{"tbl"},
		Short:       "Generate Markdown tables of inputs and outputs",
		Annotations: cli.Annotations("markdown table"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}
	return cmd
}
