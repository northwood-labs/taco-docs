// Copyright 2021 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Lines represents a conditional line reader that extracts content from a file
// based on caller-defined rules. It's the engine behind two key behaviors:
//
//  1. Reading comments above Terraform declarations (LineNum is set): it reads
//     lines immediately preceding the target line, capturing comment blocks that
//     serve as documentation descriptions.
//
//  2. Reading header/footer from .tf files (LineNum is -1): it scans the whole
//     file looking for the first block-comment, which by convention contains
//     the module description.
//
// The Condition function determines which lines are "interesting" (e.g., comment
// lines), and Parser transforms each matching line into its final form (e.g.,
// stripping comment prefixes). This separation of concerns lets the same reader
// handle multiple comment styles (// and # for inline, /* */ for block).
type Lines struct {
	FileName  string
	LineNum   int // value -1 means scan the whole file and break after finding what we were looking for
	Condition func(line string) bool
	Parser    func(line string) (string, bool)
}

// Extract opens the file and delegates to the internal extraction logic.
// It handles file-level concerns (open, stat, close) separately from the
// line-by-line parsing, keeping the parsing logic testable with io.Reader.
func (l *Lines) Extract() ([]string, error) {
	f, err := os.Open(l.FileName)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	// Empty files are a valid case (e.g., placeholder .tf files) — return
	// empty result rather than erroring.
	if stat.Size() == 0 {
		return []string{}, nil
	}
	defer func() {
		_ = f.Close()
	}()
	return l.extract(f)
}

// extract performs the actual line-by-line reading. The algorithm works by
// accumulating lines that satisfy Condition and resetting when a non-matching
// line is encountered. This naturally captures the *last contiguous block* of
// matching lines before LineNum — which is exactly the comment block immediately
// above a declaration.
//
// When LineNum is -1, it scans the entire file but stops at the first break
// after finding matching lines. This handles the "first block comment in file"
// case used for module headers.
func (l *Lines) extract(r io.Reader) ([]string, error) { //nolint:gocyclo
	// NOTE(khos2ow): this function is over our cyclomatic complexity goal.
	// Be wary when adding branches, and look for functionality that could
	// be reasonably moved into an injected dependency.

	bf := bufio.NewReader(r)
	lines := make([]string, 0)
	for lnum := 0; ; lnum++ {
		// Stop once we've reached the target line — any lines accumulated at
		// this point are the comment block immediately above the declaration.
		if l.LineNum != -1 && lnum >= l.LineNum-1 {
			break
		}
		line, err := bf.ReadString('\n')
		if errors.Is(err, io.EOF) && line == "" {
			switch lnum {
			case 0:
				return nil, errors.New("no lines in file")
			case 1:
				return nil, errors.New("only 1 line")
			default:
				if l.LineNum == -1 {
					break
				}
				return nil, fmt.Errorf("only %d lines", lnum)
			}
		}

		//nolint:gocritic
		if l.Condition(line) {
			if extracted, capture := l.Parser(line); capture {
				lines = append(lines, extracted)
			}
		} else if l.LineNum == -1 {
			// In whole-file scan mode, the first non-matching line after finding
			// matches means we've captured the entire block — stop scanning.
			break
		} else {
			// In targeted mode, a non-matching line resets the accumulator because
			// we only want the contiguous block immediately before the target line.
			lines = nil
		}
	}
	return lines, nil
}
