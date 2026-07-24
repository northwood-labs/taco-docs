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

// NewCommand registers the "tfvars json" subcommand. JSON is an alternative
// to HCL for providing variable values (via .tfvars.json files). This format
// is particularly useful when variable values are programmatically generated
// by scripts or CI pipelines that naturally produce JSON output.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "json [PATH]",
		Short:       "Generate JSON format of terraform.tfvars of inputs",
		Annotations: cli.Annotations("tfvars json"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}
	return cmd
}
