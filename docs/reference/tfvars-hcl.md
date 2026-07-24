---
title: "tfvars hcl"
description: "Generate HCL format of Terraform.tfvars of inputs"
menu:
  docs:
    parent: "tfvars"
weight: 960
toc: true
---

## Synopsis

Generate HCL format of Terraform.tfvars of inputs.

```console
terraform-docs tfvars hcl [PATH] [flags]
```

Use `--help` to see all available CLI flags.

## Example

Given the [`examples`][examples] module:

```shell
terraform-docs tfvars hcl --footer-from footer.md ./examples/
```

[examples]: https://github.com/northwood-labs/taco-docs/tree/main/examples
