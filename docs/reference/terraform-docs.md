---
title: "Terraform-docs"
description: "A utility to generate documentation from Terraform modules in various output formats"
menu:
  docs:
    parent: "reference"
weight: 950
toc: true
---

## Synopsis

A utility to generate documentation from Terraform modules in various output formats.

```console
terraform-docs [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Subcommands

* [`terraform-docs asciidoc`]({{< ref "asciidoc" >}})
  * [`terraform-docs asciidoc document`]({{< ref "asciidoc-document" >}})
  * [`terraform-docs asciidoc table`]({{< ref "asciidoc-table" >}})
* [`terraform-docs json`]({{< ref "JSON" >}})
* [`terraform-docs markdown`]({{< ref "Markdown" >}})
  * [`terraform-docs markdown document`]({{< ref "Markdown-document" >}})
  * [`terraform-docs markdown table`]({{< ref "Markdown-table" >}})
* [`terraform-docs pretty`]({{< ref "pretty" >}})
* [`terraform-docs tfvars`]({{< ref "tfvars" >}})
  * [`terraform-docs tfvars hcl`]({{< ref "tfvars-hcl" >}})
  * [`terraform-docs tfvars json`]({{< ref "tfvars-JSON" >}})
* [`terraform-docs toml`]({{< ref "TOML" >}})
* [`terraform-docs xml`]({{< ref "XML" >}})
* [`terraform-docs yaml`]({{< ref "YAML" >}})
