/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package types

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"reflect"
	"sort"
)

// Value is the interface for all Terraform variable default values. Terraform
// supports a rich type system (string, number, bool, list, map, null) but
// terraform-docs needs to serialize these values across multiple output formats
// (JSON, YAML, XML, TOML, Markdown). By wrapping each Terraform type in a
// concrete type that implements Value (plus custom marshalers), we get
// format-specific rendering without polluting the core domain model with
// serialization logic.
type Value interface {
	HasDefault() bool
	Length() int
	Raw() interface{}
}

// ValueOf wraps a raw Go interface{} (as parsed from HCL) into the appropriate
// typed Value. This type dispatch is necessary because terraform-config-inspect
// returns default values as interface{}, but we need concrete types to attach
// custom JSON/XML/YAML marshalers that produce the correct output representation
// (e.g., `null` for nil, `""` for explicit empty string).
func ValueOf(v interface{}) Value {
	if v == nil {
		return new(Nil)
	}
	value := reflect.ValueOf(v)

	// We don't really care about all the other kinds.
	//
	//nolint:exhaustive
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
		// to float64 so that marshaling is consistent regardless of how the
		// HCL parser decoded the value.
		return Number(float64(value.Int()))
	case reflect.Bool:
		return Bool(value.Bool())
	case reflect.Slice:
		return List(value.Interface().([]interface{}))
	case reflect.Map:
		return Map(value.Interface().(map[string]interface{}))
	}
	return new(Nil)
}

// TypeOf determines the Terraform type label for a variable. It prefers the
// explicitly declared type string from the .tf file (provided by terraform-inspect),
// falling back to runtime type inference from the default value. The fallback
// handles cases where type is omitted but a default is set — Terraform infers the
// type from the default value in this scenario.
func TypeOf(t string, v interface{}) String {
	if t != "" {
		return String(t)
	}
	if v != nil {
		// We don't really care about all the other kinds.
		//
		//nolint:exhaustive
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
		}
	}
	return String("any")
}

// Nil represents a variable with no default value. It marshals to `null` in JSON
// and YAML, and uses xsi:nil="true" in XML. The distinction between Nil and Empty
// is critical: Nil means "the user must provide a value" (required input), while
// Empty means "the default value is explicitly an empty string."
type Nil struct{}

// HasDefault returns false for Nil because a nil default means the variable is
// required — the user must supply a value at plan time.
func (n Nil) HasDefault() bool {
	return false
}

// Length returns the length of underlying item
func (n Nil) Length() int {
	return 0
}

// Raw underlying value of this type.
func (n Nil) Raw() interface{} {
	return nil
}

// MarshalJSON produces literal `null` to match Terraform's JSON representation.
func (n Nil) MarshalJSON() ([]byte, error) {
	return []byte(`null`), nil
}

// MarshalXML uses the xsi:nil attribute to represent null values in XML, following
// the XML Schema Instance convention for absent values.
func (n Nil) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
	return e.EncodeElement(``, start)
}

// MarshalYAML produces a YAML null value.
func (n Nil) MarshalYAML() (interface{}, error) {
	return nil, nil
}

// String represents a non-empty string value. When the underlying string is empty,
// it marshals to `null` in JSON/YAML (not `""`) because in this context an empty
// String means "description/version not specified" rather than "explicitly empty."
// For explicitly-empty defaults, the Empty type is used instead.
type String string

// nolint
func (s String) underlying() string {
	return string(s)
}

// HasDefault indicates a Terraform variable has a default value set.
func (s String) HasDefault() bool {
	return true
}

// Length returns the length of underlying item
func (s String) Length() int {
	return len(s.underlying())
}

// Raw underlying value of this type.
func (s String) Raw() interface{} {
	return s.underlying()
}

// MarshalJSON produces `null` for empty strings (which represents "no value
// specified" in fields like description) or the properly escaped JSON string.
// SetEscapeHTML(false) prevents URLs and HTML in descriptions from being mangled.
func (s String) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if len(string(s)) == 0 {
		buf.WriteString(`null`)
	} else {
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(string(s)); err != nil {
			return nil, err
		}
		buf.Truncate(buf.Len() - 1) // The json encoder adds a newline, this is not configurable
	}
	return buf.Bytes(), nil
}

// MarshalXML uses xsi:nil for empty strings to signal "no value" in XML output,
// consistent with the JSON marshaling behavior.
func (s String) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if string(s) == "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
		return e.EncodeElement(``, start)
	}
	return e.EncodeElement(string(s), start)
}

// MarshalYAML produces null for empty strings, matching the convention that
// "no value specified" renders as null across all output formats.
func (s String) MarshalYAML() (interface{}, error) {
	if len(string(s)) == 0 || string(s) == `""` {
		return nil, nil
	}
	return string(s), nil
}

// Empty represents a Terraform variable whose default is explicitly set to an
// empty string (default = ""). This is semantically different from Nil (no default)
// and from String with empty content (no value specified). Empty marshals to `""`
// in JSON, preserving the user's explicit intent.
type Empty string

