// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package print

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	sectionAll               = "all"
	sectionDataSources       = "data-sources"
	sectionFooter            = "footer"
	sectionHeader            = "header"
	sectionInputs            = "inputs"
	sectionModules           = "modules"
	sectionOutputs           = "outputs"
	sectionProviders         = "providers"
	sectionProviderFunctions = "provider-functions"
	sectionRequirements      = "requirements"
	sectionResources         = "resources"

	// OutputModeInject is the output mode that injects between markers.
	OutputModeInject = "inject"
	// OutputModeReplace is the output mode that replaces the entire file.
	OutputModeReplace = "replace"

	// OutputBeginComment is the default begin marker for inject mode.
	OutputBeginComment = "<!-- BEGIN_TF_DOCS -->"
	// OutputContent is the template placeholder for generated content.
	OutputContent = "{{ .Content }}"
	// OutputEndComment is the default end marker for inject mode.
	OutputEndComment = "<!-- END_TF_DOCS -->"

	// Sort types.

	// SortName sorts items alphabetically by name.
	SortName = "name"
	// SortRequired sorts items with required first.
	SortRequired = "required"
	// SortType sorts items by their type.
	SortType = "type"

	// minTemplateLines is the minimum number of lines an output template must
	// contain: begin comment, {{ .Content }}, and end comment.
	minTemplateLines = 3
)

// Output to file template and modes.
var (
	allSections = []string{
		sectionAll,
		sectionDataSources,
		sectionFooter,
		sectionHeader,
		sectionInputs,
		sectionModules,
		sectionOutputs,
		sectionProviders,
		sectionProviderFunctions,
		sectionRequirements,
		sectionResources,
	}

	// AllSections list.
	AllSections = strings.Join(allSections, ", ")

	OutputTemplate = fmt.Sprintf("%s\n%s\n%s", OutputBeginComment, OutputContent, OutputEndComment)
	OutputModes    = OutputModeInject + ", " + OutputModeReplace

	allSorts = []string{
		SortName,
		SortRequired,
		SortType,
	}

	// SortTypes list.
	SortTypes = strings.Join(allSorts, ", ")
)

type (
	// Config is the central data model for all user preferences. It serves as
	// the single source of truth passed to formatters, templates, and plugins.
	// mapstructure tags enable viper to decode YAML config files directly into
	// this struct without manual field-by-field assignment.
	Config struct {
		Output       output       `mapstructure:"output"`
		Sort         sort         `mapstructure:"sort"`
		OutputValues outputvalues `mapstructure:"output-values"`
		HeaderFrom   string       `mapstructure:"header-from"`
		FooterFrom   string       `mapstructure:"footer-from"`
		Content      string       `mapstructure:"content"`
		File         string       `mapstructure:"-"`
		Version      string       `mapstructure:"version"`
		Formatter    string       `mapstructure:"formatter"`
		ModuleRoot   string
		Recursive    recursive `mapstructure:"recursive"`
		Sections     sections  `mapstructure:"sections"`
		Settings     settings  `mapstructure:"settings"`
	}

	recursive struct {
		Path        string   `mapstructure:"path"`
		Exclude     []string `mapstructure:"exclude"`
		Enabled     bool     `mapstructure:"enabled"`
		IncludeMain bool     `mapstructure:"include-main"`
	}

	sections struct {
		Show []string `mapstructure:"show"`
		Hide []string `mapstructure:"hide"`

		DataSources       bool
		Header            bool
		Footer            bool
		Inputs            bool
		ModuleCalls       bool
		Outputs           bool
		Providers         bool
		ProviderFunctions bool
		Requirements      bool
		Resources         bool
	}

	output struct {
		File         string `mapstructure:"file"`
		Mode         string `mapstructure:"mode"`
		Template     string `mapstructure:"template"`
		BeginComment string
		EndComment   string
		Check        bool
	}

	outputvalues struct {
		From    string `mapstructure:"from"`
		Enabled bool   `mapstructure:"enabled"`
	}

	sort struct {
		By      string `mapstructure:"by"`
		Enabled bool   `mapstructure:"enabled"`
	}

	settings struct {
		Anchor       bool `mapstructure:"anchor"`
		AtxClosed    bool `mapstructure:"atx-closed"`
		Color        bool `mapstructure:"color"`
		Default      bool `mapstructure:"default"`
		Description  bool `mapstructure:"description"`
		Escape       bool `mapstructure:"escape"`
		HideEmpty    bool `mapstructure:"hide-empty"`
		HTML         bool `mapstructure:"html"`
		Indent       int  `mapstructure:"indent"`
		LockFile     bool `mapstructure:"lockfile"`
		ReadComments bool `mapstructure:"read-comments"`
		Required     bool `mapstructure:"required"`
		Sensitive    bool `mapstructure:"sensitive"`
		Type         bool `mapstructure:"type"`
	}
)

