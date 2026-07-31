// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package types // lint:allow_bad_package_name

import (
	"reflect"
	"testing"

	assertpkg "github.com/go-openapi/testify/assert"
)

type (
	expected struct {
		typeName   string
		valueKind  string
		hasDefault bool
	}

	testprimitive struct {
		name     string
		values   List
		types    string
		expected expected
	}

	testlist struct {
		name     string
		values   []List
		types    string
		expected expected
	}

	testmap struct {
		name     string
		values   []Map
		types    string
		expected expected
	}
)

// Verifies the type system's core ValueOf/TypeOf dispatch logic for primitive
// values. This is the foundation for rendering input defaults—incorrect type
// detection cascades into wrong output everywhere.
func testPrimitive(t *testing.T, tests []testprimitive) {
	t.Helper()

	for i := range tests {
		tt := tests[i]
		for _, tv := range tt.values {
			t.Run(tt.name, func(t *testing.T) {
				assert := assertpkg.New(t)

				actualValue := ValueOf(tv)
				actualType := TypeOf(tt.types, tv)

				assert.Equal(tt.expected.typeName, string(actualType))
				assert.Equal(tt.expected.valueKind, reflect.TypeOf(actualValue).String())
				assert.Equal(tt.expected.hasDefault, actualValue.HasDefault())
			})
		}
	}
}
