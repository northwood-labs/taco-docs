/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import (
	"fmt"
	"sort"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

// Provider represents a Terraform provider used by the module.
//
// WHY: Providers are discovered from actual resource usage (not just required_providers) so the
// documentation reflects what the module truly depends on at runtime. The Alias field exists
// because Terraform allows multiple configurations of the same provider (e.g., aws.us-east),
// and documentation needs to distinguish them.
type Provider struct {
	Name     string       `json:"name"    toml:"name"    xml:"name"    yaml:"name"`
	Alias    types.String `json:"alias"   toml:"alias"   xml:"alias"   yaml:"alias"`
	Version  types.String `json:"version" toml:"version" xml:"version" yaml:"version"`
	Position Position     `json:"-"       toml:"-"       xml:"-"       yaml:"-"`
}

// FullName returns full name of the provider, with alias if available.
//
// WHY: Terraform uses "name.alias" notation (e.g., aws.us-east) to reference aliased provider
// configurations. FullName reconstructs this so documentation displays the same identifier
// that users write in their provider meta-arguments.
func (p *Provider) FullName() string {
	if p.Alias != "" {
		return fmt.Sprintf("%s.%s", p.Name, p.Alias)
	}
	return p.Name
}

func sortProvidersByName(x []*Provider) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Name == x[j].Name {
			return x[i].Name == x[j].Name && x[i].Alias < x[j].Alias
		}
		return x[i].Name < x[j].Name
	})
}

func sortProvidersByPosition(x []*Provider) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Position.Filename == x[j].Position.Filename {
			if x[i].Position.Line == x[j].Position.Line {
				return x[i].FullName() < x[j].FullName()
			}
			return x[i].Position.Line < x[j].Position.Line
		}
		return x[i].Position.Filename < x[j].Position.Filename
	})
}

type providers []*Provider

// WHY: When sorting is disabled, position-based ordering preserves the author's original file
// layout. When enabled, alphabetical by name is the only strategy for providers since there's
// no meaningful "type" or "required" axis for providers the way there is for inputs.
func (pp providers) sort(enabled bool, by string) { //nolint:unparam
	if !enabled {
		sortProvidersByPosition(pp)
	} else {
		// always sort by name if sorting is enabled
		sortProvidersByName(pp)
	}
}
