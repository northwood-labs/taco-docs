// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package tfvars

import (
	"github.com/spf13/cobra"

	"github.com/northwood-labs/taco-docs/cmd/tfvars/hcl"
	"github.com/northwood-labs/taco-docs/cmd/tfvars/json"
	"github.com/northwood-labs/taco-docs/internal/cli"
	"github.com/northwood-labs/taco-docs/print"
)

// NewCommand registers the "tfvars" subcommand group. Unlike documentation
// formatters, tfvars output is meant to be used directly as Terraform input —
// it generates a terraform.tfvars file skeleton from a module's declared
// variables. This helps users quickly scaffold the required variable values
// for a module they're consuming, rather than manually reading docs and
// creating the file by hand.
//
// Note: this command has no PreRunE/RunE of its own because it requires an
// explicit sub-format choice (HCL or JSON) — there is no meaningful default.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "tfvars [PATH]",
		Short:       "Generate terraform.tfvars of inputs",
		Annotations: cli.Annotations("tfvars"),
	}

	// subcommands — HCL for native .tfvars format, JSON for .tfvars.json
	cmd.AddCommand(hcl.NewCommand(runtime, config))
	cmd.AddCommand(json.NewCommand(runtime, config))

	return cmd
}
