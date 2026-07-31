// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

// Package types provides typed wrappers for Terraform variable default
// values, enabling format-specific serialization to JSON, YAML, XML,
// and TOML.
package types // lint:allow_bad_package_name

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"maps"
	"reflect"
	"slices"
)

// yamlNull is a typed nil used by MarshalYAML methods to represent YAML null
// values without triggering the nilnil linter on bare nil returns.
var yamlNull any

type (
	// Value is the interface for all Terraform variable default values.
	// Terraform supports a rich type system (string, number, bool, list, map,
	// null) but terraform-docs needs to serialize these values across multiple
	// output formats (JSON, YAML, XML, TOML, Markdown). By wrapping each
	// Terraform type in a concrete type that implements Value (plus custom
	// marshalers), we get format-specific rendering without polluting the core
	// domain model with serialization logic.
	Value interface {
		// HasDefault reports whether the variable has a default value.
		HasDefault() bool

		// Length returns the number of elements in the value.
		Length() int

		// Raw returns the underlying Go value.
		Raw() any
	}

	// Nil represents a variable with no default value. It marshals to `null` in
	// JSON and YAML, and uses xsi:nil="true" in XML. The distinction between
	// Nil and Empty is critical: Nil means "the user must provide a value"
	// (required input), while Empty means "the default value is explicitly an
	// empty string.".
	Nil struct{}

	// String represents a non-empty string value. When the underlying string is
	// empty, it marshals to `null` in JSON/YAML (not `""`) because in this
	// context an empty String means "description/version not specified" rather
	// than "explicitly empty." For explicitly-empty defaults, the Empty type is
	// used instead.
	String string

	// Empty represents a Terraform variable whose default is explicitly set to
	// an empty string (default = ""). This is semantically different from Nil
	// (no default) and from String with empty content (no value specified).
	// Empty marshals to `""` in JSON, preserving the user's explicit intent.
	Empty string

	// Number represents a Terraform number value (integer or float). All
	// numeric types are unified under float64 because Terraform's type system
	// doesn't distinguish between integers and floats.
	Number float64

	// Bool represents a Terraform bool value.
	Bool bool

	// List represents a Terraform list/tuple default value. It exists as a
	// distinct type (rather than using []interface{} directly) so that custom
	// XML marshaling can wrap items in <item> tags for well-formed structure.
	List []any

	// Map represents a Terraform map/object default value. Like List, it exists
	// as a distinct type to provide custom XML marshaling where map keys become
	// element names and values become element content.
	Map map[string]any

	xmllistentry struct {
		Value   any      `xml:",chardata"` // lint:allow_format
		XMLName xml.Name `xml:"item"`
	}

	xmlmapentry struct {
		Value   any      `xml:",chardata"` // lint:allow_format
		XMLName xml.Name // lint:allow_format
	}
)

// ValueOf wraps a raw Go interface{} (as parsed from HCL) into the appropriate
// typed Value. This type dispatch is necessary because terraform-config-inspect
// returns default values as interface{}, but we need concrete types to attach
// custom JSON/XML/YAML marshalers that produce the correct output
// representation (e.g., `null` for nil, `""` for explicit empty string).
func ValueOf(v any) Value {
	if v == nil {
		return new(Nil)
	}

	value := reflect.ValueOf(v)

	// We don't really care about all the other kinds.
	switch value.Kind() {
	case reflect.String:
		// Distinguish between "no value" (empty string from zero-value) and
		// "explicitly set to empty string" — these serialize differently.
		if value.IsZero() {
			return Empty("")
		}

		return String(value.String())
	case reflect.Float32, reflect.Float64:
		return Number(value.Float())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Terraform's number type is unified — we normalize all integer types
		// to float64 so that marshaling is consistent regardless of how the HCL
		// parser decoded the value.
		return Number(float64(value.Int()))
	case reflect.Bool:
		return Bool(value.Bool())
	case reflect.Slice:
		sl, ok := value.Interface().([]any)
		if !ok {
			return new(Nil)
		}

		return List(sl)
	case reflect.Map:
		m, ok := value.Interface().(map[string]any)
		if !ok {
			return new(Nil)
		}

		return Map(m)
	default:
		return new(Nil)
	}
}

