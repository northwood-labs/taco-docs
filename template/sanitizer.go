/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"mvdan.cc/xurls/v2"
)

// SanitizeName escapes underscore character which have special meaning in
// Markdown.
func SanitizeName(name string, escape bool) string {
	if escape {
		// Escape underscore
		name = strings.ReplaceAll(name, "_", "\\_")
	}
	return name
}

// SanitizeSection converts passed 'string' to suitable Markdown or AsciiDoc
// representation for a document. (including line-break, illegal characters,
// code blocks etc). This is used for header and footer content where users
// control the source text and line-endings must be preserved exactly as written.
//
// IMPORTANT: SanitizeSection will never change the line-endings and preserve
// them as they are provided by the users.
func SanitizeSection(s string, escape bool, html bool) string {
	if s == "" {
		return "n/a"
	}
	result := processSegments(
		s,
		"```",
		func(segment string, first bool, last bool) string {
			segment = EscapeCharacters(segment, escape, false)
			segment = ConvertMultiLineText(segment, false, true, html)
			segment = NormalizeURLs(segment, escape)
			return segment
		},
		func(segment string, first bool, last bool) string {
			lastbreak := ""
			if !strings.HasSuffix(segment, "\n") {
				lastbreak = "\n"
			}

			// Adjust indentation and linebreak for indented codeblock
			// https://github.com/terraform-docs/terraform-docs/issues/521
			lastindent := ""
			lines := strings.Split(segment, "\n")
			if len(strings.TrimSpace(lines[len(lines)-1])) == 0 {
				lastbreak = ""
			}

			segment = fmt.Sprintf("```%s%s%s```", segment, lastindent, lastbreak)
			return segment
		},
	)
	return result
}

// SanitizeDocument converts passed 'string' to suitable Markdown or AsciiDoc
// representation for a document (including line-break, illegal characters,
// code blocks etc). Unlike SanitizeSection, this applies Markdown line-break
// conversion to produce proper multi-line paragraph rendering.
func SanitizeDocument(s string, escape bool, html bool) string {
	if s == "" {
		return "n/a"
	}
	result := processSegments(
		s,
		"```",
		func(segment string, first bool, last bool) string {
			segment = EscapeCharacters(segment, escape, false)
			segment = ConvertMultiLineText(segment, false, false, html)
			segment = NormalizeURLs(segment, escape)
			return segment
		},
		func(segment string, first bool, last bool) string {
			lastbreak := ""
			if !strings.HasSuffix(segment, "\n") {
				lastbreak = "\n"
			}
			segment = fmt.Sprintf("```%s%s```", segment, lastbreak)
			return segment
		},
	)
	return result
}

// SanitizeMarkdownTable converts passed 'string' to suitable Markdown
// representation for a table cell. Table cells can't contain literal newlines
// (they'd break the table structure), so line breaks are converted to <br/> or
// spaces depending on HTML mode.
func SanitizeMarkdownTable(s string, escape bool, html bool) string {
	if s == "" {
		return "n/a"
	}
	result := processSegments(
		s,
		"```",
		func(segment string, first bool, last bool) string {
			segment = EscapeCharacters(segment, escape, true)
			segment = ConvertMultiLineText(segment, true, false, html)
			segment = NormalizeURLs(segment, escape)
			return segment
		},
		func(segment string, first bool, last bool) string {
			linebreak := "<br/>"
			codestart := "<pre>"
			codeend := "</pre>"

			segment = strings.TrimSpace(segment)

			if !html {
				linebreak = ""
				codestart = " ```"
				codeend = "``` "

				if first {
					codestart = codestart[1:]
				}
				if last {
					codeend = codeend[:3]
				}

				segment = ConvertOneLineCodeBlock(segment)
			}

			segment = strings.ReplaceAll(segment, "\n", linebreak)
			segment = strings.ReplaceAll(segment, "\r", "")
			segment = fmt.Sprintf("%s%s%s", codestart, segment, codeend)
			return segment
		},
	)
	return result
}

// SanitizeAsciidocTable converts passed 'string' to suitable AsciiDoc
// representation for a table cell. AsciiDoc tables use different code block
// syntax ([source]/----) and don't need HTML line-break conversion.
func SanitizeAsciidocTable(s string, escape bool, html bool) string {
	if s == "" {
		return "n/a"
	}
	result := processSegments(
		s,
		"```",
		func(segment string, first bool, last bool) string {
			segment = EscapeCharacters(segment, escape, true)
			segment = NormalizeURLs(segment, escape)
			return segment
		},
		func(segment string, first bool, last bool) string {
			segment = strings.TrimSpace(segment)
			segment = fmt.Sprintf("[source]\n----\n%s\n----", segment)
			return segment
		},
	)
	return result
}

// ConvertMultiLineText translates natural line breaks into Markdown-compatible
// line breaks (double-space or <br/>) depending on whether we're in a table cell.
// Markdown requires trailing double-spaces for soft line breaks within a paragraph;
// tables need explicit <br/> tags since whitespace is collapsed in cells.
func ConvertMultiLineText(s string, isTable bool, isHeader bool, showHTML bool) string {
	if isTable {
		s = strings.TrimSpace(s)
	}

	// Convert line-break on a non-empty line followed by another line
	// starting with "alphanumeric" word into space-space-newline
	// which is a know convention of Markdown for multi-lines paragprah.
	// This doesn't apply on a markdown list for example, because all the
	// consecutive lines start with hyphen which is a special character.
	if !isHeader {
		s = regexp.MustCompile(`(\S*)(\r?\n)(\s*)(\w+)`).ReplaceAllString(s, "$1  $2$3$4")
		s = strings.ReplaceAll(s, "    \n", "  \n")
		s = strings.ReplaceAll(s, "  \n\n", "\n\n")
		s = strings.ReplaceAll(s, "\n  \n", "\n\n")
	}

	if !isTable {
		return s
	}

	// representation of line break. <br/> if showHTML is true, <space> if false.
	linebreak := " "

	if showHTML {
		linebreak = "<br/>"
	}

	// Convert space-space-newline to 'linebreak'.
	s = strings.ReplaceAll(s, "  \n", linebreak)

	// Convert single newline to 'linebreak'.
	return strings.ReplaceAll(s, "\n", linebreak)
}

