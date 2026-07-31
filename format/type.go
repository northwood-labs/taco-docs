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
	"fmt"
	"maps"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// initializers list of all registered engine initializer functions.
//
// A package-level map acts as a formatter registry. Each formatter
// registers itself via init(), so adding a new format requires zero changes
// to this file or any central import list—the Go linker pulls in the init()
// automatically when the package is imported.
var initializers = make(map[string]initializerFn)

// Sections provides read access to the individual documentation
// sections that a formatter produces. Each method returns the
// formatted text for that section.
type (
	Sections interface {
		// Header returns the header section based on the underlying format.
		Header() string

		// Footer returns the footer section based on the underlying format.
		Footer() string

		// Inputs returns the inputs section based on the underlying format.
		Inputs() string

		// Modules returns the modules section based on the underlying format.
		Modules() string

		// Outputs returns the outputs section based on the underlying format.
		Outputs() string

		// Providers returns the providers section based on the underlying
		// format.
		Providers() string

		// Requirements returns the requirements section based on the underlying
		// format.
		Requirements() string

		// Resources returns the resources section based on the underlying
		// format.
		Resources() string
	}

	// Type represents an output format type (e.g., json, markdown table,
	// yaml, etc).
	//
	// A common interface enables the CLI to select any output
	// format at runtime without knowing the concrete type at compile
	// time. Each formatter only needs to satisfy this contract—Generate
	// populates sections, the accessors expose them, and Render allows
	// custom content templates to recompose those sections in
	// user-defined order. This polymorphism keeps the command layer
	// format-agnostic.
	Type interface {
		Sections

		// Generate populates the formatter's sections from the Terraform module.
		Generate(module *terraform.Module) error

		// Content returns all the sections combined based on the underlying format.
		Content() string

		// Render applies the given template to produce the final output.
		Render(tmpl string) (string, error)
	}

	// initializerFn returns a concrete implementation of an Engine.
	initializerFn func(*print.Config) Type
)

// asInitializer wraps a constructor that returns a concrete type into
// an initializerFn that returns the Type interface.
func asInitializer[T Type](fn func(*print.Config) T) initializerFn {
	return func(config *print.Config) Type {
		return fn(config)
	}
}

// register a formatter engine initializer function.
//
// Accepting a map allows a single formatter to register multiple
// aliases (e.g., "md", "markdown", "markdown table") in one call, reducing
// boilerplate while keeping the registry logic centralized here.
func register(e map[string]initializerFn) {
	if e == nil {
		return
	}

	maps.Copy(initializers, e)
}

// New initializes and returns the concrete implementation of
// format.Engine based on the provided 'name', for example for name
// of 'json' it will return '*format.JSON' through 'format.NewJSON'
// function.
//
// This factory decouples callers from concrete types. The caller only
// needs a config with a formatter name string; the registry resolves it to
// the right constructor, keeping the command layer ignorant of which
// formatters exist.
func New(config *print.Config) (Type, error) {
	name := config.Formatter

	fn, ok := initializers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrFormatterNotFound, name)
	}

	return fn(config), nil
}
