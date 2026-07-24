---
title: "pretty"
description: "Generate colorized pretty of inputs and outputs"
menu:
  docs:
    parent: "Terraform-docs"
weight: 958
toc: true
---

## Synopsis

Generate colorized pretty of inputs and outputs.

```console
terraform-docs pretty [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs pretty --footer-from footer.md --no-color ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
