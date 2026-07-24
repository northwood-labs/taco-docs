/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package pretty

import (
	"github.com/spf13/cobra"

	"github.com/terraform-docs/terraform-docs/internal/cli"
	"github.com/terraform-docs/terraform-docs/print"
)

// NewCommand registers the "pretty" subcommand. The pretty formatter produces
// human-readable, colorized terminal output — it exists primarily for quick
// local inspection of a module's interface without leaving the terminal. The
// color flag allows disabling ANSI escape codes when output is piped or
// captured in environments that don't support them.
func NewCommand(runtime *cli.Runtime, config *print.Config) *cobra.Command {
	cmd := &cobra.Command{
		Args:        cobra.ExactArgs(1),
		Use:         "pretty [PATH]",
		Short:       "Generate colorized pretty of inputs and outputs",
		Annotations: cli.Annotations("pretty"),
		PreRunE:     runtime.PreRunEFunc,
		RunE:        runtime.RunEFunc,
	}

	// flags
	cmd.PersistentFlags().BoolVar(&config.Settings.Color, "color", true, "colorize printed result")

	return cmd
}
