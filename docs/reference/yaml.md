---
title: "YAML"
description: "Generate YAML of inputs and outputs"
menu:
  docs:
    parent: "Terraform-docs"
weight: 964
toc: true
---

## Synopsis

Generate YAML of inputs and outputs.

```console
terraform-docs yaml [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs yaml --footer-from footer.md ./examples/
```

[examples]: https://github.com/terraform-docs/terraform-docs/tree/main/examples
