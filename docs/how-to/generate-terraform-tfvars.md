---
title: "Generate Terraform.tfvars"
description: "How to generate Terraform.tfvars file with Terraform-docs"
menu:
  docs:
    parent: "how-to"
weight: 208
toc: false
---

Since `v0.9.0`

You can generate `terraform.tfvars` in both `hcl` and `json` format by executing
the following, respectively:

```bash
terraform-docs tfvars hcl /path/to/module

terraform-docs tfvars json /path/to/module
```

{{< alert type="info" >}}
Required input variables will be `""` (empty) in HCL and `null` in JSON format.
{{< /alert >}}
