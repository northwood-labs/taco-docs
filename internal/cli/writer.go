// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/northwood-labs/taco-docs/print"
)

type (
	// stdoutWriter writes content to [os.Stdout]. It appends a trailing newline
	// because documentation output should end cleanly when piped or displayed
	// in a terminal — without it, the shell prompt would appear on the same
	// line.
	stdoutWriter struct{}

	// fileWriter handles writing generated documentation to a file on disk. It
	// supports two operational modes that address different workflows:
	//
	// - "replace" mode: overwrites the entire file with generated content. This
	// is
	//
	//  the simpler model, suitable when the file is fully machine-generated.
	//
	// - "inject" mode: inserts generated content between begin/end comment
	// markers in an existing file. This enables the common pattern where a
	// README has hand-written sections (project overview, examples) alongside
	// an auto-generated section for inputs/outputs. The markers (e.g., <!--
	// BEGIN_TF_DOCS -->) allow
	//
	//  the tool to update only its section without disturbing the rest.
	fileWriter struct {
		writer   io.Writer
		file     string
		dir      string
		mode     string
		template string
		begin    string
		end      string
		check    bool
	}
)

// Write content to Stdout.
func (*stdoutWriter) Write(p []byte) (int, error) {
	_, err := os.Stdout.WriteString(string(p) + "\n")
	if err != nil {
		return 0, fmt.Errorf("writing to stdout: %w", err)
	}

	return len(p), nil
}

// Write content to target file. The logic branches on output mode:
//
// - For "replace": apply template (if any) then write the whole file.
//
// - For "inject": apply template then splice the result between markers in the
// existing file content, preserving everything outside the markers.
func (fw *fileWriter) Write(p []byte) (int, error) {
	filename := fw.fullFilePath()

	if fw.template == "" {
		// template is optional for mode replace — content is written as-is.
		if fw.mode == print.OutputModeReplace {
			_, err := fw.writeToFile(filename, p)
			if err != nil {
				return 0, fmt.Errorf("writing file: %w", err)
			}

			return len(p), nil
		}

		return 0, ErrTemplateMissing
	}

	// Wrap the raw content in the user's output template (which typically adds
	// the begin/end comment markers around the generated documentation).
	buf, err := fw.apply(p)
	if err != nil {
		return 0, fmt.Errorf("applying output template: %w", err)
	}

	// In replace mode, the entire file becomes the templated output.
	if fw.mode == print.OutputModeReplace {
		_, replaceErr := fw.writeToFile(filename, buf.Bytes())
		if replaceErr != nil {
			return 0, fmt.Errorf("writing file: %w", replaceErr)
		}

		return len(p), nil
	}

	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		// In inject mode, if the target file doesn't exist yet, create it with
		// just the generated content — this bootstraps the initial file.
		_, writeErr := fw.writeToFile(filename, buf.Bytes())
		if writeErr != nil {
			return 0, fmt.Errorf("writing new file: %w", writeErr)
		}

		return len(p), nil
	}

	if len(content) == 0 {
		// An empty target file is treated the same as a missing one for inject
		// mode.
		_, emptyErr := fw.writeToFile(filename, buf.Bytes())
		if emptyErr != nil {
			return 0, fmt.Errorf("writing to empty file: %w", emptyErr)
		}

		return len(p), nil
	}

	n, err := fw.inject(filename, string(content), buf.String())
	if err != nil {
		return n, fmt.Errorf("injecting content: %w", err)
	}

	return n, nil
}

// fullFilePath resolves the output file path. If the configured path is
// absolute it's used directly; otherwise it's joined with the module root
// directory. This supports both project-relative paths (common) and absolute
// paths (rare but needed for cross-project documentation aggregation).
func (fw *fileWriter) fullFilePath() string {
	if filepath.IsAbs(fw.file) {
		return fw.file
	}

	return filepath.Join(fw.dir, fw.file)
}

// apply wraps the generated content in the user's output template. The template
// typically contains the begin/end comment markers with {{ .Content }} between
// them.
func (fw *fileWriter) apply(p []byte) (bytes.Buffer, error) {
	type content struct {
		Content string
	}

	var buf bytes.Buffer

	tmpl := template.Must(template.New("content").Parse(fw.template))

	err := tmpl.ExecuteTemplate(&buf, "content", content{string(p)})
	if err != nil {
		return buf, fmt.Errorf("executing output template: %w", err)
	}

	return buf, nil
}

// inject splices generated content into an existing file between the begin and
// end comment markers. This preserves any hand-written content above and below
// the markers. The function validates marker presence and ordering to prevent
// silent data corruption from malformed files.
func (fw *fileWriter) inject(filename, content, generated string) (int, error) {
	before := strings.Index(content, fw.begin)
	after := strings.Index(content, fw.end)

	// If neither marker is present, append the generated content to the
	// existing file — this handles the first-time injection case for files that
	// don't yet have markers but already have content.
	if before < 0 && after < 0 {
		n, err := fw.writeToFile(filename, []byte(content+"\n"+generated))
		if err != nil {
			return n, fmt.Errorf("writing appended content: %w", err)
		}

		return n, nil
	}

	if before < 0 {
		return 0, ErrBeginCommentMissing
	}

	generated = content[:before] + generated

	if after < 0 {
		return 0, ErrEndCommentMissing
	}

	if after < before {
		return 0, ErrEndBeforeBegin
	}

	// Preserve everything after the end marker (including the marker itself is
	// consumed, and content after it is re-appended).
	generated += content[after+len(fw.end):]

	n, err := fw.writeToFile(filename, []byte(generated))
	if err != nil {
		return n, fmt.Errorf("writing injected content: %w", err)
	}

	return n, nil
}

// writeToFile persists content to disk (or to an injected [io.Writer] for
// testing).
//
// In "check" mode it performs a diff instead of writing — this enables CI
// pipelines to verify that generated docs are up-to-date without actually
// modifying files, failing the build if changes are detected.
func (fw *fileWriter) writeToFile(filename string, p []byte) (int, error) {
	// Check mode: compare against existing file content and report staleness
	// without modifying anything. This supports CI "lint" workflows.
	if fw.check {
		f, err := os.ReadFile(filepath.Clean(filename))
		if err != nil {
			return 0, fmt.Errorf("reading file for check: %w", err)
		}

		if !bytes.Equal(f, p) {
			return 0, fmt.Errorf("%w: %s", ErrFileOutOfDate, filename)
		}

		fmt.Printf("%s is up to date\n", filename)

		return 0, nil
	}

	// If an io.Writer was injected (for testing), use it instead of the
	// filesystem.
	if fw.writer != nil {
		n, err := fw.writer.Write(p)
		if err != nil {
			return n, fmt.Errorf("writing to injected writer: %w", err)
		}

		return n, nil
	}

	err := os.WriteFile(filename, p, 0o644) // lint:allow_possible_insecure
	if err == nil {
		fmt.Printf("%s updated successfully\n", filename)
	}

	return len(p), err
}
