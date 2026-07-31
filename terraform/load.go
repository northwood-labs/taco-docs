// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/northwood-labs/taco-docs/internal/reader"
	"github.com/northwood-labs/taco-docs/internal/types"
	"github.com/northwood-labs/taco-docs/print"
	"github.com/terraform-docs/terraform-config-inspect/tfconfig"
)

const (
	// httpRequestTimeout is the time limit for fetching remote header/footer
	// content from HTTP URLs.
	httpRequestTimeout = 5

	// minSplitParts is the minimum expected segments when splitting a
	// provider function name by "::".
	minSplitParts = 3
)

// providerFunctionContext holds shared state for provider function
// collection as we walk the HCL AST recursively.
type providerFunctionContext struct {
	Discovered map[string]*ProviderFunction
	Versions   map[string]string
	Sources    map[string]string
}

// LoadWithOptions returns new instance of Module with all the inputs and
// outputs discovered from provided 'path' containing Terraform config.
//
// This is the single entry point that orchestrates the full pipeline: parse →
// enrich → sort. Keeping it as a thin orchestrator makes the loading logic
// testable in isolation (each load* function) while giving callers a simple
// one-call API.
func LoadWithOptions(config *print.Config) (*Module, error) {
	tfmodule, err := loadModule(config.ModuleRoot)
	if err != nil {
		return nil, fmt.Errorf("loading terraform module: %w", err)
	}

	module, err := loadModuleItems(tfmodule, config)
	if err != nil {
		return nil, fmt.Errorf("loading module items: %w", err)
	}

	sortItems(module, config)

	return module, nil
}

// loadModule wraps terraform-config-inspect and filters out known-benign
// diagnostics. OpenTofu supports for_each on provider blocks, which
// terraform-config-inspect doesn't understand and reports as "Invalid provider
// reference". Filtering these prevents false-negative failures when documenting
// OpenTofu modules.
func loadModule(path string) (*tfconfig.Module, error) {
	module, diag := tfconfig.LoadModule(path)
	if diag != nil && diag.HasErrors() {
		// Filter out "Invalid provider reference" errors which can happen
		// with OpenTofu 'for_each' in providers.
		var filteredDiags tfconfig.Diagnostics

		for i := range diag {
			if diag[i].Severity == tfconfig.DiagError &&
				(diag[i].Summary == "Invalid provider reference" ||
					strings.Contains(diag[i].Detail,
						"Provider argument requires a provider name followed by an optional alias")) {
				continue
			}

			filteredDiags = append(filteredDiags, diag[i])
		}

		if filteredDiags.HasErrors() {
			return nil, filteredDiags
		}
	}

	if err := fixOpenTofuProviders(module); err != nil {
		return nil, fmt.Errorf("fixing OpenTofu providers: %w", err)
	}

	return module, nil
}

// fixOpenTofuProviders is a workaround for OpenTofu's for_each in provider
// blocks. terraform-config-inspect leaves Provider.Name empty for resources
// using for_each-iterated providers because it doesn't parse that syntax. We
// re-read the HCL to find the actual provider reference from the resource's
// "provider" attribute, restoring correct provider attribution so loadProviders
// can group resources under the right provider.
func fixOpenTofuProviders(module *tfconfig.Module) error {
	resources := []map[string]*tfconfig.Resource{module.ManagedResources, module.DataResources}
	parser := hclparse.NewParser()

	// cache parsed files to avoid re-reading/re-parsing.
	files := make(map[string]*hcl.File)

	for _, resourceMap := range resources {
		for _, r := range resourceMap {
			// Check if provider is missing or default (empty name/alias issues)
			// If r.Provider.Name is empty, it's definitely broken (since
			// resourceTypeDefaultProviderName always returns something).
			if r.Provider.Name != "" {
				continue
			}

			f, err := getParsedFile(parser, r.Pos.Filename, files)
			if err != nil {
				return fmt.Errorf("parsing file %q: %w", r.Pos.Filename, err)
			}

			if f == nil {
				continue
			}

			if name, alias, ok := findProviderInFile(f, r.Pos.Line); ok {
				r.Provider.Name = name
				r.Provider.Alias = alias
			}
		}
	}

	return nil
}

