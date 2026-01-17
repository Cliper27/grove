# Grove

**Grove** is a cross-platform tool and Go library for defining, validating, and building directory trees based on user-defined schemas.  
You can either use the **prebuilt binaries** to validate or create directory structures, or use the `parser` package in your Go code for programmatic access.

[![Go Reference](https://pkg.go.dev/badge/github.com/Cliper27/grove/internal/parser.svg)](https://pkg.go.dev/github.com/Cliper27/grove/internal/parser)

---

## Features

- Validate or build directory trees from `.gro` schema files.
- Prebuilt binaries available for Linux, macOS, and Windows.
- Go `parser` package allows integration in your own programs.
- Define required, allowed, and denied files or folders.
- Support for including other schemas, with cycle and duplicate detection.
- Parse human-readable byte units (e.g., `10MB`, `1GB`).

---

## Installation

### CLI

Download the latest release from [GitHub Releases](https://github.com/Cliper27/grove/releases) for your platform, or build from source:

```bash
go install github.com/cliper27/grove/cmd/grove@latest
```

### Go Library
For developers who want to use the parser package:
```bash
go get github.com/cliper27/grove/internal/parser
```

## Usage

### CLI
```bash
# Validate a directory against a schema
grove validate ./my-project ./schemas/go-project.gro

# Build a directory tree from a schema
grove build ./output ./schemas/go-project.gro
```

### Go Parser Package
```go
package main

import (
    "fmt"
    "github.com/cliper27/grove/internal/parser"
)

func main() {
    schema, err := parser.LoadSchema("path/to/go-project.gro")
    if err != nil {
        panic(err)
    }

    fmt.Println("Schema loaded:", schema.Name)

    bytes, err := parser.ParseByteUnits("10MB")
    if err != nil {
        panic(err)
    }
    fmt.Println("10MB =", bytes, "bytes")
}
```

## Schema Example (`go-project.gro`)
```yaml
name: go-project
options:
  maxSize: 1GB

include:
  - ./cmd.package/go-package.gro

require:
  cmd:
    schema: go-command
  internal:
    schema: go-internal
  README.md:
    maxSize: 10MB

allow:
  pkg:
    schema: go-package

deny:
  - node_modules/
  - "*.exe"
  - "^temp_[0-9]+.bin$"
```

## Contributing
Bug reports, feature requests, and pull requests are welcome! Please open an issue or submit a PR.

## License
MIT License