// nolint
func (e Empty) underlying() string {
	return string(e)
}

// HasDefault indicates a Terraform variable has a default value set.
func (e Empty) HasDefault() bool {
	return true
}

// Length returns the length of underlying item
func (e Empty) Length() int {
	return len(e.underlying())
}

// Raw underlying value of this type.
func (e Empty) Raw() interface{} {
	return e.underlying()
}

// MarshalJSON produces `""` (not `null`) because the user explicitly set the
// default to an empty string — we must preserve that distinction.
func (e Empty) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

// Number represents a Terraform number value (integer or float). All numeric
// types are unified under float64 because Terraform's type system doesn't
// distinguish between integers and floats.
type Number float64

// nolint
func (n Number) underlying() float64 {
	return float64(n)
}

// HasDefault indicates a Terraform variable has a default value set.
func (n Number) HasDefault() bool {
	return true
}

// Length returns the length of underlying item
func (n Number) Length() int {
	return 0
}

// Raw underlying value of this type.
func (n Number) Raw() interface{} {
	return n.underlying()
}

// Bool represents a Terraform bool value.
type Bool bool

// nolint
func (b Bool) underlying() bool {
	return bool(b)
}

// HasDefault indicates a Terraform variable has a default value set.
func (b Bool) HasDefault() bool {
	return true
}

// Length returns the length of underlying item
func (b Bool) Length() int {
	return 0
}

// Raw underlying value of this type.
func (b Bool) Raw() interface{} {
	return b.underlying()
}

// List represents a Terraform list/tuple default value. It exists as a distinct
// type (rather than using []interface{} directly) so that custom XML marshaling
// can wrap items in <item> tags for well-formed structure.
type List []interface{}

// Underlying returns a defensive copy of the list elements.
func (l List) Underlying() []interface{} {
	r := make([]interface{}, 0)
	for _, i := range l {
		r = append(r, i)
	}
	return r
}

// HasDefault indicates a Terraform variable has a default value set.
func (l List) HasDefault() bool {
	return true
}

// Length returns the length of underlying item
func (l List) Length() int {
	return len(l)
}

// Raw underlying value of this type.
func (l List) Raw() interface{} {
	return l.Underlying()
}

type xmllistentry struct {
	XMLName xml.Name    `xml:"item"`
	Value   interface{} `xml:",chardata"`
}

// MarshalXML wraps each list element in an <item> tag. This is necessary because
// XML has no native list syntax — without wrapper elements, the structure would be
// ambiguous when parsed back.
func (l List) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(l) == 0 {
		return e.EncodeElement(``, start)
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	for _, i := range l {
		e.Encode(xmllistentry{XMLName: xml.Name{Local: "item"}, Value: i}) //nolint:errcheck,gosec
	}
	return e.EncodeToken(start.End())
}

// Map represents a Terraform map/object default value. Like List, it exists as
// a distinct type to provide custom XML marshaling where map keys become element
// names and values become element content.
type Map map[string]interface{}

// Underlying returns a defensive copy of the map.
func (m Map) Underlying() map[string]interface{} {
	r := make(map[string]interface{})
	for k, e := range m {
		r[k] = e
	}
	return r
}

// Raw underlying value of this type.
func (m Map) Raw() interface{} {
	return m.Underlying()
}

// HasDefault indicates a Terraform variable has a default value set.
func (m Map) HasDefault() bool {
	return true
}

// Length returns the length of underlying item
func (m Map) Length() int {
	return len(m)
}

type xmlmapentry struct {
	XMLName xml.Name    `xml:","`
	Value   interface{} `xml:",chardata"`
}

// sortmapkeys ensures deterministic XML output by sorting map keys alphabetically.
// Without this, Go's random map iteration would produce non-reproducible output,
// making diffs noisy and CI checks unreliable.
type sortmapkeys []string

func (s sortmapkeys) Len() int           { return len(s) }
func (s sortmapkeys) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortmapkeys) Less(i, j int) bool { return s[i] < s[j] }

// MarshalXML converts a map to XML where each key becomes an element name.
// Nested maps and lists are handled recursively to produce well-formed XML
// at any depth. Keys are sorted for deterministic output.
func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return e.EncodeElement(``, start)
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sortmapkeys(keys))
	for _, k := range keys {
		// We don't really care about all the other kinds.
		//
		//nolint:exhaustive
		switch reflect.TypeOf(m[k]).Kind() {
		case reflect.Map:
			is := xml.StartElement{Name: xml.Name{Local: k}}
			Map(m[k].(map[string]interface{})).MarshalXML(e, is) //nolint:errcheck,gosec
		case reflect.Slice:
			is := xml.StartElement{Name: xml.Name{Local: k}}
			List(m[k].([]interface{})).MarshalXML(e, is) //nolint:errcheck,gosec
		default:
			e.Encode(xmlmapentry{XMLName: xml.Name{Local: k}, Value: m[k]}) //nolint:errcheck,gosec
		}
	}
	return e.EncodeToken(start.End())
}