// Caching parsed files avoids redundant I/O and parsing when multiple resources
// in the same file need provider fixup.
func getParsedFile(parser *hclparse.Parser, filename string, files map[string]*hcl.File) (*hcl.File, error) {
	if f, ok := files[filename]; ok {
		return f, nil
	}

	b, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filename, err)
	}

	f, _ := parser.ParseHCL(b, filename)

	files[filename] = f

	return f, nil
}

// findProviderInFile locates the "provider" attribute in the HCL block starting
// at the given line and extracts the provider name and optional alias from the
// expression traversal. This handles both simple references (provider = aws)
// and indexed references from for_each (provider = aws["key"]) by unwrapping
// IndexExpr before reading the traversal.
func findProviderInFile(f *hcl.File, line int) (string, string, bool) {
	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return "", "", false
	}

	for _, b := range body.Blocks {
		if b.DefRange().Start.Line != line {
			continue
		}

		attr, ok := b.Body.Attributes["provider"]
		if !ok {
			return "", "", false
		}

		expr := attr.Expr
		if idxExpr, ok := expr.(*hclsyntax.IndexExpr); ok {
			expr = idxExpr.Collection
		}

		// Try to get traversal.
		traversal, diags := hcl.AbsTraversalForExpr(expr)
		if diags.HasErrors() || len(traversal) == 0 {
			return "", "", false
		}

		providerName := traversal.RootName()
		alias := ""

		if len(traversal) > 1 {
			if getAttr, ok := traversal[1].(hcl.TraverseAttr); ok {
				alias = getAttr.Name
			}
		}

		return providerName, alias, true
	}

	return "", "", false
}

// loadModuleItems assembles all sub-loaders into a single Module struct. This
// function exists to keep LoadWithOptions focused on orchestration (load →
// sort) while isolating the assembly of individual sections into a testable
// unit.
func loadModuleItems(tfmodule *tfconfig.Module, config *print.Config) (*Module, error) {
	header, err := loadHeader(config)
	if err != nil {
		return nil, fmt.Errorf("loading header section: %w", err)
	}

	footer, err := loadFooter(config)
	if err != nil {
		return nil, fmt.Errorf("loading footer section: %w", err)
	}

	inputs, required, optional := loadInputs(tfmodule, config)
	modulecalls := loadModulecalls(tfmodule, config)

	outputs, err := loadOutputs(tfmodule, config)
	if err != nil {
		return nil, fmt.Errorf("loading outputs: %w", err)
	}

	providers := loadProviders(tfmodule, config)
	providerFunctions := loadProviderFunctions(tfmodule, config)
	requirements := loadRequirements(tfmodule)
	resources := loadResources(tfmodule, config)

	return &Module{
		Header:            header,
		Footer:            footer,
		Inputs:            inputs,
		ModuleCalls:       modulecalls,
		Outputs:           outputs,
		Providers:         providers,
		ProviderFunctions: providerFunctions,
		Requirements:      requirements,
		Resources:         resources,

		RequiredInputs: required,
		OptionalInputs: optional,
	}, nil
}

func getFileFormat(filename string) string {
	if filename == "" {
		return ""
	}

	last := strings.LastIndex(filename, ".")
	if last == -1 {
		return ""
	}

	return filename[last:]
}

func isFileFormatSupported(filename, section string) (bool, error) {
	if section == "" {
		return false, ErrSectionMissing
	}

	if filename == "" {
		return false, fmt.Errorf("%w: --%s-from", ErrSectionFromMissing, section)
	}

	switch getFileFormat(filename) {
	case ".adoc", ".md", ".tf", ".tofu", ".txt":
		return true, nil
	}

	return false, fmt.Errorf(
		"%w: %s",
		ErrUnsupportedFileFormat,
		section,
	)
}

func getSource(filename string) string {
	// Default source is local.
	source := "local"

	// Identify another source different from the local for the filename.
	if strings.HasPrefix(filename, "http") || strings.HasPrefix(filename, "https") ||
		strings.HasPrefix(filename, "s3") {
		source = "web"
	}

	return source
}

