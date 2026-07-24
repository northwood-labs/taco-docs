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
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/template"
	"github.com/terraform-docs/terraform-docs/terraform"
	"mvdan.cc/xurls/v2"
)

// sanitize cleans a Markdown document to soothe linters.
//
// WHY: Generated markdown is injected into existing repos that typically enforce
// markdownlint or similar rules. Trailing spaces, excessive blank lines, and bare
// URLs would cause CI linter failures for users. By normalizing whitespace and
// wrapping URLs here, the output is linter-friendly out of the box.
func sanitize(markdown string) string {
	result := markdown

	// WHY: Double trailing spaces are intentional line breaks in Markdown (soft break).
	// We must preserve them while stripping all other trailing whitespace.
	result = regexp.MustCompile(` {2}(\r?\n)`).ReplaceAllString(result, "‡‡‡DOUBLESPACES‡‡‡$1")

	// Remove trailing spaces from the end of lines
	result = regexp.MustCompile(` +(\r?\n)`).ReplaceAllString(result, "$1")
	result = regexp.MustCompile(` +$`).ReplaceAllLiteralString(result, "")

	// Restore the preserved double spaces
	result = regexp.MustCompile(`‡‡‡DOUBLESPACES‡‡‡(\r?\n)`).ReplaceAllString(result, "  $1")

	// WHY: A blank line containing only double spaces is an artifact of template
	// rendering; it serves no formatting purpose and triggers linter warnings.
	result = regexp.MustCompile(`(\r?\n)  (\r?\n)`).ReplaceAllString(result, "$1")

	// WHY: Multiple consecutive blank lines add no semantic value in generated
	// docs and violate most markdown style guides (MD012).
	result = regexp.MustCompile(`(\r?\n){3,}`).ReplaceAllString(result, "$1$1")
	result = regexp.MustCompile(`(\r?\n){2,}$`).ReplaceAllString(result, "")

	result = SanitizeBareLinks(result)

	return result
}

// SanitizeBareLinks converts bare links to Markdown representation.
//
// WHY: Bare URLs (without angle brackets or markdown link syntax) violate
// markdownlint rule MD034. Wrapping them in <> produces valid auto-links that
// render identically in all Markdown renderers while satisfying linters.
func SanitizeBareLinks(s string) string {
	urlRegex := xurls.Strict()
	matches := urlRegex.FindAllStringIndex(s, -1)
	if matches == nil {
		return s
	}

	var result strings.Builder
	lastIndex := 0

	for _, match := range matches {
		start, end := match[0], match[1]

		// WHY: Skip URLs already wrapped in angle brackets—re-wrapping would
		// produce <<url>> which is invalid Markdown.
		if start > 0 && s[start-1] == '<' && end < len(s) && s[end] == '>' {
			continue
		}

		// WHY: Skip URLs already inside markdown link syntax ]({url})—these
		// are intentional references and don't need auto-link wrapping.
		if start > 1 && s[start-2:start] == "](" && end < len(s) && s[end] == ')' {
			continue
		}

		// Append text before the URL
		result.WriteString(s[lastIndex:start])

		// Wrap the URL in <>
		url := s[start:end]
		result.WriteString("<")
		result.WriteString(url)
		result.WriteString(">")

		lastIndex = end
	}

	// Append the remaining part of the line
	result.WriteString(s[lastIndex:])
	return result.String()
}

// PrintFencedCodeBlock prints codes in fences, it automatically detects if
// the input 'code' contains '\n' it will use multi line fence, otherwise it
// wraps the 'code' inside single-tick block.
// If the fenced is multi-line it also appends an extra '\n` at the end and
// returns true accordingly, otherwise returns false for non-carriage return.
//
// WHY: Terraform types and default values can be single-line primitives or
// multi-line complex objects. Using single backticks for short values keeps
// tables compact, while triple-fence blocks preserve readability for HCL maps
// and lists. The boolean return signals callers whether extra spacing is needed.
func PrintFencedCodeBlock(code string, language string) (string, bool) {
	if strings.Contains(code, "\n") {
		return fmt.Sprintf("\n\n```%s\n%s\n```\n", language, code), true
	}
	return fmt.Sprintf("`%s`", code), false
}

