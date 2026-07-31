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
	"os"
	"path/filepath"
	"testing"

	assertpkg "github.com/go-openapi/testify/assert"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// Tests the generator's template execution logic, including custom templates
// with section references, file includes, and unknown section handling. Ensures
// the --content flag works correctly.
func TestExecuteTemplate(t *testing.T) {
	header := "this is the header"
	footer := "this is the footer"

	tests := map[string]struct {
		content  string
		template string
		expected string
		complex  bool
		wantErr  bool
	}{
		"Compatible without template": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "",
			expected: "this is the header\nthis is the footer",
			wantErr:  false,
		},
		"Compatible with template not empty section": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "{{ .Header }}",
			expected: "this is the header",
			wantErr:  false,
		},
		"Compatible with template empty section": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "{{ .Inputs }}",
			expected: "",
			wantErr:  false,
		},
		"Compatible with template and unknown section": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "{{ .Unknown }}",
			expected: "",
			wantErr:  true,
		},
		"Compatible with template include file": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "{{ include \"testdata/generator/sample-file.txt\" }}",
			expected: "Sample file to be included.",
			wantErr:  false,
		},
		"Compatible with template include unknown file": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "{{ include \"file-not-found\" }}",
			expected: "",
			wantErr:  true,
		},
		"Compatible with template optional include file": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "{{ include_optional \"testdata/generator/sample-file.txt\" \"\" }}",
			expected: "Sample file to be included.",
			wantErr:  false,
		},
		"Compatible with template optional include unknown file": {
			complex:  true,
			content:  "this is the header\nthis is the footer",
			template: "{{ include_optional \"file-not-found\" \"This file does not exist\" }}",
			expected: "This file does not exist",
			wantErr:  false,
		},
		"Incompatible without template": {
			complex:  false,
			content:  "header: \"this is the header\"\nfooter: \"this is the footer\"",
			template: "",
			expected: "header: \"this is the header\"\nfooter: \"this is the footer\"",
			wantErr:  false,
		},
		"Incompatible with template": {
			complex:  false,
			content:  "header: \"this is the header\"\nfooter: \"this is the footer\"",
			template: "{{ .Header }}",
			expected: "header: \"this is the header\"\nfooter: \"this is the footer\"",
			wantErr:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assertpkg.New(t)

			config := print.DefaultConfig()

			generator := newGenerator(config, tt.complex)

			generator.content = tt.content
			generator.header = header
			generator.footer = footer

			actual, err := generator.Render(tt.template)

			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tt.expected, actual)
			}
		})
	}
}

// Ensures the withX generator functions correctly set their respective fields
// on the generator.
func TestGeneratorFunc(t *testing.T) {
	text := "foo"

	tests := map[string]struct {
		fn     func(string) generateFunc
		actual func(*generator) string
	}{
		"withContent": {
			fn:     withContent,
			actual: func(r *generator) string { return r.content },
		},
		"withHeader": {
			fn:     withHeader,
			actual: func(r *generator) string { return r.header },
		},
		"withFooter": {
			fn:     withFooter,
			actual: func(r *generator) string { return r.footer },
		},
		"withInputs": {
			fn:     withInputs,
			actual: func(r *generator) string { return r.inputs },
		},
		"withModules": {
			fn:     withModules,
			actual: func(r *generator) string { return r.modules },
		},
		"withOutputs": {
			fn:     withOutputs,
			actual: func(r *generator) string { return r.outputs },
		},
		"withProviders": {
			fn:     withProviders,
			actual: func(r *generator) string { return r.providers },
		},
		"withRequirements": {
			fn:     withRequirements,
			actual: func(r *generator) string { return r.requirements },
		},
		"withResources": {
			fn:     withResources,
			actual: func(r *generator) string { return r.resources },
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assertpkg.New(t)

			config := print.DefaultConfig()

			config.Sections.Footer = true

			generator := newGenerator(config, false, tt.fn(text))

			assert.Equal(text, tt.actual(generator))
		})
	}
}

// Tests withModule generator function to verify module data is correctly loaded
// from disk.
func TestGeneratorFuncModule(t *testing.T) {
	t.Run("withModule", func(t *testing.T) {
		assert := assertpkg.New(t)

		config := print.DefaultConfig()

		config.ModuleRoot = filepath.Join("..", "terraform", "testdata", "full-example")

		module, err := terraform.LoadWithOptions(config)

		assert.NoError(err)

		generator := newGenerator(config, true, withModule(module))

		path := filepath.Join("..", "terraform", "testdata", "expected", "full-example-mainTf-Header.golden")
		data, err := os.ReadFile(path) // lint:allow_dynamic_filename

		assert.NoError(err)

		expected := string(data)

		assert.Equal(expected, generator.module.Header)
		assert.Empty(generator.module.Footer)
		assert.Len(generator.module.Inputs, 7)
		assert.Len(generator.module.Outputs, 4)
	})
}

// Validates forEach correctly populates all named sections. This is the
// callback mechanism used by formatters to generate each section independently.
func TestForEach(t *testing.T) {
	config := print.DefaultConfig()

	generator := newGenerator(config, false)

	err := generator.forEach(func(name string) (string, error) {
		return name, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		actual string
	}{
		"all":          {actual: generator.content},
		"header":       {actual: generator.header},
		"footer":       {actual: generator.footer},
		"inputs":       {actual: generator.inputs},
		"modules":      {actual: generator.modules},
		"outputs":      {actual: generator.outputs},
		"providers":    {actual: generator.providers},
		"requirements": {actual: generator.requirements},
		"resources":    {actual: generator.resources},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assertpkg.New(t)
			assert.Equal(name, tt.actual)
		})
	}
}
