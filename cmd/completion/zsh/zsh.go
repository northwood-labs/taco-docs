/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package zsh

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCommand registers the "completion zsh" subcommand. Zsh completions use a
// different format than bash (compdef-based), so they require their own
// generation function. The output is typically saved to a file in the user's
// fpath (e.g., /usr/local/share/zsh/site-functions/) and loaded by compinit.
//
// Parent().Parent() is used to generate completions from the root command,
// ensuring the full command tree is represented.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "zsh",
		Short: "Generate shell completion for zsh",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Parent().Parent().GenZshCompletion(os.Stdout)
		},
	}
	return cmd
}
