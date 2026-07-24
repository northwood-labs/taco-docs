// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// WHY: Validates that empty strings are detected and rendered as "" (quoted empty) in docs rather than null.
func TestEmpty(t *testing.T) {
	values := List{""}
	testPrimitive(t, []testprimitive{
		{
			name:   "value empty and type string",
			values: values,
			types:  "string",
			expected: expected{
				typeName:   "string",
				valueKind:  "types.Empty",
				hasDefault: true,
			},
		},
		{
			name:   "value empty and type empty",
			values: values,
			types:  "",
			expected: expected{
				typeName:   "string",
				valueKind:  "types.Empty",
				hasDefault: true,
			},
		},
	})
}

// WHY: Empty strings are scalar—Length must be 0.
func TestEmptyLength(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected int
	}{
		{
			name:     "empty length",
			value:    "",
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tt.expected, Empty(tt.value).Length())
		})
	}
}

// WHY: Confirms the underlying Go value for Empty type.
func TestEmptyUnderlying(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty underlying",
			value: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tt.value, Empty(tt.value).underlying())
		})
	}
}

// WHY: Ensures Empty marshals to "" in JSON output, not null or the underlying value.
func TestEmptyMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "string marshal JSON",
			value:    "foo",
			expected: "\"\"",
		},
		{
			name:     "string marshal JSON",
			value:    "",
			expected: "\"\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual, err := Empty(tt.value).MarshalJSON()

			assert.Nil(err)
			assert.Equal(tt.expected, string(actual))
		})
	}
}