func sendHTTPRequest(url string) (string, error) {
	// Creation of context.
	ctx, cancel := context.WithTimeout(context.Background(), httpRequestTimeout*time.Second)
	defer cancel()

	// Send GET request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody) // #nosec G107
	if err != nil {
		fmt.Println("Error:", err)
		return "", fmt.Errorf("creating HTTP request: %w", err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return "", fmt.Errorf("sending HTTP request: %w", err)
	}

	defer func() {
		errDefer := resp.Body.Close()
		if errDefer != nil {
			fmt.Println("Error closing response body:", errDefer)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", fmt.Errorf("reading HTTP response body: %w", err)
	}

	return string(body), nil
}

func createTempFile(config *print.Config, url, content string) (string, error) {
	// Creation of context.
	fileFormat := getFileFormat(url)

	// Pattern with temp-*.* extension.
	tempFile, err := os.CreateTemp("", "temp-*"+fileFormat)
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return "", fmt.Errorf("creating temporary file: %w", err)
	}

	// overrride file name, otherwise it will use the URL and not the temp
	// file created.
	filename := string(filepath.Separator) + filepath.Join(config.ModuleRoot, tempFile.Name())

	// Write the content to the temporary file.
	if _, err := tempFile.WriteString(content); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return "", fmt.Errorf("writing to temporary file: %w", err)
	}

	return filename, nil
}

// loadHeader is a thin wrapper that short-circuits when the header section is
// disabled, avoiding unnecessary file I/O for users who don't want a header in
// their docs.
func loadHeader(config *print.Config) (string, error) {
	if !config.Sections.Header {
		return "", nil
	}

	header, err := loadSection(config, config.HeaderFrom, "header")
	if err != nil {
		return "", fmt.Errorf("loading header: %w", err)
	}

	return header, nil
}

// loadFooter mirrors loadHeader for the footer section, keeping the
// enable/disable logic co-located with the section it controls.
func loadFooter(config *print.Config) (string, error) {
	if !config.Sections.Footer {
		return "", nil
	}

	footer, err := loadSection(config, config.FooterFrom, "footer")
	if err != nil {
		return "", fmt.Errorf("loading footer: %w", err)
	}

	return footer, nil
}

// loadSection reads header/footer content from either a .tf/.tofu file
// (extracting the leading block comment) or a standalone .md/.adoc/.txt file.
// Supporting multiple formats and remote URLs gives module authors flexibility
// in where they maintain their documentation prose without coupling to a single
// convention.
func loadSection(config *print.Config, file, section string) (string, error) { // lint:allow_complexity
	if section == "" {
		return "", ErrSectionMissing
	}

	filename := filepath.Join(config.ModuleRoot, file)
	if ok, err := isFileFormatSupported(file, section); !ok {
		return "", fmt.Errorf("checking file format: %w", err)
	}

	sourceType := getSource(file)

	if sourceType == "web" {
		// Request content of the URL.
		response, err := sendHTTPRequest(file)
		if err != nil {
			fmt.Println("Error:", err)
			return "", fmt.Errorf("fetching remote content: %w", err)
		}

		// Create temp file with the remote content.
		filename, err = createTempFile(config, file, response)
		if err != nil {
			fmt.Println("Error:", err)
			return "", fmt.Errorf("creating temp file for remote content: %w", err)
		}

		// Ensure the temporary file is removed.
		defer func() {
			errDefer := os.Remove(filename)
			if errDefer != nil {
				fmt.Println("Error removing temporary file:", errDefer)
			}
		}()
	}

	if info, err := os.Stat(filename); os.IsNotExist(err) || info.IsDir() {
		if section == "header" && file == "main.tf" {
			// absorb the error to not break workflow for default value of header and missing 'main.tf'.
			return "", nil
		}

		// user explicitly asked for a file which doesn't exist.
		return "", fmt.Errorf(
			"accessing file %q: %w",
			filename,
			err,
		)
	}

	format := getFileFormat(file)
	if format != ".tf" && format != ".tofu" {
		content, err := os.ReadFile(filepath.Clean(filename))
		if err != nil {
			return "", fmt.Errorf("reading file %q: %w", filename, err)
		}

		return string(content), nil
	}

	lines := reader.Lines{
		FileName: filename,
		LineNum:  -1,
		Condition: func(line string) bool {
			line = strings.TrimSpace(line)
			return strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "*/")
		},
		Parser: func(line string) (string, bool) {
			tmp := strings.TrimSpace(line)
			if strings.HasPrefix(tmp, "/*") || strings.HasPrefix(tmp, "*/") {
				return "", false
			}

			if tmp == "*" {
				return "", true
			}

			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, "\r\n")
			line = strings.TrimPrefix(line, "* ")

			return line, true
		},
	}

	sectionText, err := lines.Extract()
	if err != nil {
		return "", fmt.Errorf("extracting %s lines: %w", section, err)
	}

	return strings.Join(sectionText, "\n"), nil
}

