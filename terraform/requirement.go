/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"github.com/terraform-docs/terraform-docs/internal/types"
)

// Requirement represents a requirement for Terraform module.
//
// WHY: This minimal struct captures only name, registry source, and version constraint—the three
// pieces needed to render the "Requirements" documentation section. It covers both the Terraform
// core version constraint and each required provider, unifying them into a single list so
// formatters can iterate once without distinguishing the two categories.
type Requirement struct {
	Name    string       `json:"name"    toml:"name"    xml:"name"    yaml:"name"`
	Source  types.String `json:"source"  toml:"source"  xml:"source"  yaml:"source"`
	Version types.String `json:"version" toml:"version" xml:"version" yaml:"version"`
}
