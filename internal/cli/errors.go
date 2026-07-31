// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package cli

import "errors"

var (
	// ErrConfigEmpty indicates the --config flag was provided with an empty
	// value.
	ErrConfigEmpty = errors.New("value of '--config' can't be empty")

	// ErrOutputFileEmpty indicates --output-file is required but was not
	// provided.
	ErrOutputFileEmpty = errors.New("value of '--output-file' cannot be empty with '--recursive'")

	// ErrConfigFileNotFound indicates the specified configuration file does
	// not exist.
	ErrConfigFileNotFound = errors.New("config file not found")

	// ErrFormatterNotFound indicates the requested formatter name does not
	// match
	// any registered built-in or plugin formatter.
	ErrFormatterNotFound = errors.New("formatter not found")

	// ErrVersionConstraint indicates the current version does not satisfy
	// the
	// version constraint specified in the configuration file.
	ErrVersionConstraint = errors.New("version constraint not satisfied")

	// ErrTemplateMissing indicates the output template is required but was
	// not provided.
	ErrTemplateMissing = errors.New("template is missing")

	// ErrBeginCommentMissing indicates the begin comment marker is missing
	// from the target file.
	ErrBeginCommentMissing = errors.New("begin comment is missing")

	// ErrEndCommentMissing indicates the end comment marker is missing from
	// the target file.
	ErrEndCommentMissing = errors.New("end comment is missing")

	// ErrEndBeforeBegin indicates the end comment marker appears before the
	// begin comment marker.
	ErrEndBeforeBegin = errors.New("end comment is before begin comment")

	// ErrFileOutOfDate indicates the output file content does not match the
	// generated content.
	ErrFileOutOfDate = errors.New("file is out of date")
)