// TypeOf determines the Terraform type label for a variable. It prefers the
// explicitly declared type string from the .tf file (provided by
// terraform-inspect), falling back to runtime type inference from the default
// value. The fallback handles cases where type is omitted but a default is set
// — Terraform infers the type from the default value in this scenario.
func TypeOf(t string, v any) String {
	if t != "" {
		return String(t)
	}

	if v != nil {
		// We don't really care about all the other kinds.
		switch reflect.ValueOf(v).Kind() {
		case reflect.String:
			return String("string")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			return String("number")
		case reflect.Bool:
			return String("bool")
		case reflect.Slice:
			return String("list")
		case reflect.Map:
			return String("map")
		default:
			return String("any")
		}
	}

	return String("any")
}

// HasDefault returns false for Nil because a nil default means the variable is
// required — the user must supply a value at plan time.
func (Nil) HasDefault() bool {
	return false
}

// Length returns the length of underlying item.
func (Nil) Length() int {
	return 0
}

// Raw underlying value of this type.
func (Nil) Raw() any {
	return nil
}

// MarshalJSON produces literal `null` to match Terraform's JSON representation.
func (Nil) MarshalJSON() ([]byte, error) { // lint:allow_param
	return []byte(`null`), nil
}

// MarshalXML uses the xsi:nil attribute to represent null values in XML,
// following the XML Schema Instance convention for absent values.
func (Nil) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})

	if err := e.EncodeElement(``, start); err != nil {
		return fmt.Errorf("encoding nil XML element: %w", err)
	}

	return nil
}

// MarshalYAML produces a YAML null value.
func (Nil) MarshalYAML() (any, error) {
	return yamlNull, nil
}

// HasDefault indicates a Terraform variable has a default value set.
func (String) HasDefault() bool {
	return true
}

// Length returns the length of underlying item.
func (s String) Length() int {
	return len(s.underlying())
}

// Raw underlying value of this type.
func (s String) Raw() any {
	return s.underlying()
}

// MarshalJSON produces `null` for empty strings (which represents "no value
// specified" in fields like description) or the properly escaped JSON string.
// SetEscapeHTML(false) prevents URLs and HTML in descriptions from being
// mangled.
func (s String) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if string(s) == "" {
		buf.WriteString(`null`)
	} else {
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)

		if err := encoder.Encode(string(s)); err != nil {
			return nil, fmt.Errorf("encoding string as JSON: %w", err)
		}

		buf.Truncate(buf.Len() - 1) // The json encoder adds a newline, this is not configurable.
	}

	return buf.Bytes(), nil
}

// MarshalXML uses xsi:nil for empty strings to signal "no value" in XML output,
// consistent with the JSON marshaling behavior.
func (s String) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if string(s) == "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})

		if err := e.EncodeElement(``, start); err != nil {
			return fmt.Errorf("encoding empty string XML element: %w", err)
		}

		return nil
	}

	if err := e.EncodeElement(string(s), start); err != nil {
		return fmt.Errorf("encoding string XML element: %w", err)
	}

	return nil
}

// MarshalYAML produces null for empty strings, matching the convention that "no
// value specified" renders as null across all output formats.
func (s String) MarshalYAML() (any, error) { // lint:allow_param
	if string(s) == "" || string(s) == `""` {
		return yamlNull, nil
	}

	return string(s), nil
}

func (s String) underlying() string {
	return string(s)
}

// HasDefault indicates a Terraform variable has a default value set.
func (Empty) HasDefault() bool {
	return true
}

// Length returns the length of underlying item.
func (e Empty) Length() int {
	return len(e.underlying())
}

// Raw underlying value of this type.
func (e Empty) Raw() any {
	return e.underlying()
}

// MarshalJSON produces `""` (not `null`) because the user explicitly set the
// default to an empty string — we must preserve that distinction.
func (Empty) MarshalJSON() ([]byte, error) { // lint:allow_param
	return []byte(`""`), nil
}

