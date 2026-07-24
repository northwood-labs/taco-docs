/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/terraform-docs/terraform-docs/internal/plugin"
	"github.com/terraform-docs/terraform-docs/internal/version"
)

// NewCommand registers the "version" subcommand. Beyond printing the core
// binary version, it also discovers and reports installed formatter plugins.
// This is essential for debugging — when users report issues, knowing both the
// core version and plugin versions helps reproduce and diagnose problems.
// Plugin discovery errors are silently ignored here because a missing plugin
// directory is a normal state (most users don't use plugins).
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "version",
		Short: "Print the version number of terraform-docs",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("terraform-docs version %s\n", version.Full())

			// Attempt to list installed plugins so users see a complete picture
			// of their terraform-docs installation in a single command.
			plugins, err := plugin.Discover()
			if err != nil {
				return
			}
			for _, f := range plugins.All() {
				name, err := f.Name()
				if err != nil {
					name = "unknown"
				}
				version, err := f.Version()
				if err != nil {
					version = "unknown"
				}
				fmt.Printf("- plugin %s %s\n", name, version)
			}
		},
	}
	return cmd
}
