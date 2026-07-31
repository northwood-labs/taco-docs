// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package format // lint:allow_naming_conflict_stdlib lint:no_dupe

import (
	"testing"

	assertpkg "github.com/go-openapi/testify/assert"

	"github.com/northwood-labs/taco-docs/internal/testutil"
	"github.com/northwood-labs/taco-docs/print"
)

// Golden-file test ensuring XML output matches expected fixtures.
func TestXml(t *testing.T) {
	tests := map[string]struct {
		config print.Config
	}{
		// Base.
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
		"HideAll": {
			config: testutil.With(func(c *print.Config) {
				c.Sections.Header = false // Since we don't show the header, the file won't be loaded at all.
				c.HeaderFrom = "bad.tf"
			}),
		},

		// Settings.
		"OutputValues": {
			config: testutil.With(func(c *print.Config) {
				c.Sections.Outputs = true
				c.OutputValues.Enabled = true
				c.OutputValues.From = "output_values.json"
				c.Settings.Sensitive = true
			}),
		},

		// Only section.
		"OnlyDataSources": {
			config: testutil.With(func(c *print.Config) { c.Sections.DataSources = true }),
		},
		"OnlyHeader": {
			config: testutil.With(func(c *print.Config) { c.Sections.Header = true }),
		},
		"OnlyFooter": {
			config: testutil.With(func(c *print.Config) {
				c.Sections.Footer = true
				c.FooterFrom = "footer.md"
			}),
		},
		"OnlyInputs": {
			config: testutil.With(func(c *print.Config) { c.Sections.Inputs = true }),
		},
		"OnlyOutputs": {
			config: testutil.With(func(c *print.Config) { c.Sections.Outputs = true }),
		},
		"OnlyModulecalls": {
			config: testutil.With(func(c *print.Config) { c.Sections.ModuleCalls = true }),
		},
		"OnlyProviders": {
			config: testutil.With(func(c *print.Config) { c.Sections.Providers = true }),
		},
		"OnlyRequirements": {
			config: testutil.With(func(c *print.Config) { c.Sections.Requirements = true }),
		},
		"OnlyResources": {
			config: testutil.With(func(c *print.Config) { c.Sections.Resources = true }),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assertpkg.New(t)

			expected, err := testutil.GetExpected("xml", "xml-"+name)
			assert.NoError(err)

			module, err := testutil.GetModule(&tt.config)
			assert.NoError(err)

			formatter := NewXML(&tt.config)

			err = formatter.Generate(module)
			assert.NoError(err)

			assert.Equal(expected, formatter.Content())
		})
	}
}
