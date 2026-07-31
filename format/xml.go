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
	xmlsdk "encoding/xml"
	"fmt"
	"strings"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// xml represents XML format.
//
// Enterprise and legacy systems (XSLT pipelines, SOAP services, Java-based
// documentation generators) often require XML input. This formatter bridges
// Terraform module metadata into those ecosystems without requiring users to
// write custom converters.
type xml struct {
	*generator

	config *print.Config
}

func init() { // lint:allow_init
	register(map[string]initializerFn{
		"xml": asInitializer(NewXML),
	})
}

// NewXML returns new instance of XML.
//
// canRender is false because XML's structure is dictated by MarshalIndent;
// custom templates would break well-formedness.
func NewXML(config *print.Config) *xml { // lint:allow_unexported_return
	return &xml{
		generator: newGenerator(config, false),
		config:    config,
	}
}

// Generate a Terraform module as xml.
func (x *xml) Generate(module *terraform.Module) error {
	sections := copySections(x.config, module)

	out, err := xmlsdk.MarshalIndent(sections, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding module as XML: %w", err)
	}

	x.funcs(withContent(strings.TrimSuffix(string(out), "\n")))

	return nil
}
