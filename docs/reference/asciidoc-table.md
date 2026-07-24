---
title: "asciidoc table"
description: "Generate AsciiDoc tables of inputs and outputs"
menu:
  docs:
    parent: "asciidoc"
weight: 953
toc: true
---

## Synopsis

Generate AsciiDoc tables of inputs and outputs.

```console
terraform-docs asciidoc table [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs asciidoc table --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
