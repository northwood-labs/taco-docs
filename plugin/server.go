// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package plugin // lint:allow_naming_conflict_stdlib

import (
	"fmt"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

type (
	// Server is an RPC Server acting as a plugin.
	Server struct {
		impl   *formatter
		broker *goplugin.MuxBroker
	}

	// printFunc is a type alias that keeps plugin authors' code simple. They
	// only need to provide a single function with this signature rather than
	// implementing a full interface.
	printFunc func(*print.Config, *terraform.Module) (string, error)

	// ServeOpts is an option for serving a plugin.
	ServeOpts struct {
		Printer printFunc
		Name    string
		Version string
	}
)

// Serve is the single entry point for plugin binaries. One call sets up the
// entire RPC server, handshake, and connection lifecycle. Plugin authors call
// this from main() and never interact with go-plugin directly.
func Serve(opts *ServeOpts) {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: goplugin.PluginSet{
			"formatter": newFormatter(opts.Name, opts.Version, opts.Printer),
		},
	})
}

// Name, Version, and Execute map 1:1 to Client calls via go-plugin's RPC
// contract. The Server receives these calls and delegates to the formatter
// implementation, bridging the network boundary transparently.

// Name returns the version of the plugin.
func (s *Server) Name(_ any, resp *string) error { // lint:allow_param
	*resp = s.impl.Name()
	return nil
}

// Version returns the version of the plugin.
func (s *Server) Version(_ any, resp *string) error { // lint:allow_param
	*resp = s.impl.Version()
	return nil
}

// Execute returns the generated output.
func (s *Server) Execute(args *ExecuteArgs, resp *string) error {
	r, err := s.impl.Execute(args)

	*resp = r

	if err != nil {
		return fmt.Errorf("executing plugin formatter: %w", err)
	}

	return nil
}
