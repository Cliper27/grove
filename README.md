![Build](https://github.com/Cliper27/grove/actions/workflows/go-build.yml/badge.svg)
![Tests](https://github.com/Cliper27/grove/actions/workflows/go-test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Cliper27/grove)](https://goreportcard.com/report/github.com/Cliper27/grove)


# Grove

[**Grove**](https://github.com/Cliper27/grove) is a cross-platform tool for defining, validating, and building directory trees based on user-defined schemas.  
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


## Schema syntax
Schemas are written in YAML. The file extension does not matter, but `.gro` is recommended for CLI autocompletion and the future VSCode extension. One file can only define one schema.

### Fields

#### `name` _(String, required)_
The name of the schema. Simple and short.
```yaml
name: go-app
```

#### `description` _(String, optional)_
Brief description of the schema.
```yaml
description: Standard go app layout
```

#### `require` _(Mapping, optional)_
The list of patterns (see the next section) that **MUST** exist for the directory to be valid against this schema.
```yaml
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
```

#### `allow` _(Mapping, optional)_
Additional patterns that are allowed but not required. By default, any file or directory is allowed unless it is explicitly denied. For example: _"Allow the folder `cmd`, but it must follow the schema `cmd-folder`"_. It can also define descriptions, either for folders or files.
```yaml
allow:
  pkg/:
    schema: go-package
  justfile:
    description: Justfile with build recipes
```

#### `deny` _(Sequence, optional)_
Patterns that are explicitly forbidden in the directory.
```yaml
deny:
  - node_modules/
  - "*.exe"
  - "~^temp_[0-9]+.bin$"
```

#### `include` _(Sequence, optional)_
If another schema is referenced in any other field, it must be specified here so that `grove` can find it. Each include is a path relative to the schema file.
```yaml
include:
  - ./cmd.package/go-package.gro
  - ./cmd.package/go-command.gro
  - ./go-internal.gro
```


### Patterns
A pattern represents a file or folder name. It can match a single item or a group of items depending on the pattern type.

Patterns are used as keys inside `require` and `allow`, and as values inside `deny`.

A pattern can be:

#### Exact
Matches a specific file or directory name exactly.
```yaml
README.md:
src/:
```

#### Glob
Matches multiple files using glob syntax.
```yaml
"*.go":
"*.test.js":
```

#### Regex
Matches names using a regular expression. Regex patterns are prefixed with `~`.
```yaml
"~^test_[0-9]+\\.txt$":
```

### Pattern properties
Patterns defined under require or allow can have additional properties.
- `schema` _(String)_: Specifies that the matched directory must itself conform to another schema. This schema must be included in the `include` field.
```yaml
cmd/:
  schema: go-command
```
This means every matching directory must validate against the go-command schema.

- `description` _(String)_: Same as the description of a schema. If both `description` and `schema` are specified, `description` overrides the referenced schema's description.


### Schema Example (`go-project.gro`)
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
  justfile:
    description: Justfile with build recipes

deny:
  - node_modules/
  - "*.exe"
  - "~^temp_[0-9]+.bin$"  # ~ indicates regex
```
