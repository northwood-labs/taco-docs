/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package bash

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCommand registers the "completion bash" subcommand. It leverages cobra's
// built-in bash completion generator which produces a script that integrates
// with bash-completion. The script is written to stdout so users can redirect
// it to a file or source it directly in their shell session.
//
// We navigate to the root command via Parent().Parent() because the completion
// script needs the full command tree to generate completions for all subcommands
// and flags — generating from a leaf command would only produce partial output.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "bash",
		Short: "Generate shell completion for bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Parent().Parent().GenBashCompletion(os.Stdout)
		},
	}
	return cmd
}
