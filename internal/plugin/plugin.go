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

	pluginsdk "github.com/terraform-docs/terraform-docs/plugin"
)

// namePrefix is the mandatory prefix for plugin binary names. This naming
// convention (tfdocs-format-<name>) serves as a discovery mechanism — the
// plugin system can find plugins by scanning a directory for files matching
// this prefix without requiring a registry or manifest file.
const namePrefix = "tfdocs-format-"

// homePluginsRoot and localPluginsRoot define the default search paths for
// plugin binaries. The local path (./.tfdocs.d/plugins) takes priority over
// the home path (~/.tfdocs.d/plugins), allowing project-specific plugins to
// shadow globally installed ones.
var (
	homePluginsRoot  = "~/.tfdocs.d/plugins"
	localPluginsRoot = "./.tfdocs.d/plugins"
)

// List caches discovered plugins and their corresponding go-plugin clients.
// Caching prevents repeated filesystem scans and process spawning during a
// single terraform-docs invocation. The struct wraps both the high-level
// formatter clients (used to call Execute) and the raw go-plugin clients
// (needed for process lifecycle management via Kill()).
type List struct {
	formatters map[string]*pluginsdk.Client
	clients    map[string]*goplugin.Client
}

// All returns every registered plugin client. This is used by the version
// command to enumerate installed plugins for display.
func (l *List) All() []*pluginsdk.Client {
	all := make([]*pluginsdk.Client, 0)
	for _, f := range l.formatters {
		all = append(all, f)
	}
	return all
}

// Get retrieves a plugin by its formatter name. The boolean return follows
// Go map-access convention to let callers distinguish "not found" from a
// nil value.
func (l *List) Get(name string) (*pluginsdk.Client, bool) {
	client, ok := l.formatters[name]
	return client, ok
}

// Clean terminates all plugin subprocesses. This should be called during
// shutdown to avoid leaving orphan processes — go-plugin uses os/exec
// under the hood and processes won't terminate automatically when the
// parent exits.
func (l *List) Clean() {
	for _, client := range l.clients {
		client.Kill()
	}
}