// NewConfig returns neew instancee of Config with empty values.
func NewConfig() *Config {
	return &Config{
		HeaderFrom:   "main.tf",
		Recursive:    recursive{},
		Sections:     sections{},
		Output:       output{},
		OutputValues: outputvalues{},
		Sort:         sort{},
		Settings:     settings{},
	}
}

// DefaultConfig provides safe defaults that produce useful output without
// any configuration. This implements the "convention over configuration"
// principle— users get a reasonable result on first run and only need to
// customize what they want to change.
func DefaultConfig() *Config {
	return &Config{
		File:         "",
		Formatter:    "",
		Version:      "",
		HeaderFrom:   "main.tf",
		FooterFrom:   "",
		Recursive:    defaultRecursive(),
		Content:      "",
		Sections:     defaultSections(),
		Output:       defaultOutput(),
		OutputValues: defaultOutputValues(),
		Sort:         defaultSort(),
		Settings:     defaultSettings(),

		ModuleRoot: "",
	}
}

func defaultRecursive() recursive {
	return recursive{
		Enabled:     false,
		Path:        "modules",
		IncludeMain: true,
		Exclude:     nil,
	}
}

func (r *recursive) validate() error {
	if r.Enabled && r.Path == "" {
		return ErrRecursivePathEmpty
	}

	return nil
}

func defaultSections() sections {
	return sections{
		Show: nil,
		Hide: nil,

		DataSources:       true,
		Header:            true,
		Footer:            false,
		Inputs:            true,
		ModuleCalls:       true,
		Outputs:           true,
		Providers:         true,
		ProviderFunctions: true,
		Requirements:      true,
		Resources:         true,
	}
}

func (s *sections) validate() error {
	if len(s.Show) > 0 && len(s.Hide) > 0 {
		return ErrShowHideConflict
	}

	for _, item := range s.Show {
		if !contains(allSections, item) {
			return fmt.Errorf("%w: '%s'", ErrInvalidSection, item)
		}
	}

	for _, item := range s.Hide {
		if !contains(allSections, item) {
			return fmt.Errorf("%w: '%s'", ErrInvalidSection, item)
		}
	}

	return nil
}

// visibility implements "show wins over hide" priority logic. When --show is
// specified, only those sections appear (whitelist). When --hide is specified,
// everything except those sections appears (blacklist). If neither is set, all
// sections are visible. This two-mode approach gives users a concise way to
// either opt-in or opt-out of specific sections.
func (s *sections) visibility(section string) bool {
	if len(s.Show) == 0 && len(s.Hide) == 0 {
		return true
	}

	for _, n := range s.Show {
		if n == sectionAll || n == section {
			return true
		}
	}

	for _, n := range s.Hide {
		if n == sectionAll || n == section {
			return false
		}
	}
	// hidden : if s.Show NOT empty AND s.Show does NOT contain section
	// visible: if s.Hide NOT empty AND s.Hide does NOT contain section
	return len(s.Hide) > 0
}

func defaultOutput() output {
	return output{
		File:     "",
		Mode:     OutputModeInject,
		Template: OutputTemplate,
		Check:    false,

		BeginComment: OutputBeginComment,
		EndComment:   OutputEndComment,
	}
}

// validate ensures template markers are well-formed to prevent silent injection
// failures. Without this check, a malformed template would cause the generated
// content to be silently dropped or written to the wrong location in the target
// file when using --output-file with inject mode.
func (o *output) validate() error {
	if o.File == "" {
		return nil
	}

	if o.Mode == "" {
		return ErrOutputModeEmpty
	}

	// Template is optional for mode 'replace'.
	if o.Mode == OutputModeReplace && o.Template == "" {
		return nil
	}

	if o.Template == "" {
		return ErrOutputTemplateEmpty
	}

	if !strings.Contains(o.Template, OutputContent) {
		return ErrOutputTemplateNoContent
	}

	// No extra validation is needed for mode 'replace', the following only
	// applies for every other modes.
	if o.Mode == OutputModeReplace {
		return nil
	}

	o.Template = strings.ReplaceAll(o.Template, "\\n", "\n")

	lines := strings.Split(o.Template, "\n")
	tests := []struct {
		condition func() bool
		err       error
	}{
		{
			condition: func() bool {
				return len(lines) < minTemplateLines
			},
			err: ErrOutputTemplateTooShort,
		},
		{
			condition: func() bool {
				return !isInlineComment(strings.TrimSpace(lines[0]))
			},
			err: ErrOutputTemplateNoBegin,
		},
		{
			condition: func() bool {
				return !isInlineComment(strings.TrimSpace(lines[len(lines)-1]))
			},
			err: ErrOutputTemplateNoEnd,
		},
	}

	for _, t := range tests {
		if t.condition() {
			return t.err
		}
	}

	o.BeginComment = strings.TrimSpace(lines[0])
	o.EndComment = strings.TrimSpace(lines[len(lines)-1])

	return nil
}

