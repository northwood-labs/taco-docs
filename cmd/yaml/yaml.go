// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package yaml

import (
	"github.com/spf13/cobra"

	"github.com/northwood-labs/taco-docs/internal/cli"
	"github.com/northwood-labs/taco-docs/print"
)

// NewCommand registers the "yaml" subcommand. YAML output provides a
// machine-readable yet human-friendly serialization of the module interface.
// It serves a similar programmatic role as JSON but is preferred in ecosystems
// that already use YAML-heavy tooling (e.g., Ansible, Helm, or custom
// documentation pipelines). No format-specific flags are needed because YAML
// has no ambiguous encoding choices like HTML escaping.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "yaml [PATH]",
		Short:       "Generate YAML of inputs and outputs",
		Annotations: cli.Annotations("yaml"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}
	return cmd
}