// PrintFencedAsciidocCodeBlock prints codes in fences, it automatically detects if
// the input 'code' contains '\n' it will use multi line fence, otherwise it
// wraps the 'code' inside single-tick block.
// If the fenced is multi-line it also appends an extra '\n` at the end and
// returns true accordingly, otherwise returns false for non-carriage return.
//
// WHY: AsciiDoc uses a different fence syntax ([source,lang] + ---- delimiters).
// This parallel function ensures AsciiDoc formatters produce valid Asciidoctor
// source blocks without duplicating the single-vs-multi decision logic.
func PrintFencedAsciidocCodeBlock(code string, language string) (string, bool) {
	if strings.Contains(code, "\n") {
		return fmt.Sprintf("\n[source,%s]\n----\n%s\n----\n", language, code), true
	}
	return fmt.Sprintf("`%s`", code), false
}

// readTemplateItems reads all static formatter .tmpl files prefixed by specific string
// from an embed file system.
//
// WHY: Each template-based formatter stores its section templates as separate
// embedded files (e.g. markdown_table_inputs.tmpl). This function abstracts the
// embed.FS traversal so formatters don't duplicate file-glob-and-parse logic.
// The prefix stripping and underscore removal normalize filenames into the
// canonical section names ("inputs", "outputs", etc.) expected by forEach.
func readTemplateItems(efs embed.FS, prefix string) []*template.Item {
	items := make([]*template.Item, 0)

	files, err := fs.ReadDir(efs, "templates")
	if err != nil {
		return items
	}

	for _, f := range files {
		content, err := efs.ReadFile("templates/" + f.Name())
		if err != nil {
			continue
		}

		name := f.Name()
		name = strings.ReplaceAll(name, prefix, "")
		name = strings.ReplaceAll(name, "_", "")
		name = strings.ReplaceAll(name, ".tmpl", "")
		// WHY: The base template (no section suffix) renders all sections
		// combined. Naming it "all" matches the forEach mapping key.
		if name == "" {
			name = "all"
		}

		items = append(items, &template.Item{
			Name:      name,
			Text:      string(content),
			TrimSpace: true,
		})
	}
	return items
}

// copySections sets the sections that'll be printed
//
// WHY: Users configure which sections to show or hide via the config file or CLI
// flags. Rather than threading conditional logic through every template, we build
// a filtered copy of the module upfront. Templates then render unconditionally
// against a module that already contains only the desired sections—keeping
// template logic simple and the show/hide decision in one place.
func copySections(config *print.Config, src *terraform.Module) *terraform.Module {
	dest := &terraform.Module{
		Header:            "",
		Footer:            "",
		Inputs:            make([]*terraform.Input, 0),
		ModuleCalls:       make([]*terraform.ModuleCall, 0),
		Outputs:           make([]*terraform.Output, 0),
		Providers:         make([]*terraform.Provider, 0),
		ProviderFunctions: make([]*terraform.ProviderFunction, 0),
		Requirements:      make([]*terraform.Requirement, 0),
		Resources:         make([]*terraform.Resource, 0),
	}

	if config.Sections.Header {
		dest.Header = src.Header
	}
	if config.Sections.Footer {
		dest.Footer = src.Footer
	}
	if config.Sections.Inputs {
		dest.Inputs = src.Inputs
	}
	if config.Sections.ModuleCalls {
		dest.ModuleCalls = src.ModuleCalls
	}
	if config.Sections.Outputs {
		dest.Outputs = src.Outputs
	}
	if config.Sections.Providers {
		dest.Providers = src.Providers
	}
	if config.Sections.ProviderFunctions {
		dest.ProviderFunctions = src.ProviderFunctions
	}
	if config.Sections.Requirements {
		dest.Requirements = src.Requirements
	}
	if config.Sections.Resources || config.Sections.DataSources {
		dest.Resources = filterResourcesByMode(config, src.Resources)
	}

	return dest
}

// filterResourcesByMode returns the managed or data resources defined by the show argument
//
// WHY: Terraform distinguishes managed resources ("resource" blocks) from data
// sources ("data" blocks) by mode. Users may want to document one, both, or
// neither independently. Filtering here ensures the template receives only the
// resource types the user asked for.
func filterResourcesByMode(config *print.Config, module []*terraform.Resource) []*terraform.Resource {
	resources := make([]*terraform.Resource, 0)
	for _, r := range module {
		if config.Sections.Resources && r.Mode == "managed" {
			resources = append(resources, r)
		}
		if config.Sections.DataSources && r.Mode == "data" {
			resources = append(resources, r)
		}
	}
	return resources
}
