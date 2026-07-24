---
title: "Markdown document"
description: "Generate Markdown document of inputs and outputs"
menu:
  docs:
    parent: "Markdown"
weight: 956
toc: true
---

## Synopsis

Generate Markdown document of inputs and outputs.

```console
terraform-docs markdown document [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs markdown document --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
