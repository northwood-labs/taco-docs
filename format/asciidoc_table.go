/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package format

import (
	"embed"
	gotemplate "text/template"

	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/template"
	"github.com/terraform-docs/terraform-docs/terraform"
)

//go:embed templates/asciidoc_table*.tmpl
var asciidocTableFS embed.FS

// asciidocTable represents AsciiDoc Table format.
//
// WHY: Teams using Asciidoctor or Antora for their documentation sites
// need native AsciiDoc output rather than embedded Markdown. This is the
// compact tabular equivalent of markdownTable for the AsciiDoc ecosystem.
type asciidocTable struct {
	*generator

	config   *print.Config
	template *template.Template
}

// NewAsciidocTable returns new instance of Asciidoc Table.
func NewAsciidocTable(config *print.Config) Type {
	items := readTemplateItems(asciidocTableFS, "asciidoc_table")

	// WHY: AsciiDoc has its own escaping rules (e.g. | inside table cells).
	// Disabling the generic markdown escape prevents double-escaping that
	// would corrupt the AsciiDoc output.
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
			return "", err
		}
		return sanitize(rendered), nil
	})

	t.funcs(withModule(module))

	return err
}

// WHY: Multiple short aliases registered so users can type "adoc" or
// "asciidoc" interchangeably, matching common community shorthand.
func init() {
	register(map[string]initializerFn{
		"asciidoc":       NewAsciidocTable,
		"asciidoc table": NewAsciidocTable,
		"asciidoc tbl":   NewAsciidocTable,
		"adoc":           NewAsciidocTable,
		"adoc table":     NewAsciidocTable,
		"adoc tbl":       NewAsciidocTable,
	})
}
