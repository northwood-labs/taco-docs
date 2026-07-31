// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package template // lint:allow_naming_conflict_stdlib

import "errors"

var (
	// ErrBaseTemplateNotFound indicates no template items were registered.
	ErrBaseTemplateNotFound = errors.New("base template not found")

	// ErrTemplateNotFound indicates the requested named template does not
	// exist.
	ErrTemplateNotFound = errors.New("template not found")
)
