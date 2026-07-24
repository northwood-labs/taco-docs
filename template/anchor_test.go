/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// WHY: Validates Markdown anchor generation with/without escape and anchor flags. Broken anchors
// mean in-page links don't work in rendered docs, breaking navigation for large modules.
func TestAnchorMarkdown(t *testing.T) {
	tests := []struct {
		typeSection string
		name        string
		anchor      bool
		escape      bool
		expected    string
	}{
		{
			typeSection: "module",
			name:        "banana_anchor_escape",
			anchor:      true,
			escape:      true,
			expected:    "<a name=\"module_banana_anchor_escape\"></a> [banana\\_anchor\\_escape](#module\\_banana\\_anchor\\_escape)",
		},
		{
			typeSection: "module",
			name:        "banana_anchor_noescape",
			anchor:      true,
			escape:      false,
			expected:    "<a name=\"module_banana_anchor_noescape\"></a> [banana_anchor_noescape](#module_banana_anchor_noescape)",
		},
		{
			typeSection: "module",
			name:        "banana_anchor_escape",
			anchor:      false,
			escape:      true,
			expected:    "banana\\_anchor\\_escape",
		},
		{
			typeSection: "module",
			name:        "banana_anchor_noescape",
			anchor:      false,
			escape:      false,
			expected:    "banana_anchor_noescape",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			actual := CreateAnchorMarkdown(tt.typeSection, tt.name, tt.anchor, tt.escape)

			assert.Equal(tt.expected, actual)
		})
	}
}

// WHY: Same as Markdown anchors but for AsciiDoc's [[id]] <<id,label>> syntax.
func TestAnchorAsciidoc(t *testing.T) {
	tests := []struct {
		typeSection string
		name        string
		anchor      bool
		escape      bool
		expected    string
	}{
		{
			typeSection: "module",
			name:        "banana_anchor_escape",
			anchor:      true,
			escape:      true,
			expected:    "[[module\\_banana\\_anchor\\_escape]] <<module\\_banana\\_anchor\\_escape,banana\\_anchor\\_escape>>",
		},
		{
			typeSection: "module",
			name:        "banana_anchor_noescape",
			anchor:      true,
			escape:      false,
			expected:    "[[module_banana_anchor_noescape]] <<module_banana_anchor_noescape,banana_anchor_noescape>>",
		},
		{
			typeSection: "module",
			name:        "banana_anchor_escape",
			anchor:      false,
			escape:      true,
			expected:    "banana\\_anchor\\_escape",
		},
		{
			typeSection: "module",
			name:        "banana_anchor_noescape",
			anchor:      false,
			escape:      false,
			expected:    "banana_anchor_noescape",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			actual := CreateAnchorAsciidoc(tt.typeSection, tt.name, tt.anchor, tt.escape)

			assert.Equal(tt.expected, actual)
		})
	}
}
