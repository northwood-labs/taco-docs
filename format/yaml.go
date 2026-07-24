/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package format

import (
	"bytes"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"

	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
)

// yaml represents YAML format.
//
// WHY: YAML is the preferred structured data format in many infrastructure
// toolchains (Ansible, Helm, GitHub Actions). Providing native YAML output
// lets these ecosystems consume module metadata without a JSON-to-YAML
// conversion step, preserving human-readability.
type yaml struct {
	*generator

	config *print.Config
}

// NewYAML returns new instance of YAML.
//
// WHY: canRender is false because the YAML encoder dictates output structure;
// custom content templates can't meaningfully reorder a serialized YAML document.
func NewYAML(config *print.Config) Type {
	return &yaml{
		generator: newGenerator(config, false),
		config:    config,
	}
}

// Generate a Terraform module as YAML.
//
// WHY: The 2-space indent matches the YAML community convention and keeps the
// output consistent with typical Kubernetes/Helm manifests users work with.
func (y *yaml) Generate(module *terraform.Module) error {
	copy := copySections(y.config, module)

	buffer := new(bytes.Buffer)
	encoder := yamlv3.NewEncoder(buffer)
	encoder.SetIndent(2)

	if err := encoder.Encode(copy); err != nil {
		return err
	}

	y.funcs(withContent(strings.TrimSuffix(buffer.String(), "\n")))

	return nil
}

func init() {
	register(map[string]initializerFn{
		"yaml": NewYAML,
	})
}
