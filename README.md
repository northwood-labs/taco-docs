# taco-docs

A utility to generate documentation from Terraform/OpenTofu modules in various output formats. A hard-fork of [terraform-docs](https://terraform-docs.io).

**TACOS:** _**T**erraform **A**utomation and **CO**laboration **S**oftware_.

Updated documentation coming soon.

## Why a fork?

By July 2026, `terraform-docs` no longer appeared to be meaningfully maintained. There were 38 open (and unreviewed) pull requests, and 153 open issues (most had no meaningful response). By all appearances, nobody was left to steer the ship.

So, we decided to hard-fork the project, merge a number of pull requests into our fork, and plan for longer-term maintenance over the project moving forward.

**Current status:** Major refactoring before cutting a new release.

## Current status

We hard-forked this project on Thursday, 2026-07-23. Many of the open PRs were **several dozen** commits behind. We did our best to merge them with accuracy, but we need to to a more thorough review of the code and the test suite to ensure that the merges applied as expected.

It is Northwood Labs’ intention to maintain this fork moving forward, and merge additional upstream pull requests that we feel are a good fit for this project and its users. We will also start chipping away at some of the open issues with new fixes.

### Merged PRs

Need to apply some manual touching-up here and verify appropriate test coverage.

* https://github.com/northwood-labs/taco-docs/pull/650
* https://github.com/northwood-labs/taco-docs/pull/657
* https://github.com/northwood-labs/taco-docs/pull/658
* https://github.com/northwood-labs/taco-docs/pull/725
* https://github.com/northwood-labs/taco-docs/pull/763
* https://github.com/northwood-labs/taco-docs/pull/847
* https://github.com/northwood-labs/taco-docs/pull/861
* https://github.com/northwood-labs/taco-docs/pull/866
* https://github.com/northwood-labs/taco-docs/pull/894
* https://github.com/northwood-labs/taco-docs/pull/897
* https://github.com/northwood-labs/taco-docs/pull/926
* https://github.com/northwood-labs/taco-docs/pull/936
* https://github.com/northwood-labs/taco-docs/pull/937
* https://github.com/northwood-labs/taco-docs/pull/938
* https://github.com/northwood-labs/taco-docs/pull/947

### To be manually merged

These are either too old, or diverge too much after the aforementioned merges were applied. These will need to be merged by-hand to get the fixes into the source code. We will do our best to preserve attribution for the original committer.

* https://github.com/northwood-labs/taco-docs/pull/571
* https://github.com/northwood-labs/taco-docs/pull/700
* https://github.com/northwood-labs/taco-docs/pull/709
* https://github.com/northwood-labs/taco-docs/pull/723
* https://github.com/northwood-labs/taco-docs/pull/820
* https://github.com/northwood-labs/taco-docs/pull/893
* https://github.com/northwood-labs/taco-docs/pull/905
* https://github.com/northwood-labs/taco-docs/pull/935

The original README has been temporarily removed, as much of it will change.