// ConvertOneLineCodeBlock converts a multi-line code block into a one-liner.
// Line breaks are replaced with single space.
func ConvertOneLineCodeBlock(s string) string {
	split := strings.Split(s, "\n")
	result := []string{}
	for _, segment := range split {
		if len(strings.TrimSpace(segment)) == 0 {
			continue
		}
		segment = regexp.MustCompile(`(\s*)=(\s*)`).ReplaceAllString(segment, " = ")
		segment = strings.TrimLeftFunc(segment, unicode.IsSpace)
		result = append(result, segment)
	}
	return strings.Join(result, " ")
}

// EscapeCharacters prevents unintended Markdown formatting (e.g., underscores
// in variable names becoming italics) while being careful not to escape inside
// inline code spans. The processSegments split on backticks ensures code spans
// are left untouched.
func EscapeCharacters(s string, escape bool, escapePipe bool) string {
	// Escape pipe (only for 'markdown table' or 'asciidoc table')
	if escapePipe {
		s = processSegments(
			s,
			"`",
			func(segment string, first bool, last bool) string {
				return strings.ReplaceAll(segment, "|", "\\|")
			},
			func(segment string, first bool, last bool) string {
				return fmt.Sprintf("`%s`", segment)
			},
		)
	}

	if escape {
		s = processSegments(
			s,
			"`",
			func(segment string, first bool, last bool) string {
				return executePerLine(segment, func(line string) string {
					escape := func(char string) {
						c := strings.ReplaceAll(char, "*", "\\*")
						cases := []struct {
							pattern string
							index   []int
						}{
							{
								pattern: `^(\s*)(` + c + `+)(\s+)(.*)`,
								index:   []int{2},
							},
							{
								pattern: `(\s+)(` + c + `+)([^\t\n\f\r ` + c + `])(.*)([^\t\n\f\r ` + c + `])(` + c + `+)(\s+)`,
								index:   []int{6, 2},
							},
						}
						for i := range cases {
							c := cases[i]
							r := regexp.MustCompile(c.pattern)
							m := r.FindAllStringSubmatch(line, -1)
							i := r.FindAllStringSubmatchIndex(line, -1)
							for j := range m {
								for _, k := range c.index {
									line = line[:i[j][k*2]] + strings.ReplaceAll(
										m[j][k],
										char,
										"‡‡‡DONTESCAPE‡‡‡",
									) + line[i[j][(k*2)+1]:]
								}
							}
						}
						line = strings.ReplaceAll(line, char, "\\"+char)
						line = strings.ReplaceAll(line, "‡‡‡DONTESCAPE‡‡‡", char)
					}
					escape("_") // Escape underscore
					return line
				})
			},
			func(segment string, first bool, last bool) string {
				segment = fmt.Sprintf("`%s`", segment)
				return segment
			},
		)
	}

	return s
}

// NormalizeURLs undoes over-eager escaping inside URLs where backslash-escaped
// underscores would break links. This runs after EscapeCharacters because URLs
// are identified by pattern matching and selectively un-escaped—it's simpler
// than trying to detect URLs before escaping.
func NormalizeURLs(s string, escape bool) string {
	if escape {
		if urls := xurls.Strict().FindAllString(s, -1); len(urls) > 0 {
			for _, url := range urls {
				normalized := strings.ReplaceAll(url, "\\", "")
				s = strings.ReplaceAll(s, url, normalized)
			}
		}
	}
	return s
}

type segmentCallbackFn func(string, bool, bool) string

// processSegments is the core strategy for treating code blocks differently
// from prose. Code blocks must not be escaped or sanitized because they contain
// literal characters (underscores, pipes, etc.) that would be mangled by
// escaping logic meant for prose text. Splitting on the delimiter and
// alternating between normal/code callbacks achieves this cleanly.
func processSegments(s string, prefix string, normalFn segmentCallbackFn, codeFn segmentCallbackFn) string {
	// Isolate blocks of code. Dont escape anything inside them
	nextIsInCodeBlock := strings.HasPrefix(s, prefix)
	segments := strings.Split(s, prefix)
	buffer := bytes.NewBufferString("")
	for i, segment := range segments {
		if len(segment) == 0 {
			continue
		}

		first := i == 0 || len(strings.TrimSpace(segments[i-1])) == 0
		last := i == len(segments)-1 || len(strings.TrimSpace(segments[i+1])) == 0

		if !nextIsInCodeBlock {
			segment = normalFn(segment, first, last)
		} else {
			segment = codeFn(segment, first, last)
		}
		buffer.WriteString(segment)
		nextIsInCodeBlock = !nextIsInCodeBlock
	}
	return buffer.String()
}

func executePerLine(s string, fn func(string) string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = fn(l)
	}
	return strings.Join(lines, "\n")
}
