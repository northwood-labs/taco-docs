// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

// Package reader provides utilities for extracting comment blocks and
// specific line ranges from Terraform source files.
package reader

import "errors"

var (
	// ErrNoLines indicates the file contains no lines to read.
	ErrNoLines = errors.New("no lines in file")

	// ErrOnlyOneLine indicates the file contains only one line, which is
	// insufficient to extract a comment block above a declaration.
	ErrOnlyOneLine = errors.New("only 1 line")

	// ErrInsufficientLines indicates the file does not contain enough lines
	// to reach the target line number.
	ErrInsufficientLines = errors.New("insufficient lines in file")
)
