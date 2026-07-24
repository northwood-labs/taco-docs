// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package main

import (
	"os"

	"github.com/terraform-docs/terraform-docs/cmd"
)

// main is the process entry point. It delegates immediately to cmd.Execute()
// which builds the CLI tree and runs the appropriate command. The only
// responsibility here is translating a non-nil error into a non-zero exit code
// so that CI pipelines and shell scripts can detect failures.
func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
