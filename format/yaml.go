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
	"fmt"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// yaml represents YAML format.
//
// YAML is the preferred structured data format in many infrastructure
// toolchains (Ansible, Helm, GitHub Actions). Providing native YAML output lets
// these ecosystems consume module metadata without a JSON-to-YAML conversion
// step, preserving human-readability.
type yaml struct {
	*generator

	config *print.Config
}

func init() { // lint:allow_init
	register(map[string]initializerFn{
		"yaml": asInitializer(NewYAML),
	})
}

// NewYAML returns new instance of YAML.
//
// canRender is false because the YAML encoder dictates output structure; custom
// content templates can't meaningfully reorder a serialized YAML document.
func NewYAML(config *print.Config) *yaml { // lint:allow_unexported_return
	return &yaml{
		generator: newGenerator(config, false),
		config:    config,
	}
}

// Generate a Terraform module as YAML.
//
// The 2-space indent matches the YAML community convention and keeps the output
// consistent with typical Kubernetes/Helm manifests users work with.
func (y *yaml) Generate(module *terraform.Module) error {
	sections := copySections(y.config, module)

	buffer := new(bytes.Buffer)
	encoder := yamlv3.NewEncoder(buffer)
	encoder.SetIndent(2)

	if err := encoder.Encode(sections); err != nil {
		return fmt.Errorf("encoding module as YAML: %w", err)
	}

	y.funcs(withContent(strings.TrimSuffix(buffer.String(), "\n")))

	return nil
}
