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

	"github.com/terraform-docs/terraform-docs/internal/types"
	"github.com/terraform-docs/terraform-docs/print"
)

// ModuleCall represents a submodule called by Terraform module.
//
// WHY: Module calls document the composition graph—which child modules this module orchestrates.
// Consumers need to know the source (registry, git, local path) and pinned version to audit
// supply-chain dependencies without reading HCL directly.
type ModuleCall struct {
	Name        string       `json:"name"        toml:"name"        xml:"name"        yaml:"name"`
	Source      string       `json:"source"      toml:"source"      xml:"source"      yaml:"source"`
	Version     string       `json:"version"     toml:"version"     xml:"version"     yaml:"version"`
	Description types.String `json:"description" toml:"description" xml:"description" yaml:"description"`
	Position    Position     `json:"-"           toml:"-"           xml:"-"           yaml:"-"`
}

// FullName returns full name of the modulecall, with version if available.
//
// WHY: Including the version in the display name disambiguates multiple calls to the same source
// at different versions—a common pattern when migrating modules incrementally.
func (mc *ModuleCall) FullName() string {
	if mc.Version != "" {
		return fmt.Sprintf("%s,%s", mc.Source, mc.Version)
	}
	return mc.Source
}

func sortModulecallsByName(x []*ModuleCall) {
	sort.Slice(x, func(i, j int) bool {
		return x[i].Name < x[j].Name
	})
}

func sortModulecallsBySource(x []*ModuleCall) {
	sort.Slice(x, func(i, j int) bool {
		if x[i].Source == x[j].Source {
			return x[i].Name < x[j].Name
		}
		return x[i].Source < x[j].Source
	})
}

func sortModulecallsByPosition(x []*ModuleCall) {
	sort.Slice(x, func(i, j int) bool {
		return x[i].Position.Filename < x[j].Position.Filename || x[i].Position.Line < x[j].Position.Line
	})
}

type modulecalls []*ModuleCall

// WHY: Module calls support sort-by-name, sort-by-source (type), and sort-by-position.
// Source-based sorting groups calls to the same registry module together, which is useful when
// a module instantiates the same child module multiple times with different configurations.
func (mm modulecalls) sort(enabled bool, by string) {
	if !enabled {
		sortModulecallsByPosition(mm)
	} else {
		switch by {
		case print.SortName, print.SortRequired:
			sortModulecallsByName(mm)
		case print.SortType:
			sortModulecallsBySource(mm)
		default:
			sortModulecallsByPosition(mm)
		}
	}
}
