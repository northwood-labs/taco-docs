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
	"encoding/xml"
	"fmt"
	"slices"
	"strings"

	"github.com/northwood-labs/taco-docs/internal/types"
)

type (
	// Output represents a Terraform output.
	//
	// Custom MarshalJSON/XML/YAML methods exist because the --output-values
	// feature conditionally includes or excludes the Value and Sensitive
	// fields. When output values are disabled (default), these fields are
	// omitted so generated docs don't contain misleading empty entries. When
	// enabled, the "withvalue" shadow struct forces serialization of even
	// zero-valued fields (empty string, false) so users see the actual state.
	// The ShowValue flag drives this switch at marshal time.
	Output struct {
		Value       types.Value  `json:"value,omitempty"     toml:"value,omitempty"     xml:"value,omitempty"     yaml:"value,omitempty"`     // lint:ignore_length
		Name        string       `json:"name"                toml:"name"                xml:"name"                yaml:"name"`                // lint:ignore_length
		Description types.String `json:"description"         toml:"description"         xml:"description"         yaml:"description"`         // lint:ignore_length
		Position    Position     `json:"-"                   toml:"-"                   xml:"-"                   yaml:"-"`                   // lint:ignore_length
		Sensitive   bool         `json:"sensitive,omitempty" toml:"sensitive,omitempty" xml:"sensitive,omitempty" yaml:"sensitive,omitempty"` // lint:ignore_length
		ShowValue   bool         `json:"-"                   toml:"-"                   xml:"-"                   yaml:"-"`                   // lint:ignore_length
	}

	// withvalue is a shadow struct identical to Output but without omitempty on
	// Value/Sensitive. Go's encoding packages check struct tags at marshal
	// time, so we need a separate type to force serialization of zero-valued
	// fields when --output-values is active.
	withvalue struct {
		Value       types.Value  `json:"value"       toml:"value"       xml:"value"       yaml:"value"`
		Name        string       `json:"name"        toml:"name"        xml:"name"        yaml:"name"`
		Description types.String `json:"description" toml:"description" xml:"description" yaml:"description"`
		Position    Position     `json:"-"           toml:"-"           xml:"-"           yaml:"-"`
		Sensitive   bool         `json:"sensitive"   toml:"sensitive"   xml:"sensitive"   yaml:"sensitive"`
		ShowValue   bool         `json:"-"           toml:"-"           xml:"-"           yaml:"-"`
	}

	// output is used for unmarshalling `terraform outputs --json` into.
	output struct {
		Type      any  `json:"type"`
		Value     any  `json:"value"`
		Sensitive bool `json:"sensitive"`
	}

	outputs []*Output
)

// GetValue returns JSON representation of the 'Value', which is an 'interface'.
// If 'Value' is a primitive type, the primitive value of 'Value' will be
// returned and not the JSON formatted of it.
func (o *Output) GetValue() string {
	if !o.ShowValue || o.Value == nil {
		return ""
	}

	marshaled, err := json.MarshalIndent(o.Value, "", "  ")
	if err != nil {
		panic(err)
	}

	value := string(marshaled)
	if value == `null` {
		return "" // types.Nil.
	}

	return value // everything else.
}

// HasDefault indicates if a Terraform output has a default value set.
func (o *Output) HasDefault() bool {
	if !o.ShowValue || o.Value == nil {
		return false
	}

	return o.Value.HasDefault()
}

// MarshalJSON custom yaml marshal function to take '--output-values' flag into
// consideration. It means if the flag is not set Value and Sensitive fields are
// set to 'omitempty', otherwise if output values are being shown 'omitempty'
// gets explicitly removed to show even empty and false values.
func (o *Output) MarshalJSON() ([]byte, error) {
	fn := func(oo any) ([]byte, error) { // lint:allow_param
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		if err := enc.Encode(oo); err != nil {
			panic(err)
		}

		return buf.Bytes(), nil
	}
	if o.ShowValue {
		result, err := fn(withvalue(*o))
		if err != nil {
			return nil, fmt.Errorf("encoding output as JSON: %w", err)
		}

		return result, nil
	}

	o.Value = nil       // explicitly make empty.
	o.Sensitive = false // explicitly make empty.

	result, err := fn(*o)
	if err != nil {
		return nil, fmt.Errorf("encoding output as JSON: %w", err)
	}

	return result, nil
}

// MarshalXML custom xml marshal function to take '--output-values' flag into
// consideration. It means if the flag is not set Value and Sensitive fields are
// set to 'omitempty', otherwise if output values are being shown 'omitempty'
// gets explicitly removed to show even empty and false values.
func (o *Output) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	fn := func(v any, name string) error {
		return e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: name}})
	}

	err := e.EncodeToken(start)
	if err != nil {
		return fmt.Errorf("encoding XML start token: %w", err)
	}

	if err := fn(o.Name, "name"); err != nil {
		return fmt.Errorf("encoding output name as XML: %w", err)
	}

	if err := fn(o.Description, "description"); err != nil {
		return fmt.Errorf("encoding output description as XML: %w", err)
	}

	if o.ShowValue {
		if err := fn(o.Value, "value"); err != nil {
			return fmt.Errorf("encoding output value as XML: %w", err)
		}

		if err := fn(o.Sensitive, "sensitive"); err != nil {
			return fmt.Errorf("encoding output sensitive flag as XML: %w", err)
		}
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return fmt.Errorf("encoding XML end token: %w", err)
	}

	return nil
}

// MarshalYAML custom yaml marshal function to take '--output-values' flag into
// consideration. It means if the flag is not set Value and Sensitive fields are
// set to 'omitempty', otherwise if output values are being shown 'omitempty'
// gets explicitly removed to show even empty and false values.
func (o *Output) MarshalYAML() (any, error) { // lint:allow_param
	if o.ShowValue {
		return withvalue(*o), nil
	}

	o.Value = nil       // explicitly make empty.
	o.Sensitive = false // explicitly make empty.

	return *o, nil
}

func sortOutputsByName(x []*Output) {
	slices.SortFunc(x, func(a, b *Output) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func sortOutputsByPosition(x []*Output) {
	slices.SortFunc(x, func(a, b *Output) int {
		if a.Position.Filename == b.Position.Filename {
			return cmp.Compare(a.Position.Line, b.Position.Line)
		}

		return strings.Compare(a.Position.Filename, b.Position.Filename)
	})
}

// Outputs support only name-based and position-based sorting. Unlike inputs,
// outputs have no "required" or "type" dimension, so fewer strategies are
// needed.
func (oo outputs) sort(enabled bool, _ string) { // lint:allow_param lint:allow_control_coupling_antipattern
	if !enabled {
		sortOutputsByPosition(oo)
	} else {
		// always sort by name if sorting is enabled.
		sortOutputsByName(oo)
	}
}
