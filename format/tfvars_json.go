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
	"bytes"
	jsonsdk "encoding/json"
	"fmt"
	"strings"

	"github.com/iancoleman/orderedmap"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// tfvarsJSON represents Terraform tfvars JSON format.
//
// CI pipelines and automation tools (Terragrunt, Atlantis, custom scripts)
// often generate variable values programmatically. A JSON .tfvars.json file is
// easier to produce and parse from code than HCL, making this the preferred
// format for machine-generated variable definitions.
type tfvarsJSON struct {
	*generator

	config *print.Config
}

func init() { // lint:allow_init
	register(map[string]initializerFn{
		"tfvars json": asInitializer(NewTfvarsJSON),
	})
}

// NewTfvarsJSON returns new instance of TfvarsJSON.
//
// canRender is false because the output is a flat JSON object of
// variable-name-to-default-value pairs—there are no sections to reorder.
func NewTfvarsJSON(config *print.Config) *tfvarsJSON { // lint:allow_unexported_return
	return &tfvarsJSON{
		generator: newGenerator(config, false),
		config:    config,
	}
}

// Generate a Terraform module as Terraform tfvars JSON.
//
// An ordered map preserves the declaration order of inputs from the module
// source. This makes diffs stable and the output predictable—users won't see
// spurious reorderings on regeneration. SetEscapeHTML(false) prevents Go's
// encoder from mangling URLs or HTML in default values.
func (j *tfvarsJSON) Generate(module *terraform.Module) error {
	sections := orderedmap.New()
	sections.SetEscapeHTML(false)

	for _, i := range module.Inputs {
		sections.Set(i.Name, i.Default)
	}

	buffer := new(bytes.Buffer)
	encoder := jsonsdk.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(sections); err != nil {
		return fmt.Errorf("encoding tfvars as JSON: %w", err)
	}

	j.funcs(withContent(strings.TrimSuffix(buffer.String(), "\n")))

	return nil
}
