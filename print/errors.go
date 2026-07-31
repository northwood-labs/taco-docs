// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package print

import "errors"

var (
	// ErrRecursivePathEmpty indicates the --recursive-path flag value is empty
	// when recursive mode is enabled.
	ErrRecursivePathEmpty = errors.New("value of '--recursive-path' can't be empty")

	// ErrShowHideConflict indicates both --show and --hide flags were used
	// together.
	ErrShowHideConflict = errors.New("'--show' and '--hide' can't be used together")

	// ErrInvalidSection indicates a section name is not recognized.
	ErrInvalidSection = errors.New("not a valid section")

	// ErrOutputModeEmpty indicates the --output-mode flag value is empty.
	ErrOutputModeEmpty = errors.New("value of '--output-mode' can't be empty")

	// ErrOutputTemplateEmpty indicates the --output-template flag value is
	// empty.
	ErrOutputTemplateEmpty = errors.New("value of '--output-template' can't be empty")

	// ErrOutputTemplateNoContent indicates the output template does not contain
	// the required {{ .Content }} placeholder.
	ErrOutputTemplateNoContent = errors.New(
		"value of '--output-template' doesn't have '{{ .Content }}'" +
			" (note that spaces inside '{{ }}' are mandatory)",
	)

	// ErrOutputTemplateTooShort indicates the output template has fewer than
	// the minimum required lines.
	ErrOutputTemplateTooShort = errors.New(
		"value of '--output-template' should contain at least 3 lines" +
			" (begin comment, {{ .Content }}, and end comment)",
	)

	// ErrOutputTemplateNoBegin indicates the output template is missing the
	// begin comment marker.
	ErrOutputTemplateNoBegin = errors.New(
		"value of '--output-template' is missing begin comment",
	)

	// ErrOutputTemplateNoEnd indicates the output template is missing the end
	// comment marker.
	ErrOutputTemplateNoEnd = errors.New(
		"value of '--output-template' is missing end comment",
	)

	// ErrOutputValuesFromEmpty indicates --output-values-from is required but
	// was not provided.
	ErrOutputValuesFromEmpty = errors.New("value of '--output-values-from' is missing")

	// ErrInvalidSortType indicates the sort-by value is not a recognized sort
	// type.
	ErrInvalidSortType = errors.New("not a valid sort type")

	// ErrFormatterEmpty indicates the formatter value is empty.
	ErrFormatterEmpty = errors.New("value of 'formatter' can't be empty")

	// ErrHeaderFromEmpty indicates the --header-from flag value is empty.
	ErrHeaderFromEmpty = errors.New("value of '--header-from' can't be empty")

	// ErrFooterFromEmpty indicates the --footer-from flag value is empty.
	ErrFooterFromEmpty = errors.New("value of '--footer-from' can't be empty")

	// ErrFooterEqualsHeader indicates --footer-from and --header-from point to
	// the same file.
	ErrFooterEqualsHeader = errors.New(
		"value of '--footer-from' can't equal value of '--header-from",
	)

	// ErrConfigFileNotFound indicates the specified configuration file does not
	// exist.
	ErrConfigFileNotFound = errors.New("config file not found")
)
