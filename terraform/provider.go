// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package terraform

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/northwood-labs/taco-docs/internal/types"
)

type (
	// Provider represents a Terraform provider used by the module.
	//
	// Providers are discovered from actual resource usage (not just
	// required_providers) so the documentation reflects what the module truly
	// depends on at runtime. The Alias field exists because Terraform allows
	// multiple configurations of the same provider (e.g., aws.us-east), and
	// documentation needs to distinguish them.
	Provider struct {
		Name     string       `json:"name"    toml:"name"    xml:"name"    yaml:"name"`
		Alias    types.String `json:"alias"   toml:"alias"   xml:"alias"   yaml:"alias"`
		Version  types.String `json:"version" toml:"version" xml:"version" yaml:"version"`
		Position Position     `json:"-"       toml:"-"       xml:"-"       yaml:"-"`
	}

	providers []*Provider
)

// FullName returns full name of the provider, with alias if available.
//
// Terraform uses "name.alias" notation (e.g., aws.us-east) to reference aliased
// provider configurations. FullName reconstructs this so documentation displays
// the same identifier that users write in their provider meta-arguments.
func (p *Provider) FullName() string {
	if p.Alias != "" {
		return fmt.Sprintf("%s.%s", p.Name, p.Alias)
	}

	return p.Name
}

func sortProvidersByName(x []*Provider) {
	slices.SortFunc(x, func(a, b *Provider) int {
		if a.Name == b.Name {
			return strings.Compare(string(a.Alias), string(b.Alias))
		}

		return strings.Compare(a.Name, b.Name)
	})
}

func sortProvidersByPosition(x []*Provider) {
	slices.SortFunc(x, func(a, b *Provider) int {
		if a.Position.Filename == b.Position.Filename {
			if a.Position.Line == b.Position.Line {
				return strings.Compare(a.FullName(), b.FullName())
			}

			return cmp.Compare(a.Position.Line, b.Position.Line)
		}

		return strings.Compare(a.Position.Filename, b.Position.Filename)
	})
}

// When sorting is disabled, position-based ordering preserves the author's
// original file layout. When enabled, alphabetical by name is the only strategy
// for providers since there's no meaningful "type" or "required" axis for
// providers the way there is for inputs.
func (pp providers) sort(enabled bool, _ string) { // lint:allow_param lint:allow_control_coupling_antipattern
	if !enabled {
		sortProvidersByPosition(pp)
	} else {
		// always sort by name if sorting is enabled.
		sortProvidersByName(pp)
	}
}
