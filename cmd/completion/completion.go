// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package completion

import (
	"github.com/spf13/cobra"

	"github.com/northwood-labs/taco-docs/cmd/completion/bash"
	"github.com/northwood-labs/taco-docs/cmd/completion/fish"
	"github.com/northwood-labs/taco-docs/cmd/completion/zsh"
)

// NewCommand registers the "completion" subcommand group. Shell completions
// dramatically improve the developer experience for CLI tools by enabling
// tab-completion of subcommands and flags. This avoids the need to constantly
// reference --help or documentation. Each shell has its own completion syntax,
// so separate subcommands handle bash, zsh, and fish respectively.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "completion SHELL",
		Short: "Generate shell completion code for the specified shell (bash, zsh, fish)",
		Long:  longDescription,
	}

	// subcommands
	cmd.AddCommand(bash.NewCommand())
	cmd.AddCommand(zsh.NewCommand())
	cmd.AddCommand(fish.NewCommand())

	return cmd
}

const longDescription = `Outputs terraform-docs shell completion for the given shell (bash, zsh, fish)
This depends on the bash-completion binary.  Example installation instructions:
# for bash users
	$ terraform-docs completion bash > ~/.terraform-docs-completion
	$ source ~/.terraform-docs-completion

	# or the one-liner below

	$ source <(terraform-docs completion bash)

# for zsh users
	% terraform-docs completion zsh > /usr/local/share/zsh/site-functions/_terraform-docs
	% autoload -U compinit && compinit
# or if zsh-completion is installed via homebrew
	% terraform-docs completion zsh > "${fpath[1]}/_terraform-docs"

# for ohmyzsh
	$ terraform-docs completion zsh > ~/.oh-my-zsh/completions/_terraform-docs
	$ omz reload

# for fish users
	$ terraform-docs completion fish > ~/.config/fish/completions/terraform-docs.fish

Additionally, you may want to output the completion to a file and source in your .bashrc
Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2
`
