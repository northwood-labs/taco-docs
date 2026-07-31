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
	"embed"
	"fmt"
	gotemplate "text/template"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/template"
	"github.com/northwood-labs/taco-docs/terraform"
)

//go:embed templates/asciidoc_table*.tmpl
var asciidocTableFS embed.FS

// asciidocTable represents AsciiDoc Table format.
//
// Teams using Asciidoctor or Antora for their documentation sites need native
// AsciiDoc output rather than embedded Markdown. This is the compact tabular
// equivalent of markdownTable for the AsciiDoc ecosystem.
type asciidocTable struct {
	*generator

	config   *print.Config
	template *template.Template
}

// Multiple short aliases registered so users can type "adoc" or "asciidoc"
// interchangeably, matching common community shorthand.
func init() { // lint:allow_init
	register(map[string]initializerFn{
		"asciidoc":       asInitializer(NewAsciidocTable),
		"asciidoc table": asInitializer(NewAsciidocTable),
		"asciidoc tbl":   asInitializer(NewAsciidocTable),
		"adoc":           asInitializer(NewAsciidocTable),
		"adoc table":     asInitializer(NewAsciidocTable),
		"adoc tbl":       asInitializer(NewAsciidocTable),
	})
}

// NewAsciidocTable returns new instance of Asciidoc Table.
func NewAsciidocTable(config *print.Config) *asciidocTable { // lint:allow_unexported_return
	items := readTemplateItems(asciidocTableFS, "asciidoc_table")

	// AsciiDoc has its own escaping rules (e.g., | inside table cells).
	// Disabling the generic markdown escape prevents double-escaping that would
	// corrupt the AsciiDoc output.
	config.Settings.Escape = false

	tt := template.New(config, items...)
	tt.CustomFunc(gotemplate.FuncMap{
		"type": func(t string) string {
			inputType, _ := PrintFencedCodeBlock(t, "")
			return inputType
		},
		"value": func(v string) string {
			result := "n/a"
			if v != "" {
				result, _ = PrintFencedCodeBlock(v, "")
			}

			return result
		},
	})

	return &asciidocTable{
		generator: newGenerator(config, true),
		config:    config,
		template:  tt,
	}
}

// Generate a Terraform module as AsciiDoc tables.
func (t *asciidocTable) Generate(module *terraform.Module) error {
	err := t.forEach(func(name string) (string, error) {
		rendered, err := t.template.Render(name, module)
		if err != nil {
			return "", fmt.Errorf("rendering template %q: %w", name, err)
		}

		return sanitize(rendered), nil
	})

	t.funcs(withModule(module))

	if err != nil {
		return fmt.Errorf("generating asciidoc table: %w", err)
	}

	return nil
}
