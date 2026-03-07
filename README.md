# Grove

**Grove** is a cross-platform tool for defining, validating, and building directory trees based on user-defined schemas.  
You can use the prebuilt binaries to validate or create directory structures, or integrate Grove into your Go projects.


## Features

- 📋 **Structure rules** — define required, allowed, and denied files or folders
- 🔎 **Flexible matching** — support for both glob and regex patterns
- 🌳 **Directory validation** — validate entire directory trees from YAML schemas
- 🧩 **.gro schema extension** — optional extension for IDE syntax highlighting and tooling
- 💾 **Cross-platform binaries** — Linux, macOS, and Windows builds available
- 🔗 **Composable schemas** — include other schemas with cycle and duplicate detection


## Installation

### CLI

Download the latest release from [GitHub Releases](https://github.com/cliper27/grove/releases) for your platform, or build from source:

```bash
go install github.com/cliper27/grove/cmd/grove@latest
```


## Usage

### Validate a directory
```bash
grove check ./my-project ./schemas/go-project.gro
```

```bash
grove check . ./schemas/go-project.gro
```


## Schema Example (`go-project.go`)
```yaml
name: go-project
description: Standard Go project structure

include:
  - ./cmd.package/go-package.gro
  - ./cmd.package/go-command.gro
  - ./go-internal.gro

require:
  cmd/:
    schema: go-command
  internal/:
    schema: go-internal
  README.md:
    description: Project Documentation
  .gitignore:
  go.mod:
  go.sum:

allow:
  pkg/:
    schema: go-package

deny:
  - node_modules/
  - "*.exe"
  - "~^temp_[0-9]+.bin$"
```
