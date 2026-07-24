---
title: "asciidoc document"
description: "Generate AsciiDoc document of inputs and outputs"
menu:
  docs:
    parent: "asciidoc"
weight: 952
toc: true
---

## Synopsis

Generate AsciiDoc document of inputs and outputs.

```console
terraform-docs asciidoc document [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs asciidoc document --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
