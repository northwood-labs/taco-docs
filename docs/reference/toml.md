---
title: "TOML"
description: "Generate TOML of inputs and outputs"
menu:
  docs:
    parent: "Terraform-docs"
weight: 962
toc: true
---

## Synopsis

Generate TOML of inputs and outputs.

```console
terraform-docs toml [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs toml --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
