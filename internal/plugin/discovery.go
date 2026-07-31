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
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-homedir"

	pluginsdk "github.com/northwood-labs/taco-docs/plugin"
)

// Discover scans well-known directories for plugin binaries and initializes RPC
// connections to each one. The priority order (env var > local > home) allows
// CI pipelines to override plugin locations via TFDOCS_PLUGIN_DIR while
// developers use the conventional local or home-based paths. Only the first
// matching directory is used — this avoids confusion from loading the same
// plugin from multiple locations.
func Discover() (*List, error) {
	if dir := os.Getenv("TFDOCS_PLUGIN_DIR"); dir != "" {
		plugins, err := findPlugins(dir)
		if err != nil {
			return nil, fmt.Errorf("finding plugins in TFDOCS_PLUGIN_DIR: %w", err)
		}

		return plugins, nil
	}

	if _, err := os.Stat(localPluginsRoot); !os.IsNotExist(err) {
		plugins, err := findPlugins(localPluginsRoot)
		if err != nil {
			return nil, fmt.Errorf("finding plugins in local root: %w", err)
		}

		return plugins, nil
	}

	dir, err := homedir.Expand(homePluginsRoot)
	if err != nil {
		return nil, fmt.Errorf("expanding home plugins path: %w", err)
	}

	plugins, err := findPlugins(dir)
	if err != nil {
		return nil, fmt.Errorf("finding plugins in home directory: %w", err)
	}

	return plugins, nil
}

// findPlugins iterates over all files in a directory, treats each one as a
// potential plugin binary (based on the naming convention), spawns it as a
// subprocess, and establishes an RPC connection. The go-plugin library handles
// the handshake to ensure version compatibility between host and plugin.
func findPlugins(dir string) (*List, error) {
	clients := make(map[string]*goplugin.Client)
	formatters := make(map[string]*pluginsdk.Client)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading plugin directory: %w", err)
	}

	for _, f := range files {
		// Strip the mandatory prefix to get the formatter name. For example,
		// "tfdocs-format-custom" becomes "custom" as the formatter identifier.
		name := strings.ReplaceAll(f.Name(), namePrefix, "")

		path, err := getPluginPath(dir, name)
		if err != nil {
			return nil, fmt.Errorf("resolving plugin path for %q: %w", name, err)
		}

		// Each plugin runs as a separate process communicating over RPC. This
		// isolation means a crashing plugin can't bring down the host, and
		// plugins can be written in any language that implements the protocol.
		cmd := exec.CommandContext(context.TODO(), path) // lint:allow_possible_insecure

		client := pluginsdk.NewClient(&pluginsdk.ClientOpts{
			Cmd: cmd,
		})

		rpcClient, err := client.Client()
		if err != nil {
			return nil, fmt.Errorf("creating RPC client for plugin %q: %w", name, err)
		}

		raw, err := rpcClient.Dispense("formatter")
		if err != nil {
			return nil, fmt.Errorf("dispensing formatter for plugin %q: %w", name, err)
		}

		formatter, ok := raw.(*pluginsdk.Client)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrPluginUnexpectedType, name)
		}

		// Duplicate plugin names would cause silent shadowing — fail loudly.
		if _, ok := clients[name]; ok {
			return nil, fmt.Errorf("%w: %s", ErrPluginAlreadyRegistered, name)
		}

		clients[name] = client
		formatters[name] = formatter
	}

	return &List{formatters: formatters, clients: clients}, nil
}

// getPluginPath constructs the expected filesystem path for a plugin binary,
// including the platform-appropriate executable suffix. It verifies the file
// exists before returning, producing a clear error for missing plugins.
func getPluginPath(dir, name string) (string, error) {
	suffix := ""

	if runtime.GOOS == "windows" {
		suffix += ".exe"
	}

	path := filepath.Join(dir, fmt.Sprintf("%s%s%s", namePrefix, name, suffix))

	if _, err := os.Stat(path); os.IsNotExist(err) { // lint:allow_possible_insecure
		return "", os.ErrNotExist
	}

	return path, nil
}
