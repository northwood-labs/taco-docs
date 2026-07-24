---
title: "Markdown table"
description: "Generate Markdown tables of inputs and outputs"
menu:
  docs:
    parent: "Markdown"
weight: 957
toc: true
---

## Synopsis

Generate Markdown tables of inputs and outputs.

```console
terraform-docs markdown table [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs markdown table --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
