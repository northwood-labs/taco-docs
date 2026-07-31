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
	"net/rpc"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

type (
	// Client wraps RPC calls behind a simple Go interface so host code doesn't
	// have to deal with raw RPC method strings or serialization details.
	Client struct {
		rpcClient *rpc.Client
		broker    *goplugin.MuxBroker
	}

	// ClientOpts is an option for initializing a Client.
	ClientOpts struct {
		Cmd *exec.Cmd
	}

	// ExecuteArgs bundles the data the host sends to plugins for generation.
	// Keeping it in a single struct simplifies the RPC contract—one argument,
	// one response—and makes the protocol easy to extend without breaking
	// existing plugins.
	ExecuteArgs struct {
		Module *terraform.Module
		Config *print.Config
	}
)

// NewClient configures hclog to write to stderr so plugin diagnostic logs
// don't contaminate the generated documentation output on stdout. Log level
// is controlled via the TFDOCS_LOG environment variable to support debugging
// without code changes.
func NewClient(opts *ClientOpts) *goplugin.Client {
	return goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"formatter": &formatter{},
		},
		Cmd: opts.Cmd,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: os.Stderr,
			Level:  hclog.LevelFromString(os.Getenv("TFDOCS_LOG")),
		}),
	})
}

// Name calls the server-side Name method and returns its version.
func (c *Client) Name() (string, error) {
	var resp string

	err := c.rpcClient.Call("Plugin.Name", new(any), &resp)
	if err != nil {
		return "", fmt.Errorf("calling plugin Name: %w", err)
	}

	return resp, nil
}

// Version calls the server-side Version method and returns its version.
func (c *Client) Version() (string, error) {
	var resp string

	err := c.rpcClient.Call("Plugin.Version", new(any), &resp)
	if err != nil {
		return "", fmt.Errorf("calling plugin Version: %w", err)
	}

	return resp, nil
}

// Execute calls the server-side Execute method and returns generated output.
func (c *Client) Execute(args *ExecuteArgs) (string, error) {
	var resp string

	err := c.rpcClient.Call("Plugin.Execute", args, &resp)
	if err != nil {
		return "", fmt.Errorf("calling plugin Execute: %w", err)
	}

	return resp, nil
}
