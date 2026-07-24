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
	"bytes"
	jsonsdk "encoding/json"
	"strings"

	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
)

// json represents JSON format.
//
// WHY: JSON output enables programmatic consumption of module documentation—CI
// pipelines, custom renderers, and automation tools can parse it without
// scraping Markdown. It's the canonical machine-readable format.
type json struct {
	*generator

	config *print.Config
}

// NewJSON returns new instance of JSON.
//
// WHY: canRender is false because JSON's structure is dictated by the encoder;
// custom content templates would produce invalid JSON.
func NewJSON(config *print.Config) Type {
	return &json{
		generator: newGenerator(config, false),
		config:    config,
	}
}

// Generate a Terraform module as json.
//
// WHY: copySections is called first to honor the user's show/hide configuration,
// then the filtered module is serialized. SetEscapeHTML defers to the user's
// escape setting so HTML entities in descriptions aren't mangled unless requested.
func (j *json) Generate(module *terraform.Module) error {
	copy := copySections(j.config, module)

	buffer := new(bytes.Buffer)
	encoder := jsonsdk.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(j.config.Settings.Escape)

	if err := encoder.Encode(copy); err != nil {
		return err
	}

	j.funcs(withContent(strings.TrimSuffix(buffer.String(), "\n")))

	return nil
}

func init() {
	register(map[string]initializerFn{
		"json": NewJSON,
	})
}
