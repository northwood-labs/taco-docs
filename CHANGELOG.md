# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com), adheres to [Semantic Versioning](https://semver.org), and uses [Conventional Commit](https://www.conventionalcommits.org) syntax.

## Unreleased

[Compare: v0.24.0 → `HEAD`](https://github.com/northwood-labs/taco-docs/compare/v0.24.0..HEAD)

### :art: Styling

* [`5e7e589`](https://github.com/northwood-labs/taco-docs/commit/5e7e589bf9036d24d0b8ed3758bf52c5cb83903d): Run pre-commit ([@RemoteRabbit](https://github.com/RemoteRabbit))

### :books: Documentation

* [`903034b`](https://github.com/northwood-labs/taco-docs/commit/903034bd9a2b98caa29fd4e423f92a8fe30150bd): Clarify container image registry ([@Piyushkhobragade](https://github.com/Piyushkhobragade))
* [`4de00bf`](https://github.com/northwood-labs/taco-docs/commit/4de00bf2713c96b02caa9f87b67bae4079685634): Add taco-docs fork documentation with PR merge status. ([@skyzyx](https://github.com/skyzyx))
* [`63681a3`](https://github.com/northwood-labs/taco-docs/commit/63681a362ff473c6c1f9c5710a7b7f7c9ab05aaf): Add comprehensive architectural documentation with design rationale. ([@skyzyx](https://github.com/skyzyx))

### :tractor: Refactor

* [`34bef73`](https://github.com/northwood-labs/taco-docs/commit/34bef735beef5f64bb7e78b9a48ca85f8aa00cf4): Added first-pass at new project config files for Go projects. ([@skyzyx](https://github.com/skyzyx))

### <!-- 0 -->:rocket: Features

* [`67d48f8`](https://github.com/northwood-labs/taco-docs/commit/67d48f8a993eeb106de20b0fadda273b702978a0): Setup pre-commit hooks for devs ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`909d4d5`](https://github.com/northwood-labs/taco-docs/commit/909d4d598c5a1f67e9bcb4559bf47985293e9781): **terraform**: Added provider function support ([@pillatipriyanka](https://github.com/pillatipriyanka))
* [`cb0dcdb`](https://github.com/northwood-labs/taco-docs/commit/cb0dcdb4a68b543588179143a142e074ecc9e3b0): Sanitize bare URLs to Markdown format ([@liammoat](https://github.com/liammoat))

### <!-- 1 -->:bug: Bug Fixes

* [`9a52a55`](https://github.com/northwood-labs/taco-docs/commit/9a52a5553d3259be1a3cb0fd9865851512f542d9): Render explicit null variable defaults as null in tfvars hcl ([@somaz94](https://github.com/somaz94))
* [`8814c21`](https://github.com/northwood-labs/taco-docs/commit/8814c21f125867b7c7aca922a3ce1bb4e7b9bd05): **codespell**: Fix typo in readme, avaialble -> available ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`4515138`](https://github.com/northwood-labs/taco-docs/commit/45151388702cc9948dcaa5dbbb8c94e7c8a06f84): **codespell**: Fix typos in cmd/* ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`022e534`](https://github.com/northwood-labs/taco-docs/commit/022e5340a7eb4c133337a2ed0b7f12b3d2d2abb7): **codespell**: Fix typos in docs/reference, these got fixed when running the generate after the rest got fixed ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`0e6cd93`](https://github.com/northwood-labs/taco-docs/commit/0e6cd930319d9693d668aec6c7b232c3ceb6ca9b): **codespell**: Fix typos in format/ ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`03c9369`](https://github.com/northwood-labs/taco-docs/commit/03c9369de7ff05e35fb8862a9d1a179055cb9594): **codespell**: Fix typos in internal/ ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`dab33b0`](https://github.com/northwood-labs/taco-docs/commit/dab33b0f04dde5e8f54bf2d0523d108a407687cb): **codespell**: Fix typos in print/ ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`35df166`](https://github.com/northwood-labs/taco-docs/commit/35df1663b62e8b76c240ec58db430336fd11fd6e): **codespell**: Fix typos in template/ ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`81a7656`](https://github.com/northwood-labs/taco-docs/commit/81a76560af352622142f60f573b6d5005f7c6e21): **codespell**: Add codespell ignore to lines with `requireds` ([@RemoteRabbit](https://github.com/RemoteRabbit))
* [`5757539`](https://github.com/northwood-labs/taco-docs/commit/57575393c8d7ef51787e87c0762391123caae69c): Some formatters would auto remove whitespace from codeblock ([@RemoteRabbit](https://github.com/RemoteRabbit))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`b4dab24`](https://github.com/northwood-labs/taco-docs/commit/b4dab24f32e4f48db774949de3cb9ba05f320d2f): Update Go dependencies to latest versions. ([@skyzyx](https://github.com/skyzyx))
* [`18b7d51`](https://github.com/northwood-labs/taco-docs/commit/18b7d51a785925286ddb68d51e6201f56252183b): Update copyright headers and simplify code. ([@skyzyx](https://github.com/skyzyx))
* [`b642e29`](https://github.com/northwood-labs/taco-docs/commit/b642e29e3a7368447ad5c257baf9b8ac778b8d4b): Remove GitHub workflows and archive configuration files. ([@skyzyx](https://github.com/skyzyx))
* [`2a109bb`](https://github.com/northwood-labs/taco-docs/commit/2a109bbb0c7f2d17828a3cf7c359f309431554c7): Rebrand project as taco-docs. ([@skyzyx](https://github.com/skyzyx))

### Ops

* [`ea50b6a`](https://github.com/northwood-labs/taco-docs/commit/ea50b6a7aac20f4404ea6e517df0afa105e0ff58): **git**: Add .local/ to gitignore, often used for local repo notes ([@RemoteRabbit](https://github.com/RemoteRabbit))

## 0.24.0 — 2026-05-10

[Compare: v0.23.0 → v0.24.0](https://github.com/northwood-labs/taco-docs/compare/v0.23.0...v0.24.0)

### :dependabot: Building and Dependencies

* [`f9ca5a7`](https://github.com/northwood-labs/taco-docs/commit/f9ca5a7b95336353b18b8bd7babcd8c785f1adb1): **deps**: Bump `library/golang` from 1.26.2-alpine to 1.26.3-alpine ([#930](@REPO/issues/930)) ([@dependabot](https://github.com/dependabot))

### <!-- 1 -->:bug: Bug Fixes

* [`f5dca81`](https://github.com/northwood-labs/taco-docs/commit/f5dca81a13443c3072446169530c54d36ef0410e): Ignore directories without `.tf` files for docs generation ([#931](@REPO/issues/931)) ([@rudransh-shrivastava](https://github.com/rudransh-shrivastava))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`c080d5b`](https://github.com/northwood-labs/taco-docs/commit/c080d5bda4eb3594c6b7e46e816ec5fa50ffe458): Fix typo `genereted` ([#932](@REPO/issues/932)) ([@pascal-hofmann](https://github.com/pascal-hofmann))

## 0.23.0 — 2026-05-06

[Compare: v0.22.0 → v0.23.0](https://github.com/northwood-labs/taco-docs/compare/v0.22.0...v0.23.0)

### :dependabot: Building and Dependencies

* [`c355cbe`](https://github.com/northwood-labs/taco-docs/commit/c355cbe191c68d293dc16b4e55d0976ad301dbae): **deps**: Bump `stefanzweifel/git-auto-commit-action` from 5 to 7 ([@dependabot](https://github.com/dependabot))
* [`65c4048`](https://github.com/northwood-labs/taco-docs/commit/65c40482b9f72b7cbd5f5bbe30d8dc51cfd28e07): **deps**: Bump `codecov/codecov-action` from 4 to 5 ([@dependabot](https://github.com/dependabot))
* [`db8f626`](https://github.com/northwood-labs/taco-docs/commit/db8f626011c8268ca0d796ed7834dbfe88d12738): **deps**: Bump `library/alpine` from 3.23.2 to 3.23.3 ([@dependabot](https://github.com/dependabot))
* [`0ff2b14`](https://github.com/northwood-labs/taco-docs/commit/0ff2b14170bcd0efde1b80ce5b9beb5c9e61d88a): **deps**: Bump `goreleaser/goreleaser-action` from 6 to 7 ([@dependabot](https://github.com/dependabot))
* [`31b29ab`](https://github.com/northwood-labs/taco-docs/commit/31b29aba257c5dc94a43952fecc9e38b281a86a2): **deps**: Bump `actions/checkout` from 4 to 6 ([@dependabot](https://github.com/dependabot))
* [`c3618ef`](https://github.com/northwood-labs/taco-docs/commit/c3618ef2c4b26a5f03a1baf1feabb37d1ba6ec23): **deps**: Bump `codecov/codecov-action` from 5 to 6 ([@dependabot](https://github.com/dependabot))
* [`747fbde`](https://github.com/northwood-labs/taco-docs/commit/747fbde7a0f305183c8f360c735522346ff3911c): **deps**: Bump `docker/setup-qemu-action` from 3 to 4 ([@dependabot](https://github.com/dependabot))
* [`0ed4b7d`](https://github.com/northwood-labs/taco-docs/commit/0ed4b7d399a61e31dffa67449676249b603d7149): **deps**: Bump `docker/login-action` from 3 to 4 ([@dependabot](https://github.com/dependabot))
* [`a793f92`](https://github.com/northwood-labs/taco-docs/commit/a793f920ff07c650eb15783710da6075a29ffff7): **deps**: Bump `library/golang` from 1.25.5-alpine to 1.26.0-alpine ([@dependabot](https://github.com/dependabot))
* [`e891fc1`](https://github.com/northwood-labs/taco-docs/commit/e891fc1ccd6556b5bfd9d20d10f69a1d74db6551): **deps**: Bump `docker/build-push-action` from 6 to 7 ([@dependabot](https://github.com/dependabot))
* [`105608d`](https://github.com/northwood-labs/taco-docs/commit/105608d46d432b58104678eb1a5f3a076327fd86): **deps**: Bump `actions/setup-go` from 5 to 6 ([@dependabot](https://github.com/dependabot))
* [`9af6379`](https://github.com/northwood-labs/taco-docs/commit/9af6379d9cc22b2d9dde8cc8b0e33738014fa899): **deps**: Bump `github/codeql-action` from 3 to 4 ([@dependabot](https://github.com/dependabot))
* [`91b7995`](https://github.com/northwood-labs/taco-docs/commit/91b79959905200e40656453cb2a342ea3f3f09de): **deps**: Bump `library/alpine` in /scripts/release ([@dependabot](https://github.com/dependabot))
* [`851e507`](https://github.com/northwood-labs/taco-docs/commit/851e5072f06c56bd6149f7472dba89e0d05e4a79): **deps**: Bump `library/golang` from 1.26.0-alpine to 1.26.2-alpine ([@dependabot](https://github.com/dependabot))
* [`26553fe`](https://github.com/northwood-labs/taco-docs/commit/26553feb284e8abee65416c908fdddf8adde7c48): **deps**: Bump `softprops/action-gh-release` from 2 to 3 ([@dependabot](https://github.com/dependabot))
* [`6a413c3`](https://github.com/northwood-labs/taco-docs/commit/6a413c34e1b820ae3ade4d39e2f9b286f7e1be2d): **deps**: Bump `docker/setup-buildx-action` from 3 to 4 ([@dependabot](https://github.com/dependabot))
* [`be1a794`](https://github.com/northwood-labs/taco-docs/commit/be1a79478b65f008aeb42f497ba9a862bcd8dbc9): **deps**: Bump `library/alpine` from 3.23.3 to 3.23.4 ([#925](@REPO/issues/925)) ([@dependabot](https://github.com/dependabot))

### :tractor: Refactor

* [`3cc1f37`](https://github.com/northwood-labs/taco-docs/commit/3cc1f37cf815f66a15e640b06f78747862e18e1e): Address golangci-lint cyclomatic complexity error ([#928](@REPO/issues/928)) ([@pascal-hofmann](https://github.com/pascal-hofmann))

### <!-- 0 -->:rocket: Features

* [`ea95325`](https://github.com/northwood-labs/taco-docs/commit/ea953255968baed053169b35382dc89ad944f241): **cli**: Improve recursive submodule scanning to traverse nested directories ([@rudransh-shrivastava](https://github.com/rudransh-shrivastava))
* [`eae174a`](https://github.com/northwood-labs/taco-docs/commit/eae174a6fbf0c8c1bda502cd9ee140a7cf7db359): **recursive**: Add exclude config option and CLI flag to skip directories during recursive scanning ([@rudransh-shrivastava](https://github.com/rudransh-shrivastava))

### <!-- 1 -->:bug: Bug Fixes

* [`4e8d760`](https://github.com/northwood-labs/taco-docs/commit/4e8d760ae458e947c6782570c30796541b53526f): Config load with absolute path ([@gusandrioli](https://github.com/gusandrioli))
* [`d91daf9`](https://github.com/northwood-labs/taco-docs/commit/d91daf920cfbfff2a54855515227d340d698322d): Deterministic sorting for providers ([@gusandrioli](https://github.com/gusandrioli))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`ac64df2`](https://github.com/northwood-labs/taco-docs/commit/ac64df2db7b9ff0e21628e0994f1119b4e596acf): **release**: Do not bump homebrew version ([@pascal-hofmann](https://github.com/pascal-hofmann))
* [`0ec4c49`](https://github.com/northwood-labs/taco-docs/commit/0ec4c49af6ea3c5e15b5e27eb09a1fea4ba739b2): **codecov**: Use parameter files instead of deprecated file ([@pascal-hofmann](https://github.com/pascal-hofmann))
* [`1e3f6a1`](https://github.com/northwood-labs/taco-docs/commit/1e3f6a12e8c5c5108c88322084af024439112af2): Testing for sort config enabled ([@gusandrioli](https://github.com/gusandrioli))

## 0.22.0 — 2026-04-07

[Compare: v0.21.0 → v0.22.0](https://github.com/northwood-labs/taco-docs/compare/v0.21.0...v0.22.0)

### :books: Documentation

* [`5b0c7d9`](https://github.com/northwood-labs/taco-docs/commit/5b0c7d9d344b3adc96a54201fb6b6a9316114d86): Add OpenSSF Best Practices badge ([@pascal-hofmann](https://github.com/pascal-hofmann))

### :dependabot: Building and Dependencies

* [`ee533f9`](https://github.com/northwood-labs/taco-docs/commit/ee533f9dc746a039682e57a1e04dc2da89d85d77): **deps**: Bump `library/golang` from 1.24.2-alpine to 1.25.3-alpine ([@dependabot](https://github.com/dependabot))
* [`d9bfd39`](https://github.com/northwood-labs/taco-docs/commit/d9bfd39e740021e14fb9c5ff2351e1f5df4ecc7d): **deps**: Bump `library/alpine` from 3.22.2 to 3.23.0 ([@dependabot](https://github.com/dependabot))
* [`796a38d`](https://github.com/northwood-labs/taco-docs/commit/796a38d2f0bbc0cf0e4b7a2d95514873c14a232a): **deps**: Bump `library/alpine` in /scripts/release ([@dependabot](https://github.com/dependabot))
* [`dd2fe9d`](https://github.com/northwood-labs/taco-docs/commit/dd2fe9dc575129d4d58c3c87affee68a26e8c0a5): **deps**: Update Go and dependencies to latest versions ([@ricardo-kh](https://github.com/ricardo-kh))
* [`15a95e6`](https://github.com/northwood-labs/taco-docs/commit/15a95e671eae3650393f06f00516b90a7b8619fb): **deps**: Bump `library/alpine` in /scripts/release ([@dependabot](https://github.com/dependabot))
* [`ad868f2`](https://github.com/northwood-labs/taco-docs/commit/ad868f2f4324a1fdf996674198f771be0c3f4456): **deps**: Bump `library/alpine` from 3.23.0 to 3.23.2 ([@dependabot](https://github.com/dependabot))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`67e416b`](https://github.com/northwood-labs/taco-docs/commit/67e416b2f2491f23f65e20e7bab5e2d9d39e1ebe): Fix tests ([@pascal-hofmann](https://github.com/pascal-hofmann))

## 0.21.0 — 2025-12-12

[Compare: v0.20.0 → v0.21.0](https://github.com/northwood-labs/taco-docs/compare/v0.20.0...v0.21.0)

### :books: Documentation

* [`e3fe6f1`](https://github.com/northwood-labs/taco-docs/commit/e3fe6f123c1662f29c5861d2f5e32b78db496391): Update links to go wiki in `CONTRIBUTING.md` ([@pascal-hofmann](https://github.com/pascal-hofmann))

### :dependabot: Building and Dependencies

* [`553f6ad`](https://github.com/northwood-labs/taco-docs/commit/553f6ad72c138e70126fafffd4f301f89e92cece): **deps**: Bump `golang.org/x/net` from 0.36.0 to 0.38.0 ([@dependabot](https://github.com/dependabot))
* [`929b3e0`](https://github.com/northwood-labs/taco-docs/commit/929b3e0fffc92755eb21170535b71c952e45a506): **deps**: Bump `dependencies,` golang version and base containers ([@pascal-hofmann](https://github.com/pascal-hofmann))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`4ac0a6a`](https://github.com/northwood-labs/taco-docs/commit/4ac0a6a5243452eeedb4861256c928e1c664115e): **linting**: Update to golangci-lint 2.7.2 ([@pascal-hofmann](https://github.com/pascal-hofmann))
* [`dd450a5`](https://github.com/northwood-labs/taco-docs/commit/dd450a5352e7a7571d28c74b9674e3eabce6a681): Fix golangci-lint issues ([@pascal-hofmann](https://github.com/pascal-hofmann))

## 0.20.0 — 2025-04-04

[Compare: v0.19.0 → v0.20.0](https://github.com/northwood-labs/taco-docs/compare/v0.19.0...v0.20.0)

### :dependabot: Building and Dependencies

* [`616bff0`](https://github.com/northwood-labs/taco-docs/commit/616bff0feb6b47ff931a77e3c705fd02bfb0b3ea): **deps**: Bump `library/golang` from 1.23.1-alpine to 1.23.4-alpine ([@dependabot](https://github.com/dependabot))
* [`006ff31`](https://github.com/northwood-labs/taco-docs/commit/006ff31f25d6288d88cbeef04ca02139fd90a1c7): **deps**: Bump `golang.org/x/crypto` from 0.27.0 to 0.31.0 ([@dependabot](https://github.com/dependabot))
* [`e470746`](https://github.com/northwood-labs/taco-docs/commit/e4707464498082547cf1df2756000f55c215eaf1): **deps**: Bump `golang.org/x/net` from 0.29.0 to 0.33.0 ([@dependabot](https://github.com/dependabot))
* [`73ee296`](https://github.com/northwood-labs/taco-docs/commit/73ee2961b456dac63debd47410d3694d0d348673): **deps**: Bump `library/alpine` from 3.20.3 to 3.21.3 ([@dependabot](https://github.com/dependabot))
* [`06ca95c`](https://github.com/northwood-labs/taco-docs/commit/06ca95c6b7c3c498364c9b168522904b761fb758): **deps**: Bump `library/alpine` in /scripts/release ([@dependabot](https://github.com/dependabot))
* [`adb8099`](https://github.com/northwood-labs/taco-docs/commit/adb8099f144447b961270e0b4245de597b4dd6b7): **deps**: Bump `golang.org/x/net` from 0.33.0 to 0.36.0 ([@dependabot](https://github.com/dependabot))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`55d8916`](https://github.com/northwood-labs/taco-docs/commit/55d89169488edb3bb16eb5960b723bca6f729f4d): Bump `version` to v0.20.0-alpha ([@tf-docs-bot](https://github.com/tf-docs-bot))
* [`983e98a`](https://github.com/northwood-labs/taco-docs/commit/983e98a7b561e0f7065cc7012749bf6534dfded0): Bump `golang` to 1.24.2 ([@pascal-hofmann](https://github.com/pascal-hofmann))
* [`93c1839`](https://github.com/northwood-labs/taco-docs/commit/93c18398836d7d75417c92ae4bf655b2bdb51961): Update staticcheck to 2025.1.1 ([@pascal-hofmann](https://github.com/pascal-hofmann))

## 0.19.0 — 2024-09-18

[Compare: v0.18.0 → v0.19.0](https://github.com/northwood-labs/taco-docs/compare/v0.18.0...v0.19.0)

### :dependabot: Building and Dependencies

* [`c825b41`](https://github.com/northwood-labs/taco-docs/commit/c825b41689069d9d51d53bd459a8164a5a3290d6): **deps**: Bump `library/alpine` from 3.20.0 to 3.20.2 ([@dependabot](https://github.com/dependabot))
* [`a639fbd`](https://github.com/northwood-labs/taco-docs/commit/a639fbd42479a6e76c8231b8f9d3f6049ed6209f): **deps**: Bump `library/alpine` in /scripts/release ([@dependabot](https://github.com/dependabot))
* [`5441df2`](https://github.com/northwood-labs/taco-docs/commit/5441df2ac976dcda3dd8273adf8559aec1728975): **deps**: Bump `library/alpine` from 3.20.2 to 3.20.3 ([@dependabot](https://github.com/dependabot))
* [`7da557a`](https://github.com/northwood-labs/taco-docs/commit/7da557ac5eeb9f8b537a15f49215b17ba253b5d7): **deps**: Bump `docker/build-push-action` from 5 to 6 ([@dependabot](https://github.com/dependabot))
* [`a2f4573`](https://github.com/northwood-labs/taco-docs/commit/a2f4573244b63716ff3c62ebcbee022347b0e98d): **deps**: Bump `library/golang` from 1.23.0-alpine to 1.23.1-alpine ([@dependabot](https://github.com/dependabot))

### <!-- 1 -->:bug: Bug Fixes

* [`78e94da`](https://github.com/northwood-labs/taco-docs/commit/78e94da7864d47732bbd941fe5c4e081ccd77e74): Replace <br> with <br /> for markdown syntax ([@christophe-scalepad](https://github.com/christophe-scalepad))
* [`af31cc6`](https://github.com/northwood-labs/taco-docs/commit/af31cc618e8f4763da3323de4040927d849041de): Release scripts ([@khos2ow](https://github.com/khos2ow))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`a97e171`](https://github.com/northwood-labs/taco-docs/commit/a97e1713117275d3b1a1407605c2a7d3a67c5e90): Kickoff actions run ([@khos2ow](https://github.com/khos2ow))
* [`8ae3344`](https://github.com/northwood-labs/taco-docs/commit/8ae3344d97da87cf80209ed2e79c632e59e2c5f1): Bump `version` to v0.19.0-alpha ([@khos2ow](https://github.com/khos2ow))
* [`1919452`](https://github.com/northwood-labs/taco-docs/commit/19194525e499d1a32627fe90ead37951d02687d0): Enhance release workflows ([@khos2ow](https://github.com/khos2ow))
* [`186bd7e`](https://github.com/northwood-labs/taco-docs/commit/186bd7e667861ce7e8f9fcad4de74c7c08b8d756): Update teaser image ([@khos2ow](https://github.com/khos2ow))
* [`4c94478`](https://github.com/northwood-labs/taco-docs/commit/4c944787bc56574ea0a9d63f483f41e635530398): Use correct env var for repo owner ([@pascal-hofmann](https://github.com/pascal-hofmann))
* [`62756ca`](https://github.com/northwood-labs/taco-docs/commit/62756ca4d3de82b63de9b74635a473667bc104c8): Use correct var for repo owner (second try) ([@pascal-hofmann](https://github.com/pascal-hofmann))
* [`0db6eef`](https://github.com/northwood-labs/taco-docs/commit/0db6eef258c38ac6b37a33ebd1d94ad1e094e0d6): Update go to 1.23.1 ([@khos2ow](https://github.com/khos2ow))
* [`3355644`](https://github.com/northwood-labs/taco-docs/commit/3355644e986a830d22e4148a1f54f3dc8db99aeb): Update go dependencies ([@khos2ow](https://github.com/khos2ow))
* [`c2e8d0a`](https://github.com/northwood-labs/taco-docs/commit/c2e8d0ac3c1428fd62d37c8fbc0e4997fcb953be): Fix linter issues ([@khos2ow](https://github.com/khos2ow))
* [`11270e3`](https://github.com/northwood-labs/taco-docs/commit/11270e31d8952d699dd48059f2b8593b294b78c3): Update staticcheck to 2024.1.1 ([@khos2ow](https://github.com/khos2ow))
* [`49fde02`](https://github.com/northwood-labs/taco-docs/commit/49fde02ef2fb556c0c336adcf7333d39b788f535): Fix release scripts ([@khos2ow](https://github.com/khos2ow))

### Fix

* [`3c44c58`](https://github.com/northwood-labs/taco-docs/commit/3c44c5828a0ce1c3a47c8419f5af808592276b41): Let Docker image be built correctly for non-amd64 platforms ([@Tenzer](https://github.com/Tenzer))

## 0.18.0 — 2024-05-30

[Compare: v0.17.0 → v0.18.0](https://github.com/northwood-labs/taco-docs/compare/v0.17.0...v0.18.0)

### :books: Documentation

* [`3f4630c`](https://github.com/northwood-labs/taco-docs/commit/3f4630c0be51581d3aa3a4303baa3e457e74e48d): Document `terraform-docs` + `ohmyzsh` usage ([@darkandrew7](https://github.com/darkandrew7))

### :test_tube: Testing

* [`fa916db`](https://github.com/northwood-labs/taco-docs/commit/fa916db457704e8c7bc6400540766f04dd5048ec): Amend full-example/main.tf, add local value with provider function to verify terraform-docs can properly parse ([@brittandeyoung](https://github.com/brittandeyoung))

### <!-- 0 -->:rocket: Features

* [`8f74fd4`](https://github.com/northwood-labs/taco-docs/commit/8f74fd445301187ee96d31c396cad42a11f921bb): Ignore outputs, providers, resources with comments ([@khos2ow](https://github.com/khos2ow))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`d5e48a5`](https://github.com/northwood-labs/taco-docs/commit/d5e48a56a9bb0d568979efb4847b017c83d145f2): Add scripts/release/ folder to dependabot ([@khos2ow](https://github.com/khos2ow))
* [`b79a7c4`](https://github.com/northwood-labs/taco-docs/commit/b79a7c485ea90948184b520f869d2747a7e46749): Bump `golang` to 1.22.3 ([@khos2ow](https://github.com/khos2ow))
* [`fde40b1`](https://github.com/northwood-labs/taco-docs/commit/fde40b18dbe342c8bddb98940a7ea9ad7bcdaa90): Bump `alpine` to 3.20.0 ([@khos2ow](https://github.com/khos2ow))
* [`288faea`](https://github.com/northwood-labs/taco-docs/commit/288faea0e5717d3ce3dc3324c08a3d366609c183): Update dependencies ([@khos2ow](https://github.com/khos2ow))
* [`2b71b4d`](https://github.com/northwood-labs/taco-docs/commit/2b71b4da7e52f30e13390e276531fccd6f5a21ab): Add release-cut workflow ([@khos2ow](https://github.com/khos2ow))

## 0.17.0 — 2023-12-19

[Compare: v0.16.0 → v0.17.0](https://github.com/northwood-labs/taco-docs/compare/v0.16.0...v0.17.0)

### :books: Documentation

* [`5d39661`](https://github.com/northwood-labs/taco-docs/commit/5d3966130fa16f15d6dfeb5a89f77c1b05cf4c5b): Fix two typos in the 'configuration' section of the documentation ([@tbriot](https://github.com/tbriot))

## 0.16.0 — 2021-10-05

[Compare: v0.15.0 → v0.16.0](https://github.com/northwood-labs/taco-docs/compare/v0.15.0...v0.16.0)

### :books: Documentation

* [`519f25e`](https://github.com/northwood-labs/taco-docs/commit/519f25ee01a6d0bd9989e1c59261d3d88f64194e): Fix typo in `README.md` and update configuration guide ([@bcdady](https://github.com/bcdady))

### <!-- 0 -->:rocket: Features

* [`045707b`](https://github.com/northwood-labs/taco-docs/commit/045707beee6423d5298a27f8c06b421a9dfba812): Add new flag 'read-comments' to optionally process comments as description ([@pbikki](https://github.com/pbikki))

## 0.11.0 — 2021-02-10

[Compare: v0.10.1 → v0.11.0](https://github.com/northwood-labs/taco-docs/compare/v0.10.1...v0.11.0)

### :books: Documentation

* [`3d8097e`](https://github.com/northwood-labs/taco-docs/commit/3d8097e2d524df70465f2b8dc2b28e120dcae772): Add missing `v` to tag ([#329](@REPO/issues/329)) ([@winmillwill](https://github.com/winmillwill))

### :dependabot: Building and Dependencies

* [`e3e81c2`](https://github.com/northwood-labs/taco-docs/commit/e3e81c25214bf998d087646d9535c718a324658f): Bump `github.com/spf13/cobra` from 1.0.0 to 1.1.0 ([#335](@REPO/issues/335)) ([@dependabot](https://github.com/dependabot))
* [`329ef42`](https://github.com/northwood-labs/taco-docs/commit/329ef42c2ee5328cd70e32774b3ac01f7c3f71de): Bump `github.com/hashicorp/hcl/v2` from 2.6.0 to 2.7.0 ([#334](@REPO/issues/334)) ([@dependabot](https://github.com/dependabot))
* [`5433d6b`](https://github.com/northwood-labs/taco-docs/commit/5433d6b8484b9e7db04353e2ac57fbdaa244a812): Bump `gopkg.in/yaml.v3` to v3.0.0-20200615113413-eeeca48fe776 ([#337](@REPO/issues/337)) ([@khos2ow](https://github.com/khos2ow))
* [`5bc452a`](https://github.com/northwood-labs/taco-docs/commit/5bc452a37dd29e4b8a6a1f6e1e8de630f86c7bc5): Bump `github.com/spf13/cobra` from 1.1.0 to 1.1.1 ([#340](@REPO/issues/340)) ([@dependabot](https://github.com/dependabot))
* [`50595de`](https://github.com/northwood-labs/taco-docs/commit/50595de0f253ab5933a75c967d2ab5ee36a0b5c0): Bump `github.com/zclconf/go-cty` from 1.6.1 to 1.7.0 ([#341](@REPO/issues/341)) ([@dependabot](https://github.com/dependabot))
* [`3663e92`](https://github.com/northwood-labs/taco-docs/commit/3663e92721dd5eea5cb0ffcc055896d17d3f40dd): Bump `alpine` from 3.12.0 to 3.12.1 ([#342](@REPO/issues/342)) ([@dependabot](https://github.com/dependabot))
* [`2eab30c`](https://github.com/northwood-labs/taco-docs/commit/2eab30cc2b35d7c8a16d05ee31d6947fbd8f9ce3): Bump `github.com/hashicorp/hcl/v2` from 2.7.0 to 2.8.0 ([#353](@REPO/issues/353)) ([@dependabot](https://github.com/dependabot))
* [`a0050f8`](https://github.com/northwood-labs/taco-docs/commit/a0050f8e726b745b990d02e000631517b5b3912d): Bump `github.com/zclconf/go-cty` from 1.7.0 to 1.7.1 ([#357](@REPO/issues/357)) ([@dependabot](https://github.com/dependabot))
* [`7be58e6`](https://github.com/northwood-labs/taco-docs/commit/7be58e6f0cbbff0213d8813b51488aeff3dcfa6b): Bump `github.com/hashicorp/hcl/v2` from 2.8.0 to 2.8.1 ([#359](@REPO/issues/359)) ([@dependabot](https://github.com/dependabot))
* [`e5e7302`](https://github.com/northwood-labs/taco-docs/commit/e5e7302c14c558f0243983ae180e8d118a3e340e): Bump `alpine` from 3.12.1 to 3.12.3 ([#358](@REPO/issues/358)) ([@dependabot](https://github.com/dependabot))
* [`f653336`](https://github.com/northwood-labs/taco-docs/commit/f653336783ffcd923fc309e3da7fb469dc38cd60): Bump `dawidd6/action-homebrew-bump-formula` from v3.4.1 to v3.5.0 ([#356](@REPO/issues/356)) ([@dependabot](https://github.com/dependabot))
* [`55cf46a`](https://github.com/northwood-labs/taco-docs/commit/55cf46abd5aa17065ac3d90de2579009db9a4c99): Bump `dawidd6/action-homebrew-bump-formula` from v3.5.0 to v3.5.1 ([#361](@REPO/issues/361)) ([@dependabot](https://github.com/dependabot))
* [`aa89342`](https://github.com/northwood-labs/taco-docs/commit/aa89342f7f4f2e71f366a7d704925a03cdf524f3): Bump `golang` from 1.15.2-alpine to 1.15.6-alpine ([#352](@REPO/issues/352)) ([@dependabot](https://github.com/dependabot))

### <!-- 1 -->:bug: Bug Fixes

* [`4cd6f59`](https://github.com/northwood-labs/taco-docs/commit/4cd6f59f1986d7991d6eb780716df749b3106071): Show correct version when brew installs it ([#332](@REPO/issues/332)) ([@khos2ow](https://github.com/khos2ow))
* [`63750c1`](https://github.com/northwood-labs/taco-docs/commit/63750c1784493b1d9e5a3059bef40325b2e501b7): Normalize last empty line of the generated output ([#336](@REPO/issues/336)) ([@khos2ow](https://github.com/khos2ow))
* [`4a98297`](https://github.com/northwood-labs/taco-docs/commit/4a9829783739cbd8f0e3bb28aad22187c77c5382): Cleanup extra empty lines from 'pretty' output ([#338](@REPO/issues/338)) ([@khos2ow](https://github.com/khos2ow))
* [`2353afb`](https://github.com/northwood-labs/taco-docs/commit/2353afbcae40b447cb79887c2a581bbb4fd4d443): Never escape special characters in tfvars json ([#339](@REPO/issues/339)) ([@khos2ow](https://github.com/khos2ow))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`6cdab20`](https://github.com/northwood-labs/taco-docs/commit/6cdab202144bbbe2241faaac620b4d02fa28aaf8): Disable dependabot ([#363](@REPO/issues/363)) ([@khos2ow](https://github.com/khos2ow))

## 0.10.1 — 2020-09-28

[Compare: v0.10.0 → v0.10.1](https://github.com/northwood-labs/taco-docs/compare/v0.10.0...v0.10.1)

### <!-- 1 -->:bug: Bug Fixes

* [`215ab95`](https://github.com/northwood-labs/taco-docs/commit/215ab95b6921d9bd0513b8f7416c584fa42bd37a): Prevent segfault error if input arg is a file ([#327](@REPO/issues/327)) ([@khos2ow](https://github.com/khos2ow))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`dac90f1`](https://github.com/northwood-labs/taco-docs/commit/dac90f1948558a90486c0e42140ffee8fd372303): Enhance release scripts ([#326](@REPO/issues/326)) ([@khos2ow](https://github.com/khos2ow))
* [`ec40648`](https://github.com/northwood-labs/taco-docs/commit/ec40648e92f53c6d4290dd8f9b1cd5f93ae23bc8): Fix homebrew PR creation process ([#328](@REPO/issues/328)) ([@khos2ow](https://github.com/khos2ow))

## 0.10.0 — 2020-09-21

[Compare: v0.9.1 → v0.10.0](https://github.com/northwood-labs/taco-docs/compare/v0.9.1...v0.10.0)

### :books: Documentation

* [`6e259ba`](https://github.com/northwood-labs/taco-docs/commit/6e259baf877ea06911c25a32df5b755e2c495d9d): Add detail about module header usage guide ([#282](@REPO/issues/282)) ([@khos2ow](https://github.com/khos2ow))
* [`f71c2a3`](https://github.com/northwood-labs/taco-docs/commit/f71c2a360cba342c9d968e1eeab4729bf70af308): Overall improvements to documentation ([#293](@REPO/issues/293)) ([@khos2ow](https://github.com/khos2ow))

### :dependabot: Building and Dependencies

* [`f060a2d`](https://github.com/northwood-labs/taco-docs/commit/f060a2df4d539da0cdfd177a34ca31eb6c8d28f1): **deps**: Bump `github.com/zclconf/go-cty` from 1.3.1 to 1.4.0 ([#239](@REPO/issues/239)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`7b68f3c`](https://github.com/northwood-labs/taco-docs/commit/7b68f3cf0e7d5d8ec32ccd54cbe252686223fd6b): **deps**: Bump `github.com/hashicorp/hcl/v2` from 2.3.0 to 2.4.0 ([#247](@REPO/issues/247)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`29d3c52`](https://github.com/northwood-labs/taco-docs/commit/29d3c523dbe6c6cc7bd2b934246561cbf2af520a): **deps**: Bump `mvdan.cc/xurls/v2` from 2.1.0 to 2.2.0 ([#248](@REPO/issues/248)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`5aaf499`](https://github.com/northwood-labs/taco-docs/commit/5aaf499280a13903916b9e40ed4869e616bf34b2): **deps**: Bump `github.com/go-test/deep` from 1.0.5 to 1.0.6 ([#250](@REPO/issues/250)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`6034c10`](https://github.com/northwood-labs/taco-docs/commit/6034c1034be5987c06e8c9735f83351429ae139b): **deps**: Bump `github.com/hashicorp/hcl/v2` from 2.4.0 to 2.5.0 ([#254](@REPO/issues/254)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`8d9c54f`](https://github.com/northwood-labs/taco-docs/commit/8d9c54f1dfda76592a5fac50af37735b99f5243f): **deps**: Bump `github.com/spf13/cobra` from 0.0.7 to 1.0.0 ([#261](@REPO/issues/261)) ([@khos2ow](https://github.com/khos2ow))
* [`342db66`](https://github.com/northwood-labs/taco-docs/commit/342db66d9625c13005a61afab6bbefa8b7107ce6): **deps**: Bump `gopkg.in/yaml.v2` from 2.2.8 to 3.0.0 ([#260](@REPO/issues/260)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`6ed89d2`](https://github.com/northwood-labs/taco-docs/commit/6ed89d27696ac45f56de0048962cc7ac37b190d4): **deps**: Bump `github.com/hashicorp/hcl/v2` from 2.5.0 to 2.5.1 ([#262](@REPO/issues/262)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`027e6df`](https://github.com/northwood-labs/taco-docs/commit/027e6dfac604657de1394b2196f60d0da2e13ce8): **deps**: Bump `github.com/zclconf/go-cty` from 1.4.0 to 1.4.1 ([#263](@REPO/issues/263)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`36be491`](https://github.com/northwood-labs/taco-docs/commit/36be49175a824187691d81357d70946d627b6ea3): **deps**: Bump `github.com/stretchr/testify` from 1.5.1 to 1.6.0 ([#268](@REPO/issues/268)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`cd09faf`](https://github.com/northwood-labs/taco-docs/commit/cd09faf4c4dde78571984942e7c6cca1e07aa237): **deps**: Bump `github.com/zclconf/go-cty` from 1.4.1 to 1.4.2 ([#271](@REPO/issues/271)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`8f2e271`](https://github.com/northwood-labs/taco-docs/commit/8f2e271754111cc4f0eb843f476c36307eeb81b9): **deps**: Bump `github.com/hashicorp/hcl/v2` from 2.5.1 to 2.6.0 ([#273](@REPO/issues/273)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`cd912b1`](https://github.com/northwood-labs/taco-docs/commit/cd912b1d3f279fe81f619edb739aaa0343bc156c): **deps**: Bump `github.com/stretchr/testify` from 1.6.0 to 1.6.1 ([#274](@REPO/issues/274)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`2720f2c`](https://github.com/northwood-labs/taco-docs/commit/2720f2c35b052504d6cd881706e9a2aa144f9964): **deps**: Bump `github.com/zclconf/go-cty` from 1.4.2 to 1.5.0 ([#277](@REPO/issues/277)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`9357d77`](https://github.com/northwood-labs/taco-docs/commit/9357d775f1927df18bcea7a2138a5d10c777cb06): **deps**: Bump `github.com/zclconf/go-cty` from 1.5.0 to 1.5.1 ([#280](@REPO/issues/280)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`04375a6`](https://github.com/northwood-labs/taco-docs/commit/04375a6c6a2eb7b2faf5a99db658c63a986abc1c): **deps**: Bump `github.com/go-test/deep` from 1.0.6 to 1.0.7 ([#290](@REPO/issues/290)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`bf8da72`](https://github.com/northwood-labs/taco-docs/commit/bf8da72301a2ab728fc619c26b80316359e130da): **deps**: Bump `github.com/imdario/mergo` from 0.3.9 to 0.3.10 ([#295](@REPO/issues/295)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`c531384`](https://github.com/northwood-labs/taco-docs/commit/c5313844551ab40ef5168ad5de457b0c81ebb040): **deps**: Bump `github.com/imdario/mergo` from 0.3.10 to 0.3.11 ([#306](@REPO/issues/306)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`e9cd8aa`](https://github.com/northwood-labs/taco-docs/commit/e9cd8aa77475541abfe338909e1850cfaeee4202): **deps**: Bump `github.com/zclconf/go-cty` from 1.5.1 to 1.6.0 ([#316](@REPO/issues/316)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`5fcf555`](https://github.com/northwood-labs/taco-docs/commit/5fcf555e1c58da04ab9b5000ceab6be3cfe3697c): Allow dependabot to update github action versions ([#310](@REPO/issues/310)) ([@jlosito](https://github.com/jlosito))
* [`59c157c`](https://github.com/northwood-labs/taco-docs/commit/59c157c8184f56b1830395041f7bf03b43842f33): **deps**: Bump `mislav/bump-homebrew-formula-action` from v1.4 to v1.7 ([#318](@REPO/issues/318)) ([@dependabot](https://github.com/dependabot))
* [`5be6c15`](https://github.com/northwood-labs/taco-docs/commit/5be6c1596a11f9381650c552d0dee3e99bb3f52e): Create Dependabot config file ([#317](@REPO/issues/317)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`bde9025`](https://github.com/northwood-labs/taco-docs/commit/bde90257a3a57e69fa115d073f33ce50a7fa0457): Add Docker to dependabot file ([#319](@REPO/issues/319)) ([@khos2ow](https://github.com/khos2ow))
* [`4de4aa2`](https://github.com/northwood-labs/taco-docs/commit/4de4aa2ac19c79002f10954050ff5c3811e54814): Bump `golang` from 1.14.6 to 1.15.1 ([#320](@REPO/issues/320)) ([@dependabot](https://github.com/dependabot))
* [`d6e8352`](https://github.com/northwood-labs/taco-docs/commit/d6e8352c793b75326d3bd4bc972cf1b0725a18f6): Bump `github.com/zclconf/go-cty` from 1.6.0 to 1.6.1 ([#321](@REPO/issues/321)) ([@dependabot](https://github.com/dependabot))
* [`d8011f7`](https://github.com/northwood-labs/taco-docs/commit/d8011f730b8b29dc3a0ab0011d12fc552cbe65a7): Bump `golang` from 1.15.1-alpine to 1.15.2-alpine ([#323](@REPO/issues/323)) ([@dependabot](https://github.com/dependabot))

### :tractor: Refactor

* [`57a3584`](https://github.com/northwood-labs/taco-docs/commit/57a3584bed239bb21530c80da016e55bd40d7a58): Reorganize markdown format tests ([#244](@REPO/issues/244)) ([@khos2ow](https://github.com/khos2ow))
* [`c196c7c`](https://github.com/northwood-labs/taco-docs/commit/c196c7cc49ef90f112b4e55611fe633727520b7b): Add factory function to return format types ([#243](@REPO/issues/243)) ([@khos2ow](https://github.com/khos2ow))
* **[BC BREAK]** [`b6a6ad1`](https://github.com/northwood-labs/taco-docs/commit/b6a6ad1bbf77e99fa920dee784b9936170b35477): Remove deprecated flags ([#229](@REPO/issues/229)) ([@khos2ow](https://github.com/khos2ow))
* **[BC BREAK]** [`23c50e0`](https://github.com/northwood-labs/taco-docs/commit/23c50e0ad81431d1ffc2030730a9053f42b04d95): Deprecate multiple flags in favor of new ones ([#265](@REPO/issues/265)) ([@khos2ow](https://github.com/khos2ow))
* [`04a9ef4`](https://github.com/northwood-labs/taco-docs/commit/04a9ef49eb7c4bc1b81c3ebca0109a3f13a6b856): Refactor cli implemention and configuration ([#266](@REPO/issues/266)) ([@khos2ow](https://github.com/khos2ow))

### <!-- 0 -->:rocket: Features

* [`37375db`](https://github.com/northwood-labs/taco-docs/commit/37375db2832d4f32df7d737535c99ff61f8531bc): Add support for AsciiDoc renderer ([#241](@REPO/issues/241)) ([@rdelcampog](https://github.com/rdelcampog))
* [`b25909b`](https://github.com/northwood-labs/taco-docs/commit/b25909b5375b2dd347ebc45c68a0493bfa3704bc): Add new flag to sort inputs by type ([#246](@REPO/issues/246)) ([@khos2ow](https://github.com/khos2ow))
* [`b397b7d`](https://github.com/northwood-labs/taco-docs/commit/b397b7d46b41e22be1e9d2048ae6758a375f5dff): Add support for TOML renderer ([#197](@REPO/issues/197)) ([@khos2ow](https://github.com/khos2ow))
* [`e040439`](https://github.com/northwood-labs/taco-docs/commit/e0404399a7fd33ac44e74e6c352155fd515cba6a): Add new flags: --show, --show-all, --hide-all ([#267](@REPO/issues/267)) ([@khos2ow](https://github.com/khos2ow))
* [`38a86cb`](https://github.com/northwood-labs/taco-docs/commit/38a86cbdc5fcfc38eb9984111b884a31fc5e9e55): Build and push docker image ([#289](@REPO/issues/289)) ([@khos2ow](https://github.com/khos2ow))
* [`fd97ec5`](https://github.com/northwood-labs/taco-docs/commit/fd97ec5930df89ecb9658a2ff638b9666bf27b7a): Add support for `.terraform-docs.yml` config file ([#272](@REPO/issues/272)) ([@khos2ow](https://github.com/khos2ow))

### <!-- 1 -->:bug: Bug Fixes

* [`52653b5`](https://github.com/northwood-labs/taco-docs/commit/52653b5107c24bfafce9dcf29842e13a876acaab): Render special chars in variables' default value properly ([#284](@REPO/issues/284)) ()
* [`9729ec8`](https://github.com/northwood-labs/taco-docs/commit/9729ec8f929cc04c8eeee8313e475d0123e4e8bd): Normalize variables with CRLF line ending in heredoc ([#307](@REPO/issues/307)) ([@khos2ow](https://github.com/khos2ow))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`38316ec`](https://github.com/northwood-labs/taco-docs/commit/38316ec246f3eea4ed737f7d8bf54f2faabd6ff2): Custom order of changelog items ([#245](@REPO/issues/245)) ([@khos2ow](https://github.com/khos2ow))

## 0.9.1 — 2020-04-02

[Compare: v0.9.0 → v0.9.1](https://github.com/northwood-labs/taco-docs/compare/v0.9.0...v0.9.1)

### <!-- 1 -->:bug: Bug Fixes

* [`8f93043`](https://github.com/northwood-labs/taco-docs/commit/8f930437fa6945f8dd1cdd3bd83a95b1bb88e326): Make sure requirements section is sorted ([#233](@REPO/issues/233)) ([@khos2ow](https://github.com/khos2ow))
* [`80172d7`](https://github.com/northwood-labs/taco-docs/commit/80172d77f4f0a2a2c5d0cea09d390e60752d4062): Don't crash when reading header if 'main.tf' not found ([#235](@REPO/issues/235)) ([@khos2ow](https://github.com/khos2ow))

## 0.9.0 — 2020-03-31

[Compare: v0.8.2 → v0.9.0](https://github.com/northwood-labs/taco-docs/compare/v0.8.2...v0.9.0)

### :books: Documentation

* [`54ab7f9`](https://github.com/northwood-labs/taco-docs/commit/54ab7f9bbba020e44115a065c3760af4fad0c848): Auto generate formats document from examples ([#192](@REPO/issues/192)) ([@khos2ow](https://github.com/khos2ow))
* [`50254f0`](https://github.com/northwood-labs/taco-docs/commit/50254f0bece79b980f0ecd3d337153ed0e6a7488): Example git hook to keep module docs up to date ([#214](@REPO/issues/214)) ([@JamesUoM](https://github.com/JamesUoM))
* [`7c91a31`](https://github.com/northwood-labs/taco-docs/commit/7c91a310ced978d657efbd7ed6f1dc5c1a684ef2): Put reference to usage, cli, etc. in user guide ([#216](@REPO/issues/216)) ([@aaronsteers](https://github.com/aaronsteers))
* [`022f1dd`](https://github.com/northwood-labs/taco-docs/commit/022f1dd0a95b86390a960227a0f3f65d68c6c12e): Add installation guide for Windows users ([#218](@REPO/issues/218)) ([@aaronsteers](https://github.com/aaronsteers))
* [`40bd96b`](https://github.com/northwood-labs/taco-docs/commit/40bd96be44fc003e966dee330085538f86d87d8a): Enhance automatic document generation ([#227](@REPO/issues/227)) ([@khos2ow](https://github.com/khos2ow))

### :dependabot: Building and Dependencies

* [`3da3769`](https://github.com/northwood-labs/taco-docs/commit/3da3769a58a807ed826aaa02898cb84953c06e06): **deps**: Bump `github.com/stretchr/testify` from 1.4.0 to 1.5.0 ([#199](@REPO/issues/199)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`523e76e`](https://github.com/northwood-labs/taco-docs/commit/523e76eefe3f5625024b395119deb182001afa15): **deps**: Bump `github.com/stretchr/testify` from 1.5.0 to 1.5.1 ([#200](@REPO/issues/200)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`27b95c4`](https://github.com/northwood-labs/taco-docs/commit/27b95c4adf4d13ccc6f7619c07c705e4aecabad0): **deps**: Bump `github.com/spf13/cobra` from 0.0.5 to 0.0.6 ([#201](@REPO/issues/201)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`4e447c7`](https://github.com/northwood-labs/taco-docs/commit/4e447c752df5dcb08cd55b11508a4d6c679ed279): **deps**: Bump `github.com/imdario/mergo` from 0.3.8 to 0.3.9 ([#225](@REPO/issues/225)) ([@dependabot-preview](https://github.com/dependabot-preview))
* [`58119a3`](https://github.com/northwood-labs/taco-docs/commit/58119a3b46000ebc17e76bfc23c4b5e87ef4f9fc): **deps**: Bump `github.com/spf13/cobra` from 0.0.6 to 0.0.7 ([#228](@REPO/issues/228)) ([@dependabot-preview](https://github.com/dependabot-preview))

### :tractor: Refactor

* [`d4a0663`](https://github.com/northwood-labs/taco-docs/commit/d4a0663909511c951dae9a5ff336947a8faee0ba): Add tfconf.Options to load Module with ([#193](@REPO/issues/193)) ([@khos2ow](https://github.com/khos2ow))
* [`38e1897`](https://github.com/northwood-labs/taco-docs/commit/38e18970ed3f3ed774212553d77a3aeffb126e52): Introduce Format interface and expose to public pkg ([#195](@REPO/issues/195)) ([@khos2ow](https://github.com/khos2ow))
* [`743d4f3`](https://github.com/northwood-labs/taco-docs/commit/743d4f3ebc3cb2e610f69678d5a70c047d63e717): Add Default value types for better marshalling ([#196](@REPO/issues/196)) ([@khos2ow](https://github.com/khos2ow))

### <!-- 0 -->:rocket: Features

* [`42ad476`](https://github.com/northwood-labs/taco-docs/commit/42ad476a372983a2a3857401e5073640534deb84): Add support for YAML renderer ([#189](@REPO/issues/189)) ([@khos2ow](https://github.com/khos2ow))
* [`90fd2a3`](https://github.com/northwood-labs/taco-docs/commit/90fd2a322fac511fe667bf818a6a2b86335793bc): Render formatted results with go templates ([#177](@REPO/issues/177)) ([@khos2ow](https://github.com/khos2ow))
* [`31cdef0`](https://github.com/northwood-labs/taco-docs/commit/31cdef0f67699a7bb571ee5f79f6703c27bdc5be): Extract and render output values from Terraform ([#191](@REPO/issues/191)) ([@gshel](https://github.com/gshel))
* [`4ff4582`](https://github.com/northwood-labs/taco-docs/commit/4ff4582dff335be2112478a1d2de42dbf056f220): Show sensitivity of the output value in rendered result ([#207](@REPO/issues/207)) ([@khos2ow](https://github.com/khos2ow))
* [`b716a25`](https://github.com/northwood-labs/taco-docs/commit/b716a25811e12bb60be46ad6a64473f7472f8789): Add support for XML renderer ([#198](@REPO/issues/198)) ([@khos2ow](https://github.com/khos2ow))
* [`01c8fa1`](https://github.com/northwood-labs/taco-docs/commit/01c8fa1c61689a8ce1db500dd1030ef3e58e1bf7): Add support for fetching the module header from any file ([#217](@REPO/issues/217)) ([@khos2ow](https://github.com/khos2ow))
* [`b624175`](https://github.com/northwood-labs/taco-docs/commit/b6241751d92185c40b81bb5a6e1314df25dd0561): Add section for module requirements ([#222](@REPO/issues/222)) ([@khos2ow](https://github.com/khos2ow))
* [`2caf4af`](https://github.com/northwood-labs/taco-docs/commit/2caf4af15d7865139768695adddac9e4cfec61b1): Allow hiding the "Sensitive" column in markdown ([#223](@REPO/issues/223)) ([@julienduchesne](https://github.com/julienduchesne))
* [`9043f26`](https://github.com/northwood-labs/taco-docs/commit/9043f268ad25911bc55f695bf506bbb62a895b7d): Add support for tfvars hcl and json commands ([#226](@REPO/issues/226)) ([@khos2ow](https://github.com/khos2ow))

### <!-- 1 -->:bug: Bug Fixes

* [`f55fd6c`](https://github.com/northwood-labs/taco-docs/commit/f55fd6cc5bcf72b463e04fa5b3110a7918b15370): Fix type conversion for numbers ([#204](@REPO/issues/204)) ([@khos2ow](https://github.com/khos2ow))
* [`4365b49`](https://github.com/northwood-labs/taco-docs/commit/4365b4997b896407198e8178878ed7fc19beaf87): --no-header should not attempt reading main.tf file ([#224](@REPO/issues/224)) ()
* [`a3e0a56`](https://github.com/northwood-labs/taco-docs/commit/a3e0a56ce67c99eef6e730d7012db5e96ef5817c): Mark variables not required if default set to null ([#221](@REPO/issues/221)) ([@khos2ow](https://github.com/khos2ow))

### Enhance

* [`79e926e`](https://github.com/northwood-labs/taco-docs/commit/79e926ee43a651d5c787ffd8220cae1bf457c79f): Add extensive tests coverage for all the packages ([#208](@REPO/issues/208)) ([@khos2ow](https://github.com/khos2ow))

## 0.8.2 — 2020-02-03

[Compare: v0.8.1 → v0.8.2](https://github.com/northwood-labs/taco-docs/compare/v0.8.1...v0.8.2)

### <!-- 1 -->:bug: Bug Fixes

* [`bf38a75`](https://github.com/northwood-labs/taco-docs/commit/bf38a75fe36ae6dfcadf232daf208d2f9ddfc78a): Add newline between code block and trailing lines ([#184](@REPO/issues/184)) ([@khos2ow](https://github.com/khos2ow))
* [`43f69d3`](https://github.com/northwood-labs/taco-docs/commit/43f69d337f89b0c6c95d7c2fb3c6e346326de423): Preserve asterisk list in header and fix escaping ([#179](@REPO/issues/179)) ([@khos2ow](https://github.com/khos2ow))
* [`0f6b1e0`](https://github.com/northwood-labs/taco-docs/commit/0f6b1e0e82f62bfdd175305e1faaee5b82b0da3a): Add double space only at the end of paragraph lines ([#185](@REPO/issues/185)) ([@khos2ow](https://github.com/khos2ow))
* [`bf07706`](https://github.com/northwood-labs/taco-docs/commit/bf0770694749439b95e2e22852d54404439fce14): Do not escape markdown table inside module header ([#186](@REPO/issues/186)) ([@khos2ow](https://github.com/khos2ow))

## 0.8.1 — 2020-01-21

[Compare: v0.8.0 → v0.8.1](https://github.com/northwood-labs/taco-docs/compare/v0.8.0...v0.8.1)

### <!-- 1 -->:bug: Bug Fixes

* [`b604bce`](https://github.com/northwood-labs/taco-docs/commit/b604bce0c759fd162550b13e609b79aad5e57bdd): Show native map and list as default value in JSON ([#174](@REPO/issues/174)) ([@khos2ow](https://github.com/khos2ow))

## 0.8.0 — 2020-01-17

[Compare: v0.7.0 → v0.8.0](https://github.com/northwood-labs/taco-docs/compare/v0.7.0...v0.8.0)

### :books: Documentation

* **[BC BREAK]** [`d856c2c`](https://github.com/northwood-labs/taco-docs/commit/d856c2c11c4145289da376ac204a546b27752cbc): Update Module internal documentaion ([@khos2ow](https://github.com/khos2ow))
* **[BC BREAK]** [`1dd4553`](https://github.com/northwood-labs/taco-docs/commit/1dd45537d37f9d282a756308ea589e91fc6e47ea): Deprecate accepting files as commands param ([#163](@REPO/issues/163)) ([@khos2ow](https://github.com/khos2ow))
* [`0227bb9`](https://github.com/northwood-labs/taco-docs/commit/0227bb9d8db7bcaad71ecbd15900a5a2b6fe3773): Initial commit of usage documentation ([#162](@REPO/issues/162)) ([@khos2ow](https://github.com/khos2ow))

### :tractor: Refactor

* [`ea1f442`](https://github.com/northwood-labs/taco-docs/commit/ea1f442647d1ab5621660aa4744f39dca6c6223c): Move doc.Doc to tfconf.Module ([#136](@REPO/issues/136)) ([@khos2ow](https://github.com/khos2ow))

### <!-- 0 -->:rocket: Features

* [`b8da8c5`](https://github.com/northwood-labs/taco-docs/commit/b8da8c5020e3f372d14a1c2e65f8d0eda1a1fe47): Support Terraform 0.12.x configuration ([#113](@REPO/issues/113)) ([@moatra](https://github.com/moatra))
* [`2e6493c`](https://github.com/northwood-labs/taco-docs/commit/2e6493c8bbe968f8a181bbb371befdc3715e4330): Bump `golang` to latest v1.13 ([#133](@REPO/issues/133)) ([@chenrui333](https://github.com/chenrui333))
* **[BC BREAK]** [`ce32bf1`](https://github.com/northwood-labs/taco-docs/commit/ce32bf10a9699939ff7efd0aedd55d283252cf8f): Show 'providers' information ([#140](@REPO/issues/140)) ([@khos2ow](https://github.com/khos2ow))
* [`96565f8`](https://github.com/northwood-labs/taco-docs/commit/96565f8bccc291b401744e6686d4542b092652a0): Add '--no-color' flag to 'pretty' command ([#143](@REPO/issues/143)) ([@khos2ow](https://github.com/khos2ow))
* [`61554b9`](https://github.com/northwood-labs/taco-docs/commit/61554b9763717c12753edbed54fdb6d2f447d6f1): Add flags to not show different sections ([#144](@REPO/issues/144)) ([@khos2ow](https://github.com/khos2ow))
* [`453c7da`](https://github.com/northwood-labs/taco-docs/commit/453c7da2d4e178060f5f72ebab87e34d233048ef): Add '--no-escape' flag to 'json' command ([#147](@REPO/issues/147)) ([@khos2ow](https://github.com/khos2ow))

### <!-- 1 -->:bug: Bug Fixes

* [`80978e9`](https://github.com/northwood-labs/taco-docs/commit/80978e92f05492f2f012f8085e8ad76345ca9185): Reimplement '--no-sort' to be compatible with Terraform 0.12 configuration ([#141](@REPO/issues/141)) ([@khos2ow](https://github.com/khos2ow))
* [`a619384`](https://github.com/northwood-labs/taco-docs/commit/a6193849c44c5f2e1b06fabd39994797a9d7527f): Read leading comment lines if description is not provided ([#151](@REPO/issues/151)) ([@khos2ow](https://github.com/khos2ow))
* **[BC BREAK]** [`b3112d1`](https://github.com/northwood-labs/taco-docs/commit/b3112d135a6149c88a74669904226bac815ad63e): Read leading module header from main.tf ([#154](@REPO/issues/154)) ([@khos2ow](https://github.com/khos2ow))
* [`82e87aa`](https://github.com/northwood-labs/taco-docs/commit/82e87aa264924a8f0db35483a0b00eaf5153582f): Do not escape strings inside code blocks ([#155](@REPO/issues/155)) ([@khos2ow](https://github.com/khos2ow))
* [`61fc16c`](https://github.com/northwood-labs/taco-docs/commit/61fc16c101a38cc8b391d8164566b617e078b439): Show all JSON properties, empty or null ([#160](@REPO/issues/160)) ([@khos2ow](https://github.com/khos2ow))
* [`888bf9e`](https://github.com/northwood-labs/taco-docs/commit/888bf9e1d170746ae91ec9338e9e13681bfa1502): Do not wrap multiline blocks in table with <code> ([#164](@REPO/issues/164)) ([@khos2ow](https://github.com/khos2ow))
* [`1d33e9d`](https://github.com/northwood-labs/taco-docs/commit/1d33e9d03c53365beb6072f182bb6b2968ffa9da): Show empty JSON properties, as 'null' for all types ([#166](@REPO/issues/166)) ([@khos2ow](https://github.com/khos2ow))
* [`ccf2bcd`](https://github.com/northwood-labs/taco-docs/commit/ccf2bcd3da4b3f549a983233cef502023c29d6ad): Add double space at the end of multi-lines paragraph ([#169](@REPO/issues/169)) ([@khos2ow](https://github.com/khos2ow))
* [`6ff46a2`](https://github.com/northwood-labs/taco-docs/commit/6ff46a26b1c704b19a03dc2e7d3ac8fee977bd21): Do not escape any characters of a URL ([#170](@REPO/issues/170)) ([@khos2ow](https://github.com/khos2ow))

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* **[BC BREAK]** [`0e50fa9`](https://github.com/northwood-labs/taco-docs/commit/0e50fa933d973a785f58b397b45b8067b251e6e4): Rename flag to '--sort-by-required' ([#150](@REPO/issues/150)) ([@khos2ow](https://github.com/khos2ow))

### Enhance

* [`c11ab16`](https://github.com/northwood-labs/taco-docs/commit/c11ab1690250fe24d510c6f8d72560b563296017): Enable new go linters and fix the existing issues ([#132](@REPO/issues/132)) ([@khos2ow](https://github.com/khos2ow))
* [`fa8a756`](https://github.com/northwood-labs/taco-docs/commit/fa8a7569620453315bdec1efdf92b7b133e962a7): Bump `homebrew` formula version on release ([#135](@REPO/issues/135)) ([@khos2ow](https://github.com/khos2ow))
* **[BC BREAK]** [`ff80da2`](https://github.com/northwood-labs/taco-docs/commit/ff80da288f5337db56fc7988f109a064f60b55a1): Mark '--with-aggregate-type-defaults' as deprecated ([#148](@REPO/issues/148)) ([@khos2ow](https://github.com/khos2ow))

## 0.1.1 — 2017-08-16

[Compare: v0.1.0 → v0.1.1](https://github.com/northwood-labs/taco-docs/compare/v0.1.0...v0.1.1)

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`d59d826`](https://github.com/northwood-labs/taco-docs/commit/d59d8261fd3dd9490484dc76be46f5d6348c5229): Add --no-required option ([@s-urbaniak](https://github.com/s-urbaniak))

### Doc

* [`c2d16b1`](https://github.com/northwood-labs/taco-docs/commit/c2d16b1bc3d20e576252dda1d7673b08a822dce9): Snakecase -> camelcase ([@yields](https://github.com/yields))

## 0.1.0 — 2017-03-21

[Compare: v0.0.2 → v0.1.0](https://github.com/northwood-labs/taco-docs/compare/v0.0.2...v0.1.0)

### <!-- ZZZ -->:gear: Miscellaneous Tasks

* [`2b7095c`](https://github.com/northwood-labs/taco-docs/commit/2b7095c8d0caacde8b4c905a4832facc593cde1a): Add support for files ([@s-urbaniak](https://github.com/s-urbaniak))

### Doc

* [`8632ee6`](https://github.com/northwood-labs/taco-docs/commit/8632ee6d2146d66d720424fdc1ad2835f5858851): Allow top-level comments for variables when description missing ()
* [`ab50236`](https://github.com/northwood-labs/taco-docs/commit/ab50236699321ba56eaa41a61b81fa23b49fdafc): Placeholder for list types ()
* [`bda66e4`](https://github.com/northwood-labs/taco-docs/commit/bda66e47c2df4137e464fb560a5f75e6dcc7ab1c): Account for single whitespace after comment character in header ()

### Print/markdown

* [`1089829`](https://github.com/northwood-labs/taco-docs/commit/1089829ed8df30b973bcb802d8703919bfcd584a): Replace table cell newlines with HTML line breaks ()
* [`6fb15b1`](https://github.com/northwood-labs/taco-docs/commit/6fb15b1909a3b3208501ad5aea204415e01360be): Added line break conversion for outputs ()
* [`1ed94e5`](https://github.com/northwood-labs/taco-docs/commit/1ed94e57df77d98ef4942aae1f78f0f2fd594812): Better markdown description normalizations ()

## 0.0.2 — 2016-06-29

### Doc

* [`0df4d5a`](https://github.com/northwood-labs/taco-docs/commit/0df4d5ac09d6727e6135816cdee46a88c75c087e): Ignore comments with /** prefix ([@yields](https://github.com/yields))
* [`65e2192`](https://github.com/northwood-labs/taco-docs/commit/65e2192fa1b320491bb118282734534c57cd6b73): Fix map type ([@yields](https://github.com/yields))

### Outputs

* [`33b1130`](https://github.com/northwood-labs/taco-docs/commit/33b113047bd80969b6d53db8f214768a44e930a7): Use comments as description ([@yields](https://github.com/yields))

### Print

* [`ba2b483`](https://github.com/northwood-labs/taco-docs/commit/ba2b4831ca44a0c9780886cc8f066e7630f8e92f): Actually print head comment ([@yields](https://github.com/yields))
* [`7616408`](https://github.com/northwood-labs/taco-docs/commit/7616408eac153f1d0cd307a44fab2890ea9d0170): Wrap default values ([@yields](https://github.com/yields))

<p>Generated on 2026-07-24.</p>
