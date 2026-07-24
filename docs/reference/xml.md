---
title: "XML"
description: "Generate XML of inputs and outputs"
menu:
  docs:
    parent: "Terraform-docs"
weight: 963
toc: true
---

## Synopsis

Generate XML of inputs and outputs.

```console
terraform-docs xml [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs xml --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
