/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package format

import (
	"fmt"

	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
)

// Type represents an output format type (e.g. json, markdown table, yaml, etc).
//
// WHY: A common interface enables the CLI to select any output format at runtime
// without knowing the concrete type at compile time. Each formatter only needs to
// satisfy this contract—Generate populates sections, the accessors expose them, and
// Render allows custom content templates to recompose those sections in user-defined
// order. This polymorphism keeps the command layer format-agnostic.
type Type interface {
	Generate(*terraform.Module) error // generate the Terraform module

	Content() string // all the sections combined based on the underlying format

	Header() string       // header section based on the underlying format
	Footer() string       // footer section based on the underlying format
	Inputs() string       // inputs section based on the underlying format
	Modules() string      // modules section based on the underlying format
	Outputs() string      // outputs section based on the underlying format
	Providers() string    // providers section based on the underlying format
	Requirements() string // requirements section based on the underlying format
	Resources() string    // resources section based on the underlying format

	Render(tmpl string) (string, error)
}

// initializerFn returns a concrete implementation of an Engine.
type initializerFn func(*print.Config) Type

// initializers list of all registered engine initializer functions.
//
// WHY: A package-level map acts as a formatter registry. Each formatter registers
// itself via init(), so adding a new format requires zero changes to this file or
// any central import list—the Go linker pulls in the init() automatically when the
// package is imported.
var initializers = make(map[string]initializerFn)

// register a formatter engine initializer function.
//
// WHY: Accepting a map allows a single formatter to register multiple aliases
// (e.g. "md", "markdown", "markdown table") in one call, reducing boilerplate
// while keeping the registry logic centralized here.
func register(e map[string]initializerFn) {
	if e == nil {
		return
	}
	for k, v := range e {
		initializers[k] = v
	}
}

// New initializes and returns the concrete implementation of
// format.Engine based on the provided 'name', for example for name
// of 'json' it will return '*format.JSON' through 'format.NewJSON'
// function.
//
// WHY: This factory decouples callers from concrete types. The caller only needs
// a config with a formatter name string; the registry resolves it to the right
// constructor, keeping the command layer ignorant of which formatters exist.
func New(config *print.Config) (Type, error) {
	name := config.Formatter
	fn, ok := initializers[name]
	if !ok {
		return nil, fmt.Errorf("formatter '%s' not found", name)
	}
	return fn(config), nil
}