func (e Empty) underlying() string {
	return string(e)
}

// HasDefault indicates a Terraform variable has a default value set.
func (Number) HasDefault() bool {
	return true
}

// Length returns the length of underlying item.
func (Number) Length() int {
	return 0
}

// Raw underlying value of this type.
func (n Number) Raw() any {
	return n.underlying()
}

func (n Number) underlying() float64 {
	return float64(n)
}

// HasDefault indicates a Terraform variable has a default value set.
func (Bool) HasDefault() bool {
	return true
}

// Length returns the length of underlying item.
func (Bool) Length() int {
	return 0
}

// Raw underlying value of this type.
func (b Bool) Raw() any {
	return b.underlying()
}

func (b Bool) underlying() bool {
	return bool(b)
}

// Underlying returns a defensive copy of the list elements.
func (l List) Underlying() []any {
	var r []any
	for _, i := range l {
		r = append(r, i)
	}

	return r
}

// HasDefault indicates a Terraform variable has a default value set.
func (List) HasDefault() bool {
	return true
}

// Length returns the length of underlying item.
func (l List) Length() int {
	return len(l)
}

// Raw underlying value of this type.
func (l List) Raw() any {
	return l.Underlying()
}

// MarshalXML wraps each list element in an <item> tag. This is necessary
// because XML has no native list syntax — without wrapper elements, the
// structure would be ambiguous when parsed back.
func (l List) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(l) == 0 {
		if err := e.EncodeElement(``, start); err != nil {
			return fmt.Errorf("encoding empty list XML element: %w", err)
		}

		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return fmt.Errorf("encoding list start token: %w", err)
	}

	for _, i := range l {
		if encErr := e.Encode(xmllistentry{XMLName: xml.Name{Local: "item"}, Value: i}); encErr != nil {
			return fmt.Errorf("encoding list item: %w", encErr)
		}
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return fmt.Errorf("encoding list end token: %w", err)
	}

	return nil
}

// Underlying returns a defensive copy of the map.
func (m Map) Underlying() map[string]any {
	r := make(map[string]any)
	maps.Copy(r, m)

	return r
}

// Raw underlying value of this type.
func (m Map) Raw() any {
	return m.Underlying()
}

// HasDefault indicates a Terraform variable has a default value set.
func (Map) HasDefault() bool {
	return true
}

// Length returns the length of underlying item.
func (m Map) Length() int {
	return len(m)
}

// MarshalXML converts a map to XML where each key becomes an element name.
// Nested maps and lists are handled recursively to produce well-formed XML at
// any depth. Keys are sorted for deterministic output.
func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error { // lint:allow_complexity
	if len(m) == 0 {
		if err := e.EncodeElement(``, start); err != nil {
			return fmt.Errorf("encoding empty map XML element: %w", err)
		}

		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return fmt.Errorf("encoding map start token: %w", err)
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	for _, k := range keys {
		// We don't really care about all the other kinds.
		switch reflect.TypeOf(m[k]).Kind() {
		case reflect.Map:
			is := xml.StartElement{Name: xml.Name{Local: k}}
			mv, ok := m[k].(map[string]any)

			if !ok {
				break
			}

			if encErr := Map(mv).MarshalXML(e, is); encErr != nil {
				return fmt.Errorf("encoding nested map %q: %w", k, encErr)
			}
		case reflect.Slice:
			is := xml.StartElement{Name: xml.Name{Local: k}}
			sv, ok := m[k].([]any)

			if !ok {
				break
			}

			if encErr := List(sv).MarshalXML(e, is); encErr != nil {
				return fmt.Errorf("encoding nested list %q: %w", k, encErr)
			}
		default:
			if encErr := e.Encode(xmlmapentry{XMLName: xml.Name{Local: k}, Value: m[k]}); encErr != nil {
				return fmt.Errorf("encoding map entry %q: %w", k, encErr)
			}
		}
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return fmt.Errorf("encoding map end token: %w", err)
	}

	return nil
}
