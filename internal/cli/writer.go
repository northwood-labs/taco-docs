/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/terraform-docs/terraform-docs/print"
)

// stdoutWriter writes content to os.Stdout. It appends a trailing newline
// because documentation output should end cleanly when piped or displayed
// in a terminal — without it, the shell prompt would appear on the same line.
type stdoutWriter struct{}

// Write content to Stdout
func (sw *stdoutWriter) Write(p []byte) (int, error) {
	return os.Stdout.WriteString(string(p) + "\n")
}

// fileWriter handles writing generated documentation to a file on disk. It
// supports two operational modes that address different workflows:
//
//   - "replace" mode: overwrites the entire file with generated content. This is
//     the simpler model, suitable when the file is fully machine-generated.
//
//   - "inject" mode: inserts generated content between begin/end comment markers
//     in an existing file. This enables the common pattern where a README has
//     hand-written sections (project overview, examples) alongside an auto-generated
//     section for inputs/outputs. The markers (e.g., <!-- BEGIN_TF_DOCS -->) allow
//     the tool to update only its section without disturbing the rest.
type fileWriter struct {
	file string
	dir  string

	mode string

	check bool

	template string
	begin    string
	end      string

	writer io.Writer
}

// Write content to target file. The logic branches on output mode:
//   - For "replace": apply template (if any) then write the whole file.
//   - For "inject": apply template then splice the result between markers in the
//     existing file content, preserving everything outside the markers.
func (fw *fileWriter) Write(p []byte) (int, error) {
	filename := fw.fullFilePath()

	if fw.template == "" {
		// template is optional for mode replace — content is written as-is.
		if fw.mode == print.OutputModeReplace {
			return fw.write(filename, p)
		}
		return 0, errors.New("template is missing")
	}

	// Wrap the raw content in the user's output template (which typically adds
	// the begin/end comment markers around the generated documentation).
	buf, err := fw.apply(p)
	if err != nil {
		return 0, err
	}

	// In replace mode, the entire file becomes the templated output.
	if fw.mode == print.OutputModeReplace {
		return fw.write(filename, buf.Bytes())
	}

	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		// In inject mode, if the target file doesn't exist yet, create it with
		// just the generated content — this bootstraps the initial file.
		return fw.write(filename, buf.Bytes())
	}

	if len(content) == 0 {
		// An empty target file is treated the same as a missing one for inject mode.
		return fw.write(filename, buf.Bytes())
	}

	return fw.inject(filename, string(content), buf.String())
}

// fullFilePath resolves the output file path. If the configured path is absolute
// it's used directly; otherwise it's joined with the module root directory. This
// supports both project-relative paths (common) and absolute paths (rare but needed
// for cross-project documentation aggregation).
func (fw *fileWriter) fullFilePath() string {
	if filepath.IsAbs(fw.file) {
		return fw.file
	}
	return filepath.Join(fw.dir, fw.file)
}

// apply wraps the generated content in the user's output template. The template
// typically contains the begin/end comment markers with {{ .Content }} between them.
func (fw *fileWriter) apply(p []byte) (bytes.Buffer, error) {
	type content struct {
		Content string
	}

	var buf bytes.Buffer

	tmpl := template.Must(template.New("content").Parse(fw.template))
	err := tmpl.ExecuteTemplate(&buf, "content", content{string(p)})

	return buf, err
}

// inject splices generated content into an existing file between the begin and
// end comment markers. This preserves any hand-written content above and below
// the markers. The function validates marker presence and ordering to prevent
// silent data corruption from malformed files.
func (fw *fileWriter) inject(filename string, content string, generated string) (int, error) {
	before := strings.Index(content, fw.begin)
	after := strings.Index(content, fw.end)

	// If neither marker is present, append the generated content to the existing
	// file — this handles the first-time injection case for files that don't yet
	// have markers but already have content.
	if before < 0 && after < 0 {
		return fw.write(filename, []byte(content+"\n"+generated))
	}

	if before < 0 {
		return 0, errors.New("begin comment is missing")
	}

	generated = content[:before] + generated

	if after < 0 {
		return 0, errors.New("end comment is missing")
	}

	if after < before {
		return 0, errors.New("end comment is before begin comment")
	}

	// Preserve everything after the end marker (including the marker itself is
	// consumed, and content after it is re-appended).
	generated += content[after+len(fw.end):]

	return fw.write(filename, []byte(generated))
}

// write persists content to disk (or to an injected io.Writer for testing). In
// "check" mode it performs a diff instead of writing — this enables CI pipelines
// to verify that generated docs are up-to-date without actually modifying files,
// failing the build if changes are detected.
func (fw *fileWriter) write(filename string, p []byte) (int, error) {
	// Check mode: compare against existing file content and report staleness
	// without modifying anything. This supports CI "lint" workflows.
	if fw.check {
		f, err := os.ReadFile(filepath.Clean(filename))
		if err != nil {
			return 0, err
		}

		if !bytes.Equal(f, p) {
			return 0, fmt.Errorf("%s is out of date", filename)
		}

		fmt.Printf("%s is up to date\n", filename)
		return 0, nil
	}

	// If an io.Writer was injected (for testing), use it instead of the filesystem.
	if fw.writer != nil {
		return fw.writer.Write(p)
	}

	err := os.WriteFile(filename, p, 0o644)
	if err == nil {
		fmt.Printf("%s updated successfully\n", filename)
	}
	return len(p), err
}
