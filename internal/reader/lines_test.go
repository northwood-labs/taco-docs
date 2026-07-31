// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package reader

import (
	"strings"
	"testing"

	assertpkg "github.com/go-openapi/testify/assert"
)

const (
	textEmpty = ``

	textOneLine = `Lorem ipsum dolor sit amet`

	textWithLeadingComment = `
/**
 * Morbi vitae nulla in dui lobortis
 * consectetur. Integer nec tempus
 * felis. Ut quis suscipit risus.
 * Donec lobortis consequat nunc, in
 * efficitur mi maximus ac. Sed id
 * felis posuere, aliquam purus eget,
 * faucibus augue.
 */

Lorem ipsum dolor sit amet,
consectetur adipiscing elit,

# sed do eiusmod tempor incididunt
ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis

# nostrud exercitation ullamco
laboris nisi ut aliquip ex ea
commodo consequat. Duis aute irure

# dolor in reprehenderit in voluptate
velit esse cillum dolore eu fugiat
nulla pariatur.

# Excepteur sint occaecat cupidatat
# non proident, sunt in culpa qui
officia deserunt mollit anim id est laborum
`

	textWithoutLeadingComment = `
Lorem ipsum dolor sit amet,
consectetur adipiscing elit,

# sed do eiusmod tempor incididunt
ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis

# nostrud exercitation ullamco
laboris nisi ut aliquip ex ea
commodo consequat. Duis aute irure

# dolor in reprehenderit in voluptate
velit esse cillum dolore eu fugiat
nulla pariatur.

# Excepteur sint occaecat cupidatat
# non proident, sunt in culpa qui
officia deserunt mollit anim id est laborum
`
)

// Validates comment extraction from files and text by line number. This is the
// core of how descriptions are read from preceding comments in .tf files for
// inputs/outputs without description attrs.
func TestReadLinesFromFile(t *testing.T) {
	tests := []struct {
		name       string
		fileName   string
		expected   string
		lineNumber int
		wantError  bool
	}{
		{
			name:       "extract lines from file",
			fileName:   "testdata/sample.txt",
			lineNumber: -1,
			expected: "Morbi vitae nulla in dui lobortis consectetur. Integer nec tempus felis. Ut quis suscipit " +
				"risus. Donec lobortis consequat nunc, in efficitur mi maximus ac. Sed id felis posuere, aliquam " +
				"purus eget, faucibus augue.", // lint:no_const
			wantError: false,
		},
		{
			name:       "extract lines from file",
			fileName:   "testdata/sample.txt",
			lineNumber: 15,
			expected:   "sed do eiusmod tempor incididunt",
			wantError:  false,
		},
		{
			name:       "extract lines from file",
			fileName:   "testdata/no-traling-line.txt",
			lineNumber: -1,
			expected: "Morbi vitae nulla in dui lobortis consectetur. Integer nec tempus felis. Ut quis suscipit " +
				"risus. Donec lobortis consequat nunc, in efficitur mi maximus ac. Sed id felis posuere, aliquam " +
				"purus eget, faucibus augue.", // lint:no_const
			wantError: false,
		},
		{
			name:       "extract lines from file",
			fileName:   "testdata/empty.txt",
			lineNumber: 0,
			expected:   "",
			wantError:  false,
		},
		{
			name:       "extract lines from file",
			fileName:   "testdata/noop.txt",
			lineNumber: -1,
			expected:   "",
			wantError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)
			lines := Lines{
				FileName: tt.fileName,
				LineNum:  tt.lineNumber,
				Condition: func(line string) bool {
					line = strings.TrimSpace(line)

					return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/*") ||
						strings.HasPrefix(line, "*") ||
						strings.HasPrefix(line, "*/")
				},
				Parser: func(line string) (string, bool) {
					line = strings.TrimSpace(line)
					if strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*/") {
						return "", false
					}

					if line == "*" {
						return "", true
					}

					line = strings.TrimPrefix(line, "* ")
					line = strings.TrimPrefix(line, "#")
					line = strings.TrimSpace(line)

					return line, true
				},
			}

			comment, err := lines.Extract()
			if tt.wantError {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tt.expected, strings.Join(comment, " "))
			}
		})
	}
}

