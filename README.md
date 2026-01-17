# Grove

**Grove** is a cross-platform tool for defining, validating, and building directory trees based on user-defined schemas.  
You can use the prebuilt binaries to validate or create directory structures, or integrate Grove into your Go projects.

---

## Features

- Validate or build directory trees from `.gro` schema files.
- Prebuilt binaries available for Linux, macOS, and Windows.
- Define required, allowed, and denied files or folders.
- Support for including other schemas, with cycle and duplicate detection.
- Parse human-readable byte units (e.g., `10MB`, `1GB`).

---

## Installation

### CLI

Download the latest release from [GitHub Releases](https://github.com/cliper27/grove/releases) for your platform, or build from source:

```bash
go install github.com/cliper27/grove/cmd/grove@latest
```

## Usage

### Validate a directory
```bash
grove validate ./my-project ./schemas/go-project.gro
```

### Build a directory tree
```bash
grove build ./output ./schemas/go-project.gro
```

## Schema Example (`go-project.go`)
```yaml
include:
  - ./go-package.gro
  - ./go-command.gro
  - ./go-internal.gro

name: go-project
maxSize: 1GB

require:
  cmd/:
    schema: go-command
  internal/:
    schema: go-internal
  README.md:
    maxSize: 10MB
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

## Contributing
Bug reports, feature requests, and pull requests are welcome! Please open an issue or submit a PR.

## License
[MIT License](https://github.com/Cliper27/grove?tab=MIT-1-ov-file)