// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package cli

// Annotations returns a standard set of cobra command annotations. These
// annotations serve as metadata that the Runtime uses to determine which
// formatter to invoke — rather than parsing the command's Use string (which
// would be fragile when aliases are involved), PreRunEFunc reads the "command"
// annotation to know what format was requested. The "kind" annotation categorizes
// the command for potential future use in help grouping or plugin discovery.
func Annotations(cmd string) map[string]string {
	return map[string]string{
		"command": cmd,
		"kind":    "formatter",
	}
}
