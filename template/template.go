// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package template

import (
	"bytes"
	"fmt"
	"strings"
	gotemplate "text/template"

	sprig "github.com/Masterminds/sprig/v3"

	"github.com/northwood-labs/taco-docs/internal/types"
	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// Item represents a named templated which can reference other named templated too.
type Item struct {
	Name      string
	Text      string
	TrimSpace bool
}

// Template wraps Go's text/template with terraform-docs-specific concerns:
// section rendering, config-aware functions, and custom function registration.
// This abstraction lets formatters focus on defining template text without
// worrying about function wiring or template composition.
type Template struct {
	items  []*Item
	config *print.Config

	funcMap    gotemplate.FuncMap
	customFunc gotemplate.FuncMap
}

// New returns new instance of Template.
func New(config *print.Config, items ...*Item) *Template {
	return &Template{
		items:      items,
		config:     config,
		funcMap:    builtinFuncs(config),
		customFunc: make(gotemplate.FuncMap),
	}
}

// Funcs return available template out of the box and custom functions.
func (t *Template) Funcs() gotemplate.FuncMap {
	return t.funcMap
}

// CustomFunc allows formatters to add format-specific functions without
// modifying the base template engine. Functions are only added if they don't
// already exist, preventing accidental overrides of built-in functions.
func (t *Template) CustomFunc(funcs gotemplate.FuncMap) {
	for name, fn := range funcs {
		if _, found := t.customFunc[name]; !found {
			t.customFunc[name] = fn
		}
	}
	t.applyCustomFunc()
}

// applyCustomFunc is re-adding the custom functions to list of available functions.
func (t *Template) applyCustomFunc() {
	for name, fn := range t.customFunc {
		if _, found := t.funcMap[name]; !found {
			t.funcMap[name] = fn
		}
	}
}

// Render is the high-level convenience method for rendering a terraform.Module.
// It wraps Module and Config into the standard data structure that templates expect.
func (t *Template) Render(name string, module *terraform.Module) (string, error) {
	data := struct {
		Config *print.Config
		Module *terraform.Module
	}{
		Config: t.config,
		Module: module,
	}
	return t.RenderContent(name, data)
}

// RenderContent is the low-level rendering method for arbitrary data. Used by
// generator.Render for content templates where the data isn't necessarily a
// terraform.Module (e.g., custom content templates with mixed data).
func (t *Template) RenderContent(name string, data interface{}) (string, error) {
	if len(t.items) < 1 {
		return "", fmt.Errorf("base template not found")
	}

	item := t.findByName(name)
	if item == nil {
		return "", fmt.Errorf("%s template not found", name)
	}

	var buffer bytes.Buffer

	tmpl := gotemplate.New(item.Name)
	tmpl.Funcs(t.funcMap)
	gotemplate.Must(tmpl.Parse(normalize(item.Text, item.TrimSpace)))

	for _, ii := range t.items {
		tt := tmpl.New(ii.Name)
		tt.Funcs(t.funcMap)
		gotemplate.Must(tt.Parse(normalize(ii.Text, ii.TrimSpace)))
	}

	if err := tmpl.ExecuteTemplate(&buffer, item.Name, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (t *Template) findByName(name string) *Item {
	if name == "" {
		if len(t.items) > 0 {
			return t.items[0]
		}
		return nil
	}
	for _, i := range t.items {
		if i.Name == name {
			return i
		}
	}
	return nil
}

// builtinFuncs provides the standard function library available in all templates.
// It includes sprig for rich text processing (date formatting, string manipulation,
// etc.) as a fallback, while terraform-docs-specific functions take priority to
// ensure correct behavior for documentation generation.
func builtinFuncs(config *print.Config) gotemplate.FuncMap { // nolint:gocyclo
	fns := gotemplate.FuncMap{
		"default": func(_default string, value string) string {
			if value != "" {
				return value
			}
			return _default
		},
		"indent": func(extra int, char string) string {
			return GenerateIndentation(config.Settings.Indent, extra, char)
		},
		"name": func(name string) string {
			return SanitizeName(name, config.Settings.Escape)
		},
		"ternary": func(condition interface{}, trueValue string, falseValue string) string {
			var c bool
			switch x := fmt.Sprintf("%T", condition); x {
			case "string":
				c = condition.(string) != ""
			case "int":
				c = condition.(int) != 0
			case "bool":
				c = condition.(bool)
			}
			if c {
				return trueValue
			}
			return falseValue
		},
		"tostring": func(s types.String) string {
			return string(s)
		},

		// trim
		"trim": func(cut string, s string) string {
			if s != "" {
				return strings.Trim(s, cut)
			}
			return s
		},
		"trimLeft": func(cut string, s string) string {
			if s != "" {
				return strings.TrimLeft(s, cut)
			}
			return s
		},
		"trimRight": func(cut string, s string) string {
			if s != "" {
				return strings.TrimRight(s, cut)
			}
			return s
		},
		"trimPrefix": func(prefix string, s string) string {
			if s != "" {
				return strings.TrimPrefix(s, prefix)
			}
			return s
		},
		"trimSuffix": func(suffix string, s string) string {
			if s != "" {
				return strings.TrimSuffix(s, suffix)
			}
			return s
		},

		// sanitize
		"sanitizeSection": func(s string) string {
			return SanitizeSection(s, config.Settings.Escape, config.Settings.HTML)
		},
		"sanitizeDoc": func(s string) string {
			return SanitizeDocument(s, config.Settings.Escape, config.Settings.HTML)
		},
		"sanitizeMarkdownTbl": func(s string) string {
			return SanitizeMarkdownTable(s, config.Settings.Escape, config.Settings.HTML)
		},
		"sanitizeAsciidocTbl": func(s string) string {
			return SanitizeAsciidocTable(s, config.Settings.Escape, config.Settings.HTML)
		},

		// anchors
		"anchorNameMarkdown": func(prefix string, value string) string {
			return CreateAnchorMarkdown(prefix, value, config.Settings.Anchor, config.Settings.Escape)
		},
		"anchorNameAsciidoc": func(prefix string, value string) string {
			return CreateAnchorAsciidoc(prefix, value, config.Settings.Anchor, config.Settings.Escape)
		},
	}

	for name, fn := range sprig.FuncMap() {
		if _, found := fns[name]; !found {
			fns[name] = fn
		}
	}

	return fns
}

// normalize strips leading whitespace from template lines so that indented Go
// source templates don't produce indented output. This makes it possible to
// write human-readable, well-indented template definitions in .go files without
// that indentation leaking into the generated documentation.
func normalize(s string, trimSpace bool) string {
	if !trimSpace {
		return s
	}
	split := strings.Split(s, "\n")
	for i, v := range split {
		split[i] = strings.TrimSpace(v)
	}
	return strings.Join(split, "\n")
}

// GenerateIndentation generates parameterized heading depth so users can embed
// generated docs at any nesting level in their documents. For example, if docs
// are placed under a "## Terraform" section, indent=3 makes inputs/outputs
// render as "###" headings, maintaining proper document hierarchy.
func GenerateIndentation(base int, extra int, char string) string {
	if char == "" {
		return ""
	}
	if base < 1 || base > 5 {
		base = 2
	}
	var indent string
	for i := 0; i < base+extra; i++ {
		indent += char
	}
	return indent
}
