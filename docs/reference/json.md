---
title: "JSON"
description: "Generate JSON of inputs and outputs"
menu:
  docs:
    parent: "Terraform-docs"
weight: 954
toc: true
---

## Synopsis

Generate JSON of inputs and outputs.

```console
terraform-docs json [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs json --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
