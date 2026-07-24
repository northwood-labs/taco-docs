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
	xmlsdk "encoding/xml"
	"strings"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// xml represents XML format.
//
// WHY: Enterprise and legacy systems (XSLT pipelines, SOAP services,
// Java-based documentation generators) often require XML input. This
// formatter bridges Terraform module metadata into those ecosystems
// without requiring users to write custom converters.
type xml struct {
	*generator

	config *print.Config
}

// NewXML returns new instance of XML.
//
// WHY: canRender is false because XML's structure is dictated
// by MarshalIndent; custom templates would break well-formedness.
func NewXML(config *print.Config) Type {
	return &xml{
		generator: newGenerator(config, false),
		config:    config,
	}
}

// Generate a Terraform module as xml.
func (x *xml) Generate(module *terraform.Module) error {
	copy := copySections(x.config, module)

	out, err := xmlsdk.MarshalIndent(copy, "", "  ")
	if err != nil {
		return err
	}

	x.funcs(withContent(strings.TrimSuffix(string(out), "\n")))

	return nil
}

func init() {
	register(map[string]initializerFn{
		"xml": NewXML,
	})
}
