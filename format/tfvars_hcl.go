/*
Copyright 2021 The terraform-docs Authors.
Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package format

import (
	_ "embed" //nolint
	"fmt"
	"reflect"
	"strings"
	gotemplate "text/template"

	"github.com/terraform-docs/terraform-docs/internal/types"
	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/template"
	"github.com/terraform-docs/terraform-docs/terraform"
)

//go:embed templates/tfvars_hcl.tmpl
var tfvarsHCLTpl []byte

// tfvarsHCL represents Terraform tfvars HCL format.
//
// WHY: This formatter produces a ready-to-use .tfvars file that users can
// immediately feed to `terraform plan -var-file=`. It saves users from
// manually transcribing variable names and default values when scaffolding
// a new deployment from a module's documentation.
type tfvarsHCL struct {
	*generator

	config   *print.Config
	template *template.Template
}

// WHY: padding is package-level because it's computed once per Generate
// call and consumed by the "align" template function. Passing it through
// template context would require a custom struct; a package var keeps the
// template function signature simple.
var padding []int

// NewTfvarsHCL returns new instance of TfvarsHCL.
func NewTfvarsHCL(config *print.Config) Type {
	tt := template.New(config, &template.Item{
		Name:      "tfvars",
		Text:      string(tfvarsHCLTpl),
		TrimSpace: true,
	})
	tt.CustomFunc(gotemplate.FuncMap{
		// WHY: "align" pads variable names to the same width within each
		// alignment group, producing idiomatic HCL that reads like a
		// hand-formatted .tfvars file.
		"align": func(s string, i int) string {
			return fmt.Sprintf("%-*s", padding[i], s)
		},
		// WHY: Empty defaults must render as "" (empty string literal) rather
		// than blank, otherwise the HCL would be syntactically invalid.
		"value": func(s string) string {
			if s == "" {
				return "\"\""
			}
			return s
		},
		// WHY: Descriptions are rendered as HCL comments (# prefix) above
		// each variable so users understand what they're filling in without
		// switching back to the module source.
		"convertToComment": func(s types.String) string {
			return "\n# " + strings.ReplaceAll(string(s), "\n", "\n# ")
		},
		"showDescription": func() bool {
			return config.Settings.Description
		},
	})

	return &tfvarsHCL{
		generator: newGenerator(config, false),
		config:    config,
		template:  tt,
	}
}

// Generate a Terraform module as Terraform tfvars HCL.
//
// WHY: alignments() is called before rendering so the template can use
// pre-computed padding widths. Sanitize is applied to handle any stray
// whitespace the template may emit.
func (h *tfvarsHCL) Generate(module *terraform.Module) error {
	alignments(module.Inputs, h.config)

	rendered, err := h.template.Render("tfvars", module)
	if err != nil {
		return err
	}

	h.funcs(withContent(strings.TrimSuffix(sanitize(rendered), "\n")))

	return nil
}

// isMultilineFormat checks if an input's default value will span multiple
// lines when serialized.
//
// WHY: Multi-line values (lists, maps with elements) break visual alignment
// with neighboring single-line assignments. The alignment algorithm uses this
// to reset padding groups at multi-line boundaries, preventing awkward gaps.
func isMultilineFormat(input *terraform.Input) bool {
	isList := input.Type == "list" || reflect.TypeOf(input.Default).Name() == "List"
	isMap := input.Type == "map" || reflect.TypeOf(input.Default).Name() == "Map"
	return (isList || isMap) && input.Default.Length() > 0
}

// alignments computes per-input padding widths for column alignment.
//
// WHY: HCL style conventions align the `=` sign for consecutive simple
// assignments but reset alignment when a multi-line value or comment block
// intervenes. This groups inputs into "alignment runs" separated by multi-line
// values or descriptions, computing the max name length within each run.
func alignments(inputs []*terraform.Input, config *print.Config) {
	padding = make([]int, len(inputs))
	maxlen := 0
	index := 0
	for i, input := range inputs {
		isDescribed := config.Settings.Description && input.Description.Length() > 0
		l := len(input.Name)
		if isMultilineFormat(input) || isDescribed {
			// WHY: When we hit a boundary (multi-line or described input),
			// flush the current group by applying the accumulated max to
			// all preceding entries, then start a new group.
			for j := index; j < i; j++ {
				padding[j] = maxlen
			}
			padding[i] = l
			maxlen = 0
			index = i + 1
		} else if l > maxlen {
			maxlen = l
		}
	}
	// WHY: Flush the final group—inputs after the last boundary still need
	// their padding set.
	for i := index; i < len(inputs); i++ {
		padding[i] = maxlen
	}
}

func init() {
	register(map[string]initializerFn{
		"tfvars hcl": NewTfvarsHCL,
	})
}
