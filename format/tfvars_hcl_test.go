// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package format

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwood-labs/taco-docs/internal/testutil"
	"github.com/northwood-labs/taco-docs/print"
)

// WHY: Golden-file test ensuring tfvars HCL output matches expected fixtures. This format generates
// .tfvars files for use with terraform plan -var-file, so correctness is critical for workflows.
func TestTfvarsHcl(t *testing.T) {
	tests := map[string]struct {
		config print.Config
	}{
		// Base
		"Base": {
			config: testutil.WithSections(),
		},
		"Empty": {
			config: testutil.WithDefaultSections(
				testutil.With(func(c *print.Config) {
					c.ModuleRoot = "empty"
				}),
			),
		},

		// Settings
		"EscapeCharacters": {
			config: testutil.With(func(c *print.Config) {
				c.Settings.Escape = true
			}),
		},
		"PrintDescription": {
			config: testutil.WithSections(
				testutil.With(func(c *print.Config) {
					c.Settings.Description = true
				}),
			),
		},
		"SortByName": {
			config: testutil.WithSections(
				testutil.With(func(c *print.Config) {
					c.Sort.Enabled = true
					c.Sort.By = print.SortName
				}),
			),
		},
		"SortByRequired": {
			config: testutil.WithSections(
				testutil.With(func(c *print.Config) {
					c.Sort.Enabled = true
					c.Sort.By = print.SortRequired
				}),
			),
		},
		"SortByType": {
			config: testutil.WithSections(
				testutil.With(func(c *print.Config) {
					c.Sort.Enabled = true
					c.Sort.By = print.SortType
				}),
			),
		},

		// No section
		"NoInputs": {
			config: testutil.WithSections(
				testutil.With(func(c *print.Config) {
					c.Sections.Inputs = false
				}),
			),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			expected, err := testutil.GetExpected("tfvars", "hcl-"+name)
			assert.Nil(err)

			module, err := testutil.GetModule(&tt.config)
			assert.Nil(err)

			formatter := NewTfvarsHCL(&tt.config)

			err = formatter.Generate(module)
			assert.Nil(err)

			assert.Equal(expected, formatter.Content())
		})
	}
}
