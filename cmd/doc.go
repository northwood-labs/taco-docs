// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

// Package cmd implements the CLI command tree for taco-docs. It wires together
// all output-format subcommands, persistent flags, and utility commands
// (completion, version) under a single root cobra.Command.
package cmd