// isInlineComment recognizes multiple comment syntaxes so that begin/end
// markers work correctly in both Markdown and AsciiDoc files. Users can choose
// whichever comment style their target format supports.
//
// Detect if a particular line is a Markdown comment.
//
// ref: https://www.jamestharpe.com/markdown-comments/
func isInlineComment(line string) bool {
	switch {
	// AsciiDoc specific.
	case strings.HasPrefix(line, "//"):
		return true

	// Markdown specific.
	default:
		cases := [][]string{
			{"<!--", "-->"},
			{"[]: # (", ")"},
			{"[]: # \"", "\""},
			{"[]: # '", "'"},
			{"[//]: # (", ")"},
			{"[comment]: # (", ")"},
		}
		for _, c := range cases {
			if strings.HasPrefix(line, c[0]) && strings.HasSuffix(line, c[1]) {
				return true
			}
		}
	}

	return false
}

func defaultOutputValues() outputvalues {
	return outputvalues{
		Enabled: false,
		From:    "",
	}
}

func (o *outputvalues) validate() error {
	if o.Enabled && o.From == "" {
		return ErrOutputValuesFromEmpty
	}

	return nil
}

func defaultSort() sort {
	return sort{
		Enabled: true,
		By:      SortName,
	}
}

func (s *sort) validate() error {
	if !contains(allSorts, s.By) {
		return fmt.Errorf("%w: '%s'", ErrInvalidSortType, s.By)
	}

	return nil
}

func defaultSettings() settings {
	return settings{
		Anchor:       true,
		AtxClosed:    false,
		Color:        true,
		Default:      true,
		Description:  false,
		Escape:       true,
		HideEmpty:    false,
		HTML:         true,
		Indent:       2,
		LockFile:     true,
		ReadComments: true,
		Required:     true,
		Sensitive:    true,
		Type:         true,
	}
}

func (*settings) validate() error {
	return nil
}

// Parse translates user-facing show/hide lists into boolean flags that
// templates can check efficiently. Templates use these booleans to decide which
// sections to render, avoiding repeated list lookups during template execution.
func (c *Config) Parse() {
	// sections.
	c.Sections.DataSources = c.Sections.visibility("data-sources")
	c.Sections.Header = c.Sections.visibility("header")
	c.Sections.Inputs = c.Sections.visibility("inputs")
	c.Sections.ModuleCalls = c.Sections.visibility("modules")
	c.Sections.Outputs = c.Sections.visibility("outputs")
	c.Sections.Providers = c.Sections.visibility("providers")
	c.Sections.ProviderFunctions = c.Sections.visibility("provider-functions")
	c.Sections.Requirements = c.Sections.visibility("requirements")
	c.Sections.Resources = c.Sections.visibility("resources")

	// Footer section is optional and should only be enabled if --footer-from is
	// explicitly set, either via CLI or config file.
	if c.FooterFrom != "" {
		c.Sections.Footer = c.Sections.visibility("footer")
	}
}

// Validate catches misconfiguration early with clear error messages before
// expensive module parsing begins. Failing fast here means users don't wait for
// Terraform file traversal only to discover an invalid option.
func (c *Config) Validate() error {
	// formatter.
	if c.Formatter == "" {
		return ErrFormatterEmpty
	}

	// header-from.
	if c.HeaderFrom == "" {
		return ErrHeaderFromEmpty
	}

	// footer-from, not a 'default' section so can be empty.
	if c.Sections.Footer && c.FooterFrom == "" {
		return ErrFooterFromEmpty
	}

	if c.FooterFrom == c.HeaderFrom {
		return ErrFooterEqualsHeader
	}

	for _, fn := range []func() error{
		c.Recursive.validate,
		c.Sections.validate,
		c.Output.validate,
		c.OutputValues.validate,
		c.Sort.validate,
		c.Settings.validate,
	} {
		if err := fn(); err != nil {
			return fmt.Errorf("validating config: %w", err)
		}
	}

	return nil
}

// ReadConfig is a standalone config reader for use outside the CLI (e.g., in
// tests or the plugin SDK). It encapsulates the full read → unmarshal →
// validate → parse lifecycle so callers don't need to replicate the sequencing
// themselves.
func ReadConfig(rootDir, filename string) (*Config, error) {
	cfg := NewConfig()

	v := viper.New()
	v.SetConfigFile(filepath.Join(rootDir, filename))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := errors.AsType[*os.PathError](err); ok {
			return nil, fmt.Errorf("%w: %s", ErrConfigFileNotFound, filename)
		}

		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config, %w", err)
	}

	cfg.ModuleRoot = rootDir

	// process and validate configuration.
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	cfg.Parse()

	return cfg, nil
}
