// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package xml

import (
	"github.com/spf13/cobra"

	"github.com/northwood-labs/taco-docs/internal/cli"
	"github.com/northwood-labs/taco-docs/print"
)

// NewCommand registers the "xml" subcommand. XML output is provided for
// enterprise environments and legacy systems that require XML-based data
// exchange (e.g., XSLT transformations, SOAP integrations, or documentation
// systems that consume XML). No format-specific flags are needed since XML
// encoding is fully deterministic for the module data model.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "xml [PATH]",
		Short:       "Generate XML of inputs and outputs",
		Annotations: cli.Annotations("xml"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}
	return cmd
}
