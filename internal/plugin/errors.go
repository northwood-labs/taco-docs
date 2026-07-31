// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package plugin // lint:allow_naming_conflict_stdlib

import "errors"

var (
	// ErrPluginAlreadyRegistered indicates a plugin with the same name has
	// already been registered.
	ErrPluginAlreadyRegistered = errors.New("plugin is already registered")

	// ErrPluginUnexpectedType indicates the plugin returned an unexpected type
	// instead of the expected formatter client.
	ErrPluginUnexpectedType = errors.New("unexpected formatter type")
)
