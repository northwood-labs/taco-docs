/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package fish

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCommand registers the "completion fish" subcommand. Fish shell uses its own
// completion format and stores scripts in ~/.config/fish/completions/. The second
// argument (true) to GenFishCompletion includes command descriptions alongside
// completion candidates, providing richer context during tab-completion.
//
// Parent().Parent() is used to generate completions from the root command,
// ensuring the full command tree is represented.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "fish",
		Short: "Generate shell completion for fish",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Parent().Parent().GenFishCompletion(os.Stdout, true)
		},
	}
	return cmd
}
