![Build](https://github.com/Cliper27/grove/actions/workflows/go-build.yml/badge.svg)
![Tests](https://github.com/Cliper27/grove/actions/workflows/go-test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Cliper27/grove)](https://goreportcard.com/report/github.com/Cliper27/grove)


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

### Get help
```bash
grove help
```

Output:
```bash
Validate project directory structure using schemas

Usage:
  grove [command]

Available Commands:
  check       Validate a directory against a schema
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print grove version

Flags:
  -h, --help      help for grove
  -v, --version   version for grove

Use "grove [command] --help" for more information about a command.
```

```bash
grove check --help
```

Output:
```bash
Validate a directory against a schema

Usage:
  grove check <dir> <schema> [flags]

Flags:
      --format string   Output format. Options are 'json' or 'tree'
  -h, --help            help for check
  -n, --no-color        Suppress cmd colors. Colors are automatically suppressed when not using a terminal that supports them.
  -o, --output string   Output to specified file
  -q, --quiet           Suppress stdout
```


### Validate a directory
```bash
grove check . go-project.gro
```

```bash
grove check ./my-project ./schemas/go-project.gro --format "tree"
```

```bash
grove check . ./schemas/go-project.gro --format "json" -q -o "result.json"
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
