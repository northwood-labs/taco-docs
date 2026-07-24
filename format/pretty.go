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
	_ "embed" //nolint
	"fmt"
	"regexp"
	gotemplate "text/template"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/template"
	"github.com/northwood-labs/taco-docs/terraform"
)

//go:embed templates/pretty.tmpl
var prettyTpl []byte

// pretty represents colorized pretty format.
//
// WHY: During local development, authors need a quick visual overview
// of a module's interface without opening rendered Markdown. Colorized
// terminal output makes types, defaults, and descriptions easy to scan
// at a glance—this is the "human preview" format.
type pretty struct {
	*generator

	config   *print.Config
	template *template.Template
}

// NewPretty returns new instance of Pretty.
func NewPretty(config *print.Config) Type {
	tt := template.New(config, &template.Item{
		Name:      "pretty",
		Text:      string(prettyTpl),
		TrimSpace: true,
	})
	tt.CustomFunc(gotemplate.FuncMap{
		// WHY: colorize wraps text in ANSI escape codes only when
		// color is enabled, so the same template works for both TTY
		// and piped/redirected output without conditional template logic.
		"colorize": func(c string, s string) string {
			r := "\033[0m"
			if !config.Settings.Color {
				c = ""
				r = ""
			}
			return fmt.Sprintf("%s%s%s", c, s, r)
		},
	})

	return &pretty{
		generator: newGenerator(config, true),
		config:    config,
		template:  tt,
	}
}

// Generate a Terraform module document.
//
// WHY: Unlike file-based formatters, pretty output goes directly to
// stdout so trailing newlines must be stripped to avoid blank lines
// after the last entry when printed in a terminal.
func (p *pretty) Generate(module *terraform.Module) error {
	rendered, err := p.template.Render("pretty", module)
	if err != nil {
		return err
	}

	p.funcs(withContent(regexp.MustCompile(`(\r?\n)*$`).ReplaceAllString(rendered, "")))
	p.funcs(withModule(module))

	return nil
}

func init() {
	register(map[string]initializerFn{
		"pretty": NewPretty,
	})
}
