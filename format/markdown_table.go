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

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/template"
	"github.com/northwood-labs/taco-docs/terraform"
)

//go:embed templates/markdown_table*.tmpl
var markdownTableFS embed.FS

// markdownTable represents Markdown Table format.
//
// WHY: This is the most popular output format because GitHub, GitLab,
// and Bitbucket all render Markdown tables natively in README files.
// The compact tabular layout fits modules with many variables while
// remaining scannable without scrolling through verbose per-variable
// subsections.
type markdownTable struct {
	*generator

	config   *print.Config
	template *template.Template
}

// NewMarkdownTable returns new instance of Markdown Table.
func NewMarkdownTable(config *print.Config) Type {
	items := readTemplateItems(markdownTableFS, "markdown_table")

	tt := template.New(config, items...)
	tt.CustomFunc(gotemplate.FuncMap{
		// WHY: Types and values must be rendered as inline code so they
		// don't break table cell alignment or get interpreted as Markdown
		// formatting (e.g. a type containing `*` would become italic).
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

	return &markdownTable{
		generator: newGenerator(config, true),
		config:    config,
		template:  tt,
	}
}

// Generate a Terraform module as Markdown tables.
//
// WHY: forEach renders each section independently so that users with
// custom content templates can reference individual sections (e.g.
// {{ .Inputs }}) and reorder them freely.
func (t *markdownTable) Generate(module *terraform.Module) error {
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

// WHY: Multiple aliases ("markdown", "md", "md table", etc.) are registered
// because this is the default and most common format. Users expect short
// names to work, and older tutorials reference different variants.
func init() {
	register(map[string]initializerFn{
		"markdown":       NewMarkdownTable,
		"markdown table": NewMarkdownTable,
		"markdown tbl":   NewMarkdownTable,
		"md":             NewMarkdownTable,
		"md table":       NewMarkdownTable,
		"md tbl":         NewMarkdownTable,
	})
}
