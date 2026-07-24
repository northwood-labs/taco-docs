// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package plugin

import (
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// Server is an RPC Server acting as a plugin.
type Server struct {
	impl   *formatter
	broker *goplugin.MuxBroker
}

// printFunc is a type alias that keeps plugin authors' code simple. They only
// need to provide a single function with this signature rather than implementing
// a full interface.
type printFunc func(*print.Config, *terraform.Module) (string, error)

// ServeOpts is an option for serving a plugin.
type ServeOpts struct {
	Name    string
	Version string
	Printer printFunc
}

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
func (s *Server) Name(args interface{}, resp *string) error {
	*resp = s.impl.Name()
	return nil
}

// Version returns the version of the plugin.
func (s *Server) Version(args interface{}, resp *string) error {
	*resp = s.impl.Version()
	return nil
}

// Execute returns the generated output.
func (s *Server) Execute(args *ExecuteArgs, resp *string) error {
	r, err := s.impl.Execute(args)
	*resp = r
	return err
}