// loadInputs converts terraform-config-inspect's raw Variable structs into our
// enriched Input model. Key enrichments: (1) reading comments above the
// declaration as a fallback description when the variable lacks an explicit
// description attribute (controlled by readComments), (2) normalizing CRLF to
// LF for cross-platform consistency, and (3) pre-partitioning into
// required/optional slices for downstream template convenience.
func loadInputs(tfmodule *tfconfig.Module, config *print.Config) ([]*Input, []*Input, []*Input) {
	inputs := make([]*Input, 0, len(tfmodule.Variables))
	required := make([]*Input, 0, len(tfmodule.Variables))
	optional := make([]*Input, 0, len(tfmodule.Variables))

	for _, input := range tfmodule.Variables {
		comments := loadComments(input.Pos.Filename, input.Pos.Line)

		// skip over inputs that are marked as being ignored.
		if strings.Contains(comments, "terraform-docs-ignore") {
			continue
		}

		// convert CRLF to LF early on
		// (https://github.com/northwood-labs/taco-docs/issues/305)
		inputDescription := strings.ReplaceAll(input.Description, "\r\n", "\n")
		if inputDescription == "" && config.Settings.ReadComments {
			inputDescription = comments
		}

		i := &Input{
			Name:        input.Name,
			Type:        types.TypeOf(input.Type, input.Default),
			Description: types.String(inputDescription),
			Default:     types.ValueOf(input.Default),
			Required:    input.Required,
			Position: Position{
				Filename: input.Pos.Filename,
				Line:     input.Pos.Line,
			},
			Sensitive: input.Sensitive,
		}

		inputs = append(inputs, i)

		if i.HasDefault() {
			optional = append(optional, i)
		} else {
			required = append(required, i)
		}
	}

	return inputs, required, optional
}

// formatSource separates an inline "?ref=" version suffix from a module source
// URL. Many git-sourced modules embed the version in the source string rather
// than using a separate "version" field, so we extract it to populate Version
// consistently across all module call types.
func formatSource(s, v string) (string, string) {
	substr := "?ref="

	if v != "" {
		return s, v
	}

	pos := strings.LastIndex(s, substr)
	if pos == -1 {
		return s, ""
	}

	adjustedPos := pos + len(substr)
	if adjustedPos >= len(s) {
		return s, ""
	}

	source := s[0:pos]
	version := s[adjustedPos:]

	return source, version
}

func loadModulecalls(tfmodule *tfconfig.Module, config *print.Config) []*ModuleCall {
	var modules []*ModuleCall

	var source, version string

	for _, m := range tfmodule.ModuleCalls {
		comments := loadComments(m.Pos.Filename, m.Pos.Line)

		// skip over modules that are marked as being ignored.
		if strings.Contains(comments, "terraform-docs-ignore") {
			continue
		}

		description := ""
		if config.Settings.ReadComments {
			description = comments
		}

		source, version = formatSource(m.Source, m.Version)

		modules = append(modules, &ModuleCall{
			Name:        m.Name,
			Source:      source,
			Version:     version,
			Description: types.String(description),
			Position: Position{
				Filename: m.Pos.Filename,
				Line:     m.Pos.Line,
			},
		})
	}

	return modules
}

