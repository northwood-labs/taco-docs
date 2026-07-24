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
	"strings"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

// ProviderFunction represents a provider function call that is used by the module.
//
// WHY: Provider functions (provider::aws::arn_parse, etc.) are a newer Terraform/OpenTofu feature
// that lets providers expose pure functions. Since they are not resources or data sources, they
// need their own model so the documentation can list which provider-supplied functions a module
// relies on—helping consumers understand SDK dependencies beyond just CRUD resources.
type ProviderFunction struct {
	ProviderName   string       `json:"provider" toml:"provider" xml:"provider" yaml:"provider"`
	Function       string       `json:"function" toml:"function" xml:"function" yaml:"function"`
	ProviderSource string       `json:"source"   toml:"source"   xml:"source"   yaml:"source"`
	Version        types.String `json:"version"  toml:"version"  xml:"version"  yaml:"version"`
	Position       Position     `json:"-"        toml:"-"        xml:"-"        yaml:"-"`
}

// Spec returns the provider function spec in the same syntax as Terraform provider function calls.
//
// WHY: Reconstructing the canonical "provider::name::function" syntax gives documentation the
// exact identifier users would type in their HCL, making it immediately copy-pasteable.
func (p *ProviderFunction) Spec() string {
	return fmt.Sprintf("provider::%s::%s", p.ProviderName, p.Function)
}

// URL returns a best guess at the URL for provider function documentation.
//
// WHY: Linking directly to registry docs saves readers a manual search. The slash-count guard
// filters out non-registry sources (e.g., private registries with deeper paths) where the
// standard URL pattern would be invalid.
func (p *ProviderFunction) URL() string {
	if strings.Count(p.ProviderSource, "/") > 1 {
		return ""
	}
	return fmt.Sprintf(
		"https://registry.terraform.io/providers/%s/%s/docs/functions/%s",
		p.ProviderSource,
		p.Version,
		p.Function,
	)
}

func sortProviderFunctionsByName(x []*ProviderFunction) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].ProviderName == x[j].ProviderName {
			return x[i].Function < x[j].Function
		}
		return x[i].ProviderName < x[j].ProviderName
	})
}

type providerFunctions []*ProviderFunction

// WHY: Provider functions are always sorted by provider name then function name. Since they
// are a relatively new feature with no "required" or "type" dimension, a single deterministic
// ordering keeps documentation stable across regeneration runs.
func (pf providerFunctions) sort(enabled bool, by string) { //nolint:unparam
	sortProviderFunctionsByName(pf)
}
