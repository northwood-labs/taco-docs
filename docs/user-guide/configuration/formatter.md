---
title: "formatter"
description: "formatter configuration"
menu:
  docs:
    parent: "configuration"
weight: 123
toc: true
---

Since `v0.10.0`

The following options are supported out of the box by Terraform-docs and can be
used for `FORMATTER_NAME`:

* `asciidoc` — [reference]({{< ref "asciidoc" >}})
* `asciidoc document` — [reference]({{< ref "asciidoc-document" >}})
* `asciidoc table` — [reference]({{< ref "asciidoc-table" >}})
* `json` — [reference]({{< ref "JSON" >}})
* `markdown` — [reference]({{< ref "Markdown" >}})
* `markdown document` — [reference]({{< ref "Markdown-document" >}})
* `markdown table` — [reference]({{< ref "Markdown-table" >}})
* `pretty` — [reference]({{< ref "pretty" >}})
* `tfvars hcl` — [reference]({{< ref "tfvars-hcl" >}})
* `tfvars json` — [reference]({{< ref "tfvars-JSON" >}})
* `toml` — [reference]({{< ref "TOML" >}})
* `xml` — [reference]({{< ref "XML" >}})
* `yaml` — [reference]({{< ref "YAML" >}})

{{< alert type="info" >}}
Short version of formatters can also be used:

* `adoc` instead of `asciidoc`
* `md` instead of `markdown`
* `doc` instead of `document`
* `tbl` instead of `table`
{{< /alert >}}

{{< alert type="info" >}}
You need to pass name of a plugin as `formatter` in order to be able to
use the plugin. For example, if plugin binary file is called `tfdocs-format-foo`,
formatter name must be set to `foo`.
{{< /alert >}}

## Options

Available options with their default values.

```yaml
formatter: ""
```

{{< alert type="info" >}}
`formatter` is required and cannot be empty in `.terraform-docs.yml`.
{{< /alert >}}

## Examples

Format as Markdown table:

```yaml
formatter: "markdown table"
```

Format as Markdown document:

```yaml
formatter: "md doc"
```

Format as AsciiDoc document:

```yaml
formatter: "asciidoc document"
```

Format as `tfdocs-format-myplugin`:

```yaml
formatter: "myplugin"
```
