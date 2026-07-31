// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package terraform

import "errors"

var (
	// ErrSectionMissing indicates a required section name was not provided.
	ErrSectionMissing = errors.New("section is missing")

	// ErrSectionFromMissing indicates the --header-from or --footer-from value
	// is missing.
	ErrSectionFromMissing = errors.New("section source file value is missing")

	// ErrUnsupportedFileFormat indicates the file format is not supported for
	// reading section content.
	ErrUnsupportedFileFormat = errors.New(
		"only .adoc, .md, .tf, .tofu and .txt formats are supported",
	)
)
