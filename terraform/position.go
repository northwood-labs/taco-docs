/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

// Position represents position of Terraform item (input, output, provider, etc) in a file.
//
// WHY: Position tracks the original source-file location of each declaration so the sort-by-position
// strategy can preserve the author's intended ordering. JSON/TOML/XML/YAML tags are all "-" because
// position is internal metadata for sorting—it is not documentation content and should never appear
// in generated output.
type Position struct {
	Filename string `json:"-" toml:"-" xml:"-" yaml:"-"`
	Line     int    `json:"-" toml:"-" xml:"-" yaml:"-"`
}
