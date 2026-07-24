/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// WHY: Validates Bool type detection and HasDefault behavior. Booleans must be recognized whether
// the HCL type annotation is present or empty, so generated docs show "true"/"false" correctly.
func TestBool(t *testing.T) {
	values := List{true, false}
	testPrimitive(t, []testprimitive{
		{
			name:   "value not nil and type bool",
			values: values,
			types:  "bool",
			expected: expected{
				typeName:   "bool",
				valueKind:  "types.Bool",
				hasDefault: true,
			},
		},
		{
			name:   "value not nil and type empty",
			values: values,
			types:  "",
			expected: expected{
				typeName:   "bool",
				valueKind:  "types.Bool",
				hasDefault: true,
			},
		},
	})
}

// WHY: Booleans are scalar—Length must always be 0 to distinguish from collection types.
func TestBoolLength(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected int
	}{
		{
			name:     "bool length",
			value:    true,
			expected: 0,
		},
		{
			name:     "bool length",
			value:    false,
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tt.expected, Bool(tt.value).Length())
		})
	}
}

// WHY: Confirms the underlying Go value is accessible for serialization to JSON/XML/YAML.
func TestBoolUnderlying(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{
			name:  "bool underlying",
			value: true,
		},
		{
			name:  "bool underlying",
			value: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tt.value, Bool(tt.value).underlying())
		})
	}
}
