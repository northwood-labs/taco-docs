// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package terraform

import (
	"fmt"
	"sort"
	"strings"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

// Resource represents a managed or data type that is created by the module.
//
// WHY: Documenting every resource a module manages helps consumers understand blast radius
// and required provider permissions without reading the full HCL. The struct carries enough
// metadata (provider source, version, mode) to generate registry hyperlinks and human-readable
// resource addresses.
type Resource struct {
	Type           string       `json:"type"        toml:"type"        xml:"type"        yaml:"type"`
	Name           string       `json:"name"        toml:"name"        xml:"name"        yaml:"name"`
	ProviderName   string       `json:"provider"    toml:"provider"    xml:"provider"    yaml:"provider"`
	ProviderSource string       `json:"source"      toml:"source"      xml:"source"      yaml:"source"`
	Mode           string       `json:"mode"        toml:"mode"        xml:"mode"        yaml:"mode"`
	Version        types.String `json:"version"     toml:"version"     xml:"version"     yaml:"version"`
	Description    types.String `json:"description" toml:"description" xml:"description" yaml:"description"`
	Position       Position     `json:"-"           toml:"-"           xml:"-"           yaml:"-"`
}

// Spec returns the resource spec addresses a specific resource in the config.
// It takes the form: resource_type.resource_name[resource index]
// For more details, see:
// https://www.terraform.io/docs/cli/state/resource-addressing.html#resource-spec
//
// WHY: Reassembling the full type (provider_name + "_" + short_type) produces the canonical
// resource address that Terraform CLI uses, making documentation consistent with plan/state output.
func (r *Resource) Spec() string {
	return r.ProviderName + "_" + r.Type + "." + r.Name
}

// GetMode returns normalized resource type as "resource" or "data source".
//
// WHY: Terraform's internal mode strings ("managed", "data") are implementation detail; users
// and documentation refer to these as "resource" and "data source". This translation keeps
// generated docs aligned with Terraform's official terminology.
func (r *Resource) GetMode() string {
	switch r.Mode {
	case "managed":
		return "resource"
	case "data":
		return "data source"
	default:
		return "invalid"
	}
}

// URL returns a best guess at the URL for resource documentation.
//
// WHY: Auto-generating registry links means users can click through to full resource docs
// without searching manually. The slash-count guard prevents generating broken URLs for
// private or non-standard registry sources whose path structure differs from the public registry.
func (r *Resource) URL() string {
	kind := ""
	switch r.Mode {
	case "managed":
		kind = "resources"
	case "data":
		kind = "data-sources"
	default:
		return ""
	}

	if strings.Count(r.ProviderSource, "/") > 1 {
		return ""
	}
	return fmt.Sprintf(
		"https://registry.terraform.io/providers/%s/%s/docs/%s/%s",
		r.ProviderSource,
		r.Version,
		kind,
		r.Type,
	)
}

func sortResourcesByType(x []*Resource) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Mode == x[j].Mode {
			if x[i].Spec() == x[j].Spec() {
				return x[i].Name <= x[j].Name
			}
			return x[i].Spec() < x[j].Spec()
		}
		return x[i].Mode > x[j].Mode
	})
}

type resources []*Resource

// WHY: Resources are always sorted by type regardless of the user's sort preference. This
// groups managed resources and data sources separately, then alphabetizes within each group
// by their Spec() address—producing a predictable, scannable list in documentation.
func (rr resources) sort(enabled bool, by string) { //nolint:unparam
	// always sort by type
	sortResourcesByType(rr)
}
