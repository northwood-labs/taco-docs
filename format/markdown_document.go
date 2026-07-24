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
	"embed"
	gotemplate "text/template"

	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/template"
	"github.com/terraform-docs/terraform-docs/terraform"
)

//go:embed templates/markdown_document*.tmpl
var markdownDocumentFS embed.FS

// markdownDocument represents Markdown Document format.
//
// WHY: For modules with lengthy descriptions, complex types, or many
// validation rules, a table becomes unreadable. The document format
// gives each input/output its own subsection with full-width rendering,
// making verbose documentation scannable via a table of contents.
type markdownDocument struct {
	*generator

	config   *print.Config
	template *template.Template
}

// NewMarkdownDocument returns new instance of Markdown Document.
func NewMarkdownDocument(config *print.Config) Type {
	items := readTemplateItems(markdownDocumentFS, "markdown_document")

	tt := template.New(config, items...)
	tt.CustomFunc(gotemplate.FuncMap{
		// WHY: Multi-line types/values need fenced code blocks (```hcl)
		// for proper syntax highlighting in document view, unlike the
		// table format which uses inline backticks for compactness.
		"type": func(t string) string {
			result, extraline := PrintFencedCodeBlock(t, "hcl")
			if !extraline {
				result += "\n"
			}
			return result
		},
		"value": func(v string) string {
			if v == "n/a" {
				return v
			}
			result, extraline := PrintFencedCodeBlock(v, "json")
			if !extraline {
				result += "\n"
			}
			return result
		},
		// WHY: The "Required" badge is only shown when the user opts
		// in via config, since it adds noise for modules where most
		// inputs are optional.
		"isRequired": func() bool {
			return config.Settings.Required
		},
	})

	return &markdownDocument{
		generator: newGenerator(config, true),
		config:    config,
		template:  tt,
	}
}

// Generate a Terraform module as Markdown document.
func (d *markdownDocument) Generate(module *terraform.Module) error {
	err := d.forEach(func(name string) (string, error) {
		rendered, err := d.template.Render(name, module)
		if err != nil {
			return "", err
		}
		return sanitize(rendered), nil
	})

	d.funcs(withModule(module))

	return err
}

func init() {
	register(map[string]initializerFn{
		"markdown document": NewMarkdownDocument,
		"markdown doc":      NewMarkdownDocument,
		"md document":       NewMarkdownDocument,
		"md doc":            NewMarkdownDocument,
	})
}
