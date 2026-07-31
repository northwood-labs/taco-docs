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
	"bytes"
	"cmp"
	"encoding/json"
	"slices"
	"strings"

	"github.com/northwood-labs/taco-docs/internal/types"
	"github.com/northwood-labs/taco-docs/print"
)

type (
	// Input represents a Terraform input.
	//
	// Multiple sort strategies (name, required, type, position) exist because
	// different audiences organize documentation differently—alphabetical for
	// quick lookup, by-required for onboarding new users who need to know what
	// they must supply, by-type for grouping related variables, and by-position
	// for matching the author's original source order.
	Input struct {
		Default     types.Value  `json:"default"     toml:"default"     xml:"default"     yaml:"default"`
		Name        string       `json:"name"        toml:"name"        xml:"name"        yaml:"name"`
		Type        types.String `json:"type"        toml:"type"        xml:"type"        yaml:"type"`
		Description types.String `json:"description" toml:"description" xml:"description" yaml:"description"`
		Position    Position     `json:"-"           toml:"-"           xml:"-"           yaml:"-"`
		Required    bool         `json:"required"    toml:"required"    xml:"required"    yaml:"required"`
		Sensitive   bool         `json:"sensitive"   toml:"sensitive"   xml:"sensitive"   yaml:"sensitive"`
	}

	inputs []*Input
)

// GetValue returns JSON representation of the 'Default' value, which is an
// 'interface'. If 'Default' is a primitive type, the primitive value of
// 'Default' will be returned and not the JSON formatted of it.
func (i *Input) GetValue() string {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(i.Default)
	if err != nil {
		panic(err)
	}

	value := strings.TrimSpace(buf.String())
	if value == `null` {
		if i.Required {
			return ""
		}

		return `null` // explicit 'null' value.
	}

	return value // everything else.
}

// HasDefault indicates if a Terraform variable has a default value set.
func (i *Input) HasDefault() bool {
	return i.Default.HasDefault() || !i.Required
}

func sortInputsByName(x []*Input) {
	slices.SortFunc(x, func(a, b *Input) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func sortInputsByRequired(x []*Input) {
	slices.SortFunc(x, func(a, b *Input) int {
		if a.HasDefault() == b.HasDefault() {
			return strings.Compare(a.Name, b.Name)
		}

		if !a.HasDefault() && b.HasDefault() {
			return -1
		}

		return 1
	})
}

func sortInputsByPosition(x []*Input) {
	slices.SortFunc(x, func(a, b *Input) int {
		if a.Position.Filename == b.Position.Filename {
			return cmp.Compare(a.Position.Line, b.Position.Line)
		}

		return strings.Compare(a.Position.Filename, b.Position.Filename)
	})
}

func sortInputsByType(x []*Input) {
	slices.SortFunc(x, func(a, b *Input) int {
		if a.Type == b.Type {
			return strings.Compare(a.Name, b.Name)
		}

		return strings.Compare(string(a.Type), string(b.Type))
	})
}

func (ii inputs) sort(enabled bool, by string) { // lint:allow_param lint:allow_control_coupling_antipattern
	if !enabled {
		sortInputsByPosition(ii)
	} else {
		switch by {
		case print.SortType:
			sortInputsByType(ii)
		case print.SortRequired:
			sortInputsByRequired(ii)
		case print.SortName:
			sortInputsByName(ii)
		default:
			sortInputsByPosition(ii)
		}
	}
}
