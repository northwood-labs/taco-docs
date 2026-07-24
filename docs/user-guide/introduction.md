---
title: "Introduction"
description: "Generate documentation from Terraform modules in various output formats"
menu:
  docs:
    parent: "user-guide"
weight: 100
toc: true
---

`terraform-docs` is a utility to generate documentation from Terraform modules in
various output formats.

{{< img-simple src="teaser.png" >}}

## Configuration

You can also have consistent execution through a `.terraform-docs.yml` file.

Once you set it up and configured it, every time you or your teammates want to
regenerate documentation (manually, through a pre-commit hook, or as part
of a CI pipeline) all you need to do is run `terraform-docs /module/path`.

{{< img-simple src="config.png" >}}

Read all about [configuration].

## Formats

One of the most popular format is [Markdown table], which is a very good fit for
generating README of module.

{{< img-simple src="Markdown-table.png" >}}

which produces:

{{< img-simple src="Markdown-table-output.png" >}}

Read all about available [formats].

## Compatibility

Terraform-docs compatibility matrix with Terraform can be found below:

| `terraform-docs` | Terraform         |
|------------------|-------------------|
| `>= 0.13`        | `>= 0.15`         |
| `>= 0.8, < 0.13` | `>= 0.12, < 0.15` |
| `< 0.8`          | `< 0.12`          |

[configuration]: {{< ref "configuration" >}}
[formats]: {{< ref "Terraform-docs" >}}
[Markdown table]: {{< ref "Markdown-table" >}}
