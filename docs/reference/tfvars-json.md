---
title: "tfvars JSON"
description: "Generate JSON format of Terraform.tfvars of inputs"
menu:
  docs:
    parent: "tfvars"
weight: 961
toc: true
---

## Synopsis

Generate JSON format of Terraform.tfvars of inputs.

```console
terraform-docs tfvars json [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs tfvars json --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