// loadOutputs enriches raw output declarations with optional actual values from
// `terraform output -json`. When --output-values is active, it merges the
// runtime state (value + sensitivity) into the documentation model so users can
// see what a deployed module actually exports—useful for internal wikis or
// post-apply documentation.
func loadOutputs(tfmodule *tfconfig.Module, config *print.Config) ([]*Output, error) {
	outputs := make([]*Output, 0, len(tfmodule.Outputs))
	values := make(map[string]*output)

	if config.OutputValues.Enabled {
		var err error

		values, err = loadOutputValues(config)
		if err != nil {
			return nil, fmt.Errorf("loading output values: %w", err)
		}
	}

	for _, o := range tfmodule.Outputs {
		comments := loadComments(o.Pos.Filename, o.Pos.Line)

		// skip over outputs that are marked as being ignored.
		if strings.Contains(comments, "terraform-docs-ignore") {
			continue
		}

		// convert CRLF to LF early on
		// (https://github.com/northwood-labs/taco-docs/issues/584)
		description := strings.ReplaceAll(o.Description, "\r\n", "\n")
		if description == "" && config.Settings.ReadComments {
			description = comments
		}

		output := &Output{
			Name:        o.Name,
			Description: types.String(description),
			Position: Position{
				Filename: o.Pos.Filename,
				Line:     o.Pos.Line,
			},
			ShowValue: config.OutputValues.Enabled,
		}

		if config.OutputValues.Enabled {
			if value, ok := values[output.Name]; ok {
				output.Sensitive = value.Sensitive
				output.Value = types.ValueOf(value.Value)
			} else {
				output.Value = types.ValueOf("null")
			}

			if output.Sensitive {
				output.Value = types.ValueOf(`<sensitive>`)
			}
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

func loadOutputValues(config *print.Config) (map[string]*output, error) {
	var (
		out []byte
		err error
	)
	if config.OutputValues.From == "" {
		cmd := exec.CommandContext(context.TODO(), "terraform", "output", "-json")

		cmd.Dir = config.ModuleRoot
		if out, err = cmd.Output(); err != nil {
			return nil, fmt.Errorf("caught error while reading the terraform outputs: %w", err)
		}
	} else if out, err = os.ReadFile(config.OutputValues.From); err != nil {
		return nil, fmt.Errorf(
			"caught error while reading the terraform outputs file at %s: %w",
			config.OutputValues.From,
			err,
		)
	}

	var terraformOutputs map[string]*output

	err = json.Unmarshal(out, &terraformOutputs)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling terraform outputs: %w", err)
	}

	return terraformOutputs, nil
}

// loadProviders discovers providers from actual resource usage rather than only
// from required_providers. This ensures documentation reflects runtime
// dependencies even when the author forgot to declare a provider in
// required_providers. It also merges version info from .terraform.lock.hcl (if
// --lock-file is set) to show the exact installed version, giving readers a
// more accurate picture than just the constraint string.
func loadProviders(tfmodule *tfconfig.Module, config *print.Config) []*Provider { // lint:allow_complexity
	type provider struct {
		Name        string   `hcl:"name,label"`
		Version     string   `hcl:"version"`
		Constraints *string  `hcl:"constraints"`
		Hashes      []string `hcl:"hashes"`
	}

	type lockfile struct {
		Provider []provider `hcl:"provider,block"`
	}

	lock := make(map[string]provider)

	if config.Settings.LockFile {
		var lf lockfile

		filename := filepath.Join(config.ModuleRoot, ".terraform.lock.hcl")
		if err := hclsimple.DecodeFile(filename, nil, &lf); err == nil {
			for i := range lf.Provider {
				segments := strings.Split(lf.Provider[i].Name, "/")
				name := segments[len(segments)-1]

				lock[name] = lf.Provider[i]
			}
		}
	}

	resources := []map[string]*tfconfig.Resource{tfmodule.ManagedResources, tfmodule.DataResources}
	discovered := make(map[string]*Provider)

	for _, resource := range resources {
		for _, r := range resource {
			comments := loadComments(r.Pos.Filename, r.Pos.Line)

			// skip over resources that are marked as being ignored.
			if strings.Contains(comments, "terraform-docs-ignore") {
				continue
			}

			version := ""
			if l, ok := lock[r.Provider.Name]; ok {
				version = l.Version
			} else if rv, ok := tfmodule.RequiredProviders[r.Provider.Name]; ok && len(rv.VersionConstraints) > 0 {
				version = strings.Join(rv.VersionConstraints, " ")
			}

			key := fmt.Sprintf("%s.%s", r.Provider.Name, r.Provider.Alias)
			if existing, ok := discovered[key]; ok {
				// keep the earliest position across all resources of this
				// provider.
				if r.Pos.Filename < existing.Position.Filename ||
					(r.Pos.Filename == existing.Position.Filename && r.Pos.Line < existing.Position.Line) {
					existing.Position = Position{Filename: r.Pos.Filename, Line: r.Pos.Line}
				}

				continue
			}

			discovered[key] = &Provider{
				Name:    r.Provider.Name,
				Alias:   types.String(r.Provider.Alias),
				Version: types.String(version),
				Position: Position{
					Filename: r.Pos.Filename,
					Line:     r.Pos.Line,
				},
			}
		}
	}

	keys := make([]string, 0, len(discovered))
	for key := range discovered {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	providers := make([]*Provider, 0, len(discovered))
	for _, key := range keys {
		providers = append(providers, discovered[key])
	}

	return providers
}

// loadProviderFunctions walks the HCL AST of every .tf file to find
// "provider::name::func" function call expressions. These are a newer
// Terraform/OpenTofu feature not tracked by terraform-config-inspect, so we
// must perform our own AST traversal. Discovering them lets us document which
// provider-supplied functions a module relies on—information otherwise
// invisible in standard module interfaces.
func loadProviderFunctions( // lint:allow_complexity
	tfmodule *tfconfig.Module,
	config *print.Config,
) []*ProviderFunction {
	providerVersions := make(map[string]string)
	providerSources := make(map[string]string)

	for name, provider := range tfmodule.RequiredProviders {
		if len(provider.VersionConstraints) > 0 {
			providerVersions[name] = strings.Join(provider.VersionConstraints, " ")
		}

		if provider.Source != "" {
			providerSources[name] = provider.Source
		} else {
			providerSources[name] = fmt.Sprintf("%s/%s", "hashicorp", name)
		}
	}

	discovered := make(map[string]*ProviderFunction)
	parser := hclparse.NewParser()

	_ = filepath.WalkDir(config.ModuleRoot, func(path string, d fs.DirEntry, err error) error { // lint:allow_unhandled
		if err != nil {
			return err
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return fs.SkipDir
			}

			return nil
		}

		if filepath.Ext(path) != ".tf" {
			return nil
		}

		file, diags := parser.ParseHCLFile(path)
		if diags.HasErrors() {
			return nil
		}

		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			return nil
		}

		collectProviderFunctions(body, path, &providerFunctionContext{
			Discovered: discovered,
			Versions:   providerVersions,
			Sources:    providerSources,
		})

		return nil
	})

	keys := make([]string, 0, len(discovered))
	for key := range discovered {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	providerFunctions := make([]*ProviderFunction, 0, len(discovered))
	for _, key := range keys {
		providerFunctions = append(providerFunctions, discovered[key])
	}

	return providerFunctions
}

func collectProviderFunctions(
	body *hclsyntax.Body,
	filename string,
	ctx *providerFunctionContext,
) {
	for _, attr := range body.Attributes {
		collectProviderFunctionsFromExpr(attr.Expr, filename, ctx)
	}

	for _, block := range body.Blocks {
		collectProviderFunctions(block.Body, filename, ctx)
	}
}

func collectProviderFunctionsFromExpr( // lint:allow_complexity
	expr hclsyntax.Expression,
	filename string,
	ctx *providerFunctionContext,
) {
	switch t := expr.(type) {
	case *hclsyntax.FunctionCallExpr:
		if strings.HasPrefix(t.Name, "provider::") {
			parts := strings.SplitN(t.Name, "::", minSplitParts)
			if len(parts) == minSplitParts {
				providerName := parts[1]
				functionName := parts[2]

				key := fmt.Sprintf("%s.%s", providerName, functionName)
				if _, ok := ctx.Discovered[key]; !ok {
					version := types.String(ctx.Versions[providerName])

					source := ctx.Sources[providerName]
					if source == "" {
						source = fmt.Sprintf("%s/%s", "hashicorp", providerName)
					}

					ctx.Discovered[key] = &ProviderFunction{
						ProviderName:   providerName,
						Function:       functionName,
						ProviderSource: source,
						Version:        version,
						Position: Position{
							Filename: filename,
							Line:     t.Range().Start.Line,
						},
					}
				}
			}
		}

		for _, arg := range t.Args {
			collectProviderFunctionsFromExpr(arg, filename, ctx)
		}
	case *hclsyntax.TemplateExpr:
		for _, part := range t.Parts {
			collectProviderFunctionsFromExpr(part, filename, ctx)
		}
	case *hclsyntax.TemplateWrapExpr:
		collectProviderFunctionsFromExpr(t.Wrapped, filename, ctx)
	case *hclsyntax.TupleConsExpr:
		for _, expr := range t.Exprs {
			collectProviderFunctionsFromExpr(expr, filename, ctx)
		}
	case *hclsyntax.ObjectConsExpr:
		for _, item := range t.Items {
			collectProviderFunctionsFromExpr(item.KeyExpr, filename, ctx)
			collectProviderFunctionsFromExpr(item.ValueExpr, filename, ctx)
		}
	case *hclsyntax.ConditionalExpr:
		collectProviderFunctionsFromExpr(t.Condition, filename, ctx)
		collectProviderFunctionsFromExpr(t.TrueResult, filename, ctx)
		collectProviderFunctionsFromExpr(t.FalseResult, filename, ctx)
	case *hclsyntax.ForExpr:
		if t.KeyExpr != nil {
			collectProviderFunctionsFromExpr(t.KeyExpr, filename, ctx)
		}

		if t.ValExpr != nil {
			collectProviderFunctionsFromExpr(t.ValExpr, filename, ctx)
		}

		collectProviderFunctionsFromExpr(t.CollExpr, filename, ctx)

		if t.CondExpr != nil {
			collectProviderFunctionsFromExpr(t.CondExpr, filename, ctx)
		}
	case *hclsyntax.SplatExpr:
		collectProviderFunctionsFromExpr(t.Source, filename, ctx)

		if t.Each != nil {
			collectProviderFunctionsFromExpr(t.Each, filename, ctx)
		}
	case *hclsyntax.IndexExpr:
		collectProviderFunctionsFromExpr(t.Collection, filename, ctx)
		collectProviderFunctionsFromExpr(t.Key, filename, ctx)
	case *hclsyntax.UnaryOpExpr:
		collectProviderFunctionsFromExpr(t.Val, filename, ctx)
	case *hclsyntax.BinaryOpExpr:
		collectProviderFunctionsFromExpr(t.LHS, filename, ctx)
		collectProviderFunctionsFromExpr(t.RHS, filename, ctx)
	case *hclsyntax.ParenthesesExpr:
		collectProviderFunctionsFromExpr(t.Expression, filename, ctx)
	case *hclsyntax.RelativeTraversalExpr:
		collectProviderFunctionsFromExpr(t.Source, filename, ctx)
	default:
	}
}

func loadRequirements(tfmodule *tfconfig.Module) []*Requirement {
	var requirements []*Requirement
	for _, core := range tfmodule.RequiredCore {
		requirements = append(requirements, &Requirement{
			Name:    "terraform",
			Source:  "hashicorp/terraform",
			Version: types.String(core),
		})
	}

	names := make([]string, 0, len(tfmodule.RequiredProviders))
	for n := range tfmodule.RequiredProviders {
		names = append(names, n)
	}

	slices.Sort(names)

	for _, name := range names {
		for _, version := range tfmodule.RequiredProviders[name].VersionConstraints {
			var source string
			if tfmodule.RequiredProviders[name].Source != "" {
				source = tfmodule.RequiredProviders[name].Source
			} else {
				source = fmt.Sprintf("%s/%s", "hashicorp", name)
			}

			requirements = append(requirements, &Requirement{
				Name:    name,
				Source:  types.String(source),
				Version: types.String(version),
			})
		}
	}

	return requirements
}

// loadResources deduplicates resources by a composite key
// (provider.mode.type.name) because the same resource declaration appears in
// both ManagedResources and DataResources maps from terraform-config-inspect.
// Using a map with a stable key prevents duplicate entries in docs while
// building registry URLs for direct documentation links.
func loadResources(tfmodule *tfconfig.Module, config *print.Config) []*Resource {
	allResources := []map[string]*tfconfig.Resource{tfmodule.ManagedResources, tfmodule.DataResources}
	discovered := make(map[string]*Resource)

	for _, resource := range allResources {
		for _, r := range resource {
			comments := loadComments(r.Pos.Filename, r.Pos.Line)

			// skip over resources that are marked as being ignored.
			if strings.Contains(comments, "terraform-docs-ignore") {
				continue
			}

			var version string
			if rv, ok := tfmodule.RequiredProviders[r.Provider.Name]; ok {
				version = resourceVersion(rv.VersionConstraints)
			}

			var source string
			if tfmodule.RequiredProviders[r.Provider.Name].Source != "" {
				source = tfmodule.RequiredProviders[r.Provider.Name].Source
			} else {
				source = fmt.Sprintf("%s/%s", "hashicorp", r.Provider.Name)
			}

			rType := strings.TrimPrefix(r.Type, r.Provider.Name+"_")
			key := fmt.Sprintf("%s.%s.%s.%s", r.Provider.Name, r.Mode, rType, r.Name)

			description := ""
			if config.Settings.ReadComments {
				description = comments
			}

			discovered[key] = &Resource{
				Type:           rType,
				Name:           r.Name,
				Mode:           r.Mode.String(),
				ProviderName:   r.Provider.Name,
				ProviderSource: source,
				Version:        types.String(version),
				Description:    types.String(description),
				Position: Position{
					Filename: r.Pos.Filename,
					Line:     r.Pos.Line,
				},
			}
		}
	}

	resourceKeys := make([]string, 0, len(discovered))
	for key := range discovered {
		resourceKeys = append(resourceKeys, key)
	}

	slices.Sort(resourceKeys)

	resources := make([]*Resource, 0, len(discovered))
	for _, key := range resourceKeys {
		resources = append(resources, discovered[key])
	}

	return resources
}

// resourceVersion normalizes the raw version constraint strings into a single
// version suitable for registry URL construction. It uses the last constraint
// (most specific) and strips operators to extract a bare version number. Falls
// back to "latest" when the constraint is too complex or absent, ensuring
// registry links still resolve to something useful.
func resourceVersion(constraints []string) string { // lint:allow_complexity
	if len(constraints) == 0 {
		return "latest"
	}

	versionParts := strings.Split(constraints[len(constraints)-1], " ")
	switch len(versionParts) {
	case 1:
		if _, err := strconv.Atoi(versionParts[0][0:1]); err != nil {
			if versionParts[0][0:1] == "=" {
				return versionParts[0][1:]
			}

			return "latest"
		}

		return versionParts[0]
	case 2:
		if versionParts[0] == "=" {
			return versionParts[1]
		}
	default:
	}

	return "latest"
}

// loadComments extracts the contiguous block of # or // comments immediately
// above a declaration. These comments serve as fallback descriptions when the
// variable/output/resource lacks an explicit "description" attribute—a pattern
// many module authors use for brevity. Errors are absorbed (returning "")
// because a missing comment should never break doc generation.
func loadComments(filename string, lineNum int) string {
	lines := reader.Lines{
		FileName: filename,
		LineNum:  lineNum,
		Condition: func(line string) bool {
			return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//")
		},
		Parser: func(line string) (string, bool) {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimSpace(line)

			return line, true
		},
	}

	comment, err := lines.Extract()
	if err != nil {
		return "" // absorb the error, we don't need to bubble it up or break the execution.
	}

	return strings.Join(comment, " ")
}

// sortItems applies the user's configured sort preferences to every collection
// type in the module. Sorting is deferred to this final step (after all loading
// is complete) so that load functions can append items in discovery order
// without worrying about ordering, and the sort strategy can be changed without
// touching any loader logic.
func sortItems(tfmodule *Module, config *print.Config) {
	// inputs.
	inputs(tfmodule.Inputs).sort(config.Sort.Enabled, config.Sort.By)
	inputs(tfmodule.RequiredInputs).sort(config.Sort.Enabled, config.Sort.By)
	inputs(tfmodule.OptionalInputs).sort(config.Sort.Enabled, config.Sort.By)

	// outputs.
	outputs(tfmodule.Outputs).sort(config.Sort.Enabled, config.Sort.By)

	// providers.
	providers(tfmodule.Providers).sort(config.Sort.Enabled, config.Sort.By)

	// provider functions.
	providerFunctions(tfmodule.ProviderFunctions).sort(config.Sort.Enabled, config.Sort.By)

	// resources.
	resources(tfmodule.Resources).sort(config.Sort.Enabled, config.Sort.By)

	// modules.
	modulecalls(tfmodule.ModuleCalls).sort(config.Sort.Enabled, config.Sort.By)
}
