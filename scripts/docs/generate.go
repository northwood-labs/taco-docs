// Copyright 2018-2026 The terraform-docs Authors.
// Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the MIT license (the "License"); you may not
// use this file except in compliance with the License.
//
// You may obtain a copy of the License at the LICENSE file in
// the root directory of this source tree.

// Package main generates reference documentation for CLI commands.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/northwood-labs/taco-docs/cmd"
	"github.com/northwood-labs/taco-docs/format"
	"github.com/northwood-labs/taco-docs/print"
	"github.com/northwood-labs/taco-docs/terraform"
)

// These are practiaclly a copy/paste of https://github.com/spf13/cobra/blob/master/doc/md_docs.go
// The reason we've decided to bring them over and not use them directly
// from cobra module was that we wanted to inject custom "Example" section
// with generated output based on the "examples" folder.

var baseWeight = 950

type (
	reference struct {
		InheritedOptions string
		Usage            string
		Description      string
		Parent           string
		Synopsis         string
		UseLine          string
		Options          string
		Example          string
		Command          string
		Name             string
		Subcommands      []command
		Weight           int
		HasChildren      bool
		Runnable         bool
	}

	command struct {
		Name     string
		Link     string
		Children []command
	}
)

func main() {
	if err := generate(cmd.NewCommand(), baseWeight, "terraform-docs"); err != nil {
		log.Fatal(err)
	}
}

func ignore(command *cobra.Command) bool {
	switch {
	case !command.IsAvailableCommand():
		return true
	case command.IsAdditionalHelpTopicCommand():
		return true
	case command.Annotations["kind"] == "":
		return true
	case command.Annotations["kind"] != "formatter":
		return true
	}

	return false
}

func generate(command *cobra.Command, weight int, basename string) error {
	for _, c := range command.Commands() {
		if ignore(c) {
			continue
		}

		b := extractFilename(c.CommandPath())

		baseWeight++
		if err := generate(c, baseWeight, b); err != nil {
			return fmt.Errorf("generating docs for %q: %w", b, err)
		}
	}

	filename := filepath.Join("docs", "reference", basename+".md")

	f, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("creating file %q: %w", filename, err)
	}
	defer f.Close() // lint:allow_defer_close

	if _, err := f.WriteString(""); err != nil {
		return fmt.Errorf("writing to file %q: %w", filename, err)
	}

	if err := generateMarkdown(command, weight, f); err != nil {
		return fmt.Errorf("generating markdown for %q: %w", basename, err)
	}

	return nil
}

func generateMarkdown(command *cobra.Command, weight int, w io.Writer) error {
	command.InitDefaultHelpCmd()
	command.InitDefaultHelpFlag()

	commandPath := command.CommandPath()
	name := strings.ReplaceAll(commandPath, "terraform-docs ", "")

	short := command.Short
	long := command.Long

	if long == "" {
		long = short
	}

	parent := "reference"
	if command.Parent() != nil {
		parent = command.Parent().Name()
	}

	ref := &reference{
		Name:        name,
		Command:     commandPath,
		Description: short,
		Parent:      parent,
		Synopsis:    long,
		Runnable:    command.Runnable(),
		HasChildren: len(command.Commands()) > 0,
		UseLine:     command.UseLine(),
		Weight:      weight,
	}

	// Options.
	if f := command.NonInheritedFlags(); f.HasAvailableFlags() {
		ref.Options = f.FlagUsages()
	}

	// Inherited Options.
	if f := command.InheritedFlags(); f.HasAvailableFlags() {
		ref.InheritedOptions = f.FlagUsages()
	}

	if ref.HasChildren {
		subcommands(ref, command.Commands())
	} else {
		if err := example(ref); err != nil {
			return fmt.Errorf("generating example: %w", err)
		}
	}

	file := "format.tmpl"
	paths := []string{filepath.Join("scripts", "docs", file)}

	t := template.Must(template.New(file).ParseFiles(paths...))

	if err := t.Execute(w, ref); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}

func example(ref *reference) error {
	flag := " --footer-from footer.md"
	if ref.Name == "pretty" {
		flag += " --no-color"
	}

	ref.Usage = fmt.Sprintf("%s%s ./examples/", ref.Command, flag)

	config := print.DefaultConfig()

	config.ModuleRoot = "./examples"
	config.Formatter = ref.Name
	config.Settings.Color = false
	config.Sections.Show = append(config.Sections.Show, "all")
	config.Sections.Footer = true
	config.FooterFrom = "footer.md"
	config.Parse()

	tfmodule, err := terraform.LoadWithOptions(config)
	if err != nil {
		return fmt.Errorf("loading terraform module: %w", err)
	}

	formatter, err := format.New(config)
	if err != nil {
		return fmt.Errorf("creating formatter: %w", err)
	}

	if err := formatter.Generate(tfmodule); err != nil {
		return fmt.Errorf("generating example output: %w", err)
	}

	segments := strings.Split(formatter.Content(), "\n")
	buf := new(bytes.Buffer)

	for _, s := range segments {
		if s == "" {
			buf.WriteString("\n")
		} else {
			fmt.Fprintf(buf, "    %s\n", s)
		}
	}

	ref.Example = buf.String()

	return nil
}

func subcommands(ref *reference, children []*cobra.Command) {
	var subs []command

	for _, child := range children {
		if ignore(child) {
			continue
		}

		var subchild []command

		for _, c := range child.Commands() {
			if ignore(c) {
				continue
			}

			cname := c.CommandPath()
			link := extractFilename(cname)

			subchild = append(subchild, command{Name: cname, Link: link})
		}

		cname := child.CommandPath()
		link := extractFilename(cname)

		subs = append(subs, command{Name: cname, Link: link, Children: subchild})
	}

	ref.Subcommands = subs
}

func extractFilename(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "terraform-docs-", "")

	return s
}
