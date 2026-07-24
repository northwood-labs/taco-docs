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
	"strings"

	tomlsdk "github.com/BurntSushi/toml"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// toml represents TOML format.
//
// WHY: TOML is the native configuration language for tools like
// Cargo, Hugo, and pyproject.toml-based Python workflows. Providing
// TOML output lets these ecosystems ingest module metadata directly
// without format conversion.
type toml struct {
	*generator

	config *print.Config
}

// NewTOML returns new instance of TOML.
//
// WHY: canRender is false because the TOML encoder controls
// document structure; custom templates can't reorder it.
func NewTOML(config *print.Config) Type {
	return &toml{
		generator: newGenerator(config, false),
		config:    config,
	}
}

// Generate a Terraform module as toml.
func (t *toml) Generate(module *terraform.Module) error {
	copy := copySections(t.config, module)

	buffer := new(bytes.Buffer)
	encoder := tomlsdk.NewEncoder(buffer)

	if err := encoder.Encode(copy); err != nil {
		return err
	}

	t.funcs(withContent(strings.TrimSuffix(buffer.String(), "\n")))

	return nil
}

func init() {
	register(map[string]initializerFn{
		"toml": NewTOML,
	})
}
