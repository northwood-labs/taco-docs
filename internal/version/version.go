// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package version

import (
	"fmt"
	"runtime"
)

// current version — these are the source-of-truth for the release version.
// The prerelease field is set to a non-empty string (e.g., "rc1") during
// pre-release builds and cleared for stable releases. This is the only
// place version numbers are defined; ldflags injects the commit hash at
// build time.
const (
	coreVersion = "0.24.0"
	prerelease  = ""
)

// commit is provisioned by ldflags at build time (-ldflags "-X ...version.commit=abc123").
// This allows users and bug reports to identify the exact source revision.
var commit string

// Core returns the bare semantic version without pre-release suffix or metadata.
// This is used for version constraint checking against config files — constraints
// should match against the stable version number, not the full string.
func Core() string {
	return coreVersion
}

// Short returns the version with pre-release suffix (if any). This is appropriate
// for display in contexts where the commit and platform aren't needed.
func Short() string {
	v := coreVersion

	if prerelease != "" {
		v += "-" + prerelease
	}

	return v
}

// Full returns the complete version string including pre-release, commit hash,
// and OS/architecture. This goes into --version output and the User-Agent header,
// giving complete provenance for debugging environment-specific issues.
func Full() string {
	if commit != "" && commit[:1] != " " {
		commit = " " + commit
	}

	return fmt.Sprintf("v%s%s %s/%s", Short(), commit, runtime.GOOS, runtime.GOARCH)
}
