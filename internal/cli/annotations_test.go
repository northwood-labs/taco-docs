// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// WHY: Verifies cobra command annotations are correctly constructed. These annotations are used by
// the CLI framework to identify formatter subcommands and their associated command strings.
func TestCommandAnnotations(t *testing.T) {
	tests := []struct {
		name    string
		command string
	}{
		{
			name:    "command annotations",
			command: "foo",
		},
		{
			name:    "command annotations",
			command: "foo bar",
		},
		{
			name:    "command annotations",
			command: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual := Annotations(tt.command)
			assert.Equal(tt.command, actual["command"])
			assert.Equal("formatter", actual["kind"])
		})
	}
}