// Same as file-based extraction but from in-memory text. Tests edge cases like
// empty content, out-of-range line numbers, and text without leading comments.
func TestReadLinesFromText(t *testing.T) {
	tests := []struct {
		name        string
		textContent string
		expected    string
		errorText   string
		lineNumber  int
		wantError   bool
	}{
		{
			name:        "extract lines from text",
			textContent: textWithLeadingComment,
			lineNumber:  -1,
			expected: "Morbi vitae nulla in dui lobortis consectetur. Integer nec tempus felis. Ut quis suscipit " +
				"risus. Donec lobortis consequat nunc, in efficitur mi maximus ac. Sed id felis posuere, aliquam " +
				"purus eget, faucibus augue.", // lint:no_const
			wantError: false,
			errorText: "",
		},
		{
			name:        "extract lines from text",
			textContent: textEmpty,
			lineNumber:  -1,
			expected:    "",
			wantError:   true,
			errorText:   "no lines in file",
		},
		{
			name:        "extract lines from text",
			textContent: textOneLine,
			lineNumber:  10,
			expected:    "",
			wantError:   true,
			errorText:   "only 1 line",
		},
		{
			name:        "extract lines from text",
			textContent: textWithoutLeadingComment,
			lineNumber:  -1,
			expected:    "",
			wantError:   false,
			errorText:   "",
		},
		{
			name:        "extract lines from text",
			textContent: textWithLeadingComment,
			lineNumber:  15,
			expected:    "sed do eiusmod tempor incididunt",
			wantError:   false,
			errorText:   "",
		},
		{
			name:        "extract lines from text",
			textContent: textWithLeadingComment,
			lineNumber:  19,
			expected:    "nostrud exercitation ullamco",
			wantError:   false,
			errorText:   "",
		},
		{
			name:        "extract lines from text",
			textContent: textWithLeadingComment,
			lineNumber:  23,
			expected:    "dolor in reprehenderit in voluptate",
			wantError:   false,
			errorText:   "",
		},
		{
			name:        "extract lines from text",
			textContent: textWithLeadingComment,
			lineNumber:  28,
			expected:    "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui",
			wantError:   false,
			errorText:   "",
		},
		{
			name:        "extract lines from text",
			textContent: textWithLeadingComment,
			lineNumber:  54,
			expected:    "",
			wantError:   true,
			errorText:   "only 28 lines",
		},
		{
			name:        "extract lines from text",
			textContent: textEmpty,
			lineNumber:  10,
			expected:    "",
			wantError:   true,
			errorText:   "no lines in file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assertpkg.New(t)
			lines := Lines{
				LineNum: tt.lineNumber,
				Condition: func(line string) bool {
					line = strings.TrimSpace(line)

					return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/*") ||
						strings.HasPrefix(line, "*") ||
						strings.HasPrefix(line, "*/")
				},
				Parser: func(line string) (string, bool) {
					line = strings.TrimSpace(line)
					if strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*/") {
						return "", false
					}

					if line == "*" {
						return "", true
					}

					line = strings.TrimPrefix(line, "* ")
					line = strings.TrimPrefix(line, "#")
					line = strings.TrimSpace(line)

					return line, true
				},
			}
			r := strings.NewReader(strings.TrimSpace(tt.textContent))

			comment, err := lines.extractLines(r)
			if tt.wantError {
				assert.Error(err)
				assert.Equal(tt.errorText, err.Error())
			} else {
				assert.NoError(err)
				assert.Equal(tt.expected, strings.Join(comment, " "))
			}
		})
	}
}
