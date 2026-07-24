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
	"encoding/gob"
	"net/rpc"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/terraform-docs/terraform-docs/internal/types"
)

// Ensure formatter fully satisfy plugin interface.
var _ goplugin.Plugin = &formatter{}

// handshakeConfig acts as a magic cookie that prevents users from accidentally
// executing plugin binaries directly from the command line. go-plugin checks
// this cookie and refuses to serve if it's absent, providing a clear error
// instead of cryptic RPC failures. ProtocolVersion is bumped on breaking
// changes to prevent incompatible host/plugin combinations from connecting.
var handshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  7,
	MagicCookieKey:   "TFDOCS_PLUGIN",
	MagicCookieValue: "A7U5oTDDJwdL6UKOw6RXATDa86NEo4xLK3rz7QqegT1N4EY66qb6UeAJDSxLwtXH",
}

// formatter wraps the plugin's name, version, and print function to satisfy
// go-plugin's Plugin interface (Server/Client factory pattern). This indirection
// lets go-plugin manage the RPC lifecycle while keeping plugin authors' code
// focused on the formatting logic itself.
type formatter struct {
	name    string
	version string
	printer printFunc
}

func newFormatter(name string, version string, printer printFunc) *formatter {
	return &formatter{
		name:    name,
		version: version,
		printer: printer,
	}
}

func (f *formatter) Name() string {
	return f.name
}

func (f *formatter) Version() string {
	return f.version
}

func (f *formatter) Execute(args *ExecuteArgs) (string, error) {
	return f.printer(args.Config, args.Module)
}

// Server returns an RPC server acting as a plugin.
func (f *formatter) Server(b *goplugin.MuxBroker) (interface{}, error) {
	return &Server{impl: f, broker: b}, nil
}

// Client returns an RPC client for the host.
func (*formatter) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &Client{rpcClient: c, broker: b}, nil
}

// init registers custom types with gob encoding. go-plugin uses gob for RPC
// serialization; interface values (like terraform.Module fields that hold
// types.String, types.List, etc.) must be registered upfront or gob will fail
// to serialize them at runtime with an opaque "gob: type not registered" error.
func init() {
	gob.Register(new(types.Bool))
	gob.Register(new(types.Empty))
	gob.Register(new(types.List))
	gob.Register(new(types.Map))
	gob.Register(new(types.Nil))
	gob.Register(new(types.Number))
	gob.Register(new(types.String))
}
