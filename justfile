set shell := ["bash", "-c"]
set windows-shell := ["powershell", "-Command"]

# Project info
appName := "grove"
version := "1.1.0"
distDir := "dist"
buildDir := "build"

# List all commands
default:
    @just --list

# =======
# HELPERS
# =======
# Create folder if it doesn't exist: Windows
[group("Helpers")]
_mkdir-windows name:
    @if (-not (Test-Path -Path {{name}})) { New-Item -ItemType Directory -Path {{name}} | Out-Null }

# Create folder if it doesn't exist: macOS
[group("Helpers")]
_mkdir-macos name:
    @mkdir -p "{{name}}"

# Create folder if it doesn't exist: Linux
[group("Helpers")]
_mkdir-linux name:
    @mkdir -p "{{name}}"


# Write-Host
[group("Helpers")]
_print-windows string:
    @Write-Host "{{string}}"

# echo
[group("Helpers")]
_print-macos string:
    @echo "{{string}}"

# echo
[group("Helpers")]
_print-linux string:
    @echo "{{string}}"


# Delete folder Windows
[group("Helpers")]
_clean-windows dir:
    @if (Test-Path {{dir}}) { Remove-Item -Recurse -Force {{dir}} }

# Delete folder macOS
[group("Helpers")]
_clean-macos dir:
    @[ -d "{{dir}}" ] && rm -rf "{{dir}}"

# Delete folder Linux
[group("Helpers")]
_clean-linux dir:
    @[ -d "{{dir}}" ] && rm -rf "{{dir}}"


# ========
# CLEANING
# ========

# Delete folder
[group("Clean")]
clean dir:
    @just _clean-{{os()}} "{{dir}}"

# Delete `build` and `dist` folders
[group("Clean")]
clean-all:
    @just clean {{buildDir}}
    @just clean {{distDir}}
    @just _print-{{os()}} "Clean complete."


# ========
# BUILDING
# ========

# Build any os and arch from unix
[group("Build")]
_build-unix os arch:
    GOOS={{os}} GOARCH={{arch}} \
    go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" \
    -o {{buildDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}{{ if os == "windows" { ".exe" } else { "" } }} ./cmd/grove

# Build any os and arch from windows
[group("Build")]
_build-windows os arch:
    $env:GOOS="{{os}}"; $env:GOARCH="{{arch}}"; \
    go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" \
    -o {{buildDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}{{ if os == "windows" { ".exe" } else { "" } }} ./cmd/grove

# Build any OS and architecture
[group("Build")]
build os arch:
    @just _print-{{os()}} "Building {{os}} {{arch}}..."
    @just _mkdir-{{os()}} {{buildDir}}
    @just _build-{{ if os() == "windows" { "windows" } else { "unix" } }} {{os}} {{arch}}

# Build any OS and architecture
[group("Build")]
build-all:
    @just clean {{buildDir}}
    @just build "windows" "amd64"
    @just build "windows" "arm64"
    @just build "linux" "amd64"
    @just build "linux" "arm64"
    @just build "darwin" "amd64"
    @just build "darwin" "arm64"
    @just _print-{{os()}} "Done!"


# =========
# PACKAGING
# =========

# Generate Windows `.zip` from Windows
[group("Package")]
_zip_windows_from_windows os arch:
    @if (Test-Path {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.zip) { Remove-Item {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.zip }
    @Compress-Archive -Path {{buildDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.exe -DestinationPath {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.zip

# Generate Unix `.tar.gz` from Windows
[group("Package")]
_zip_unix_from_windows os arch:
    @tar -czf {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.tar.gz -C {{buildDir}} {{appName}}-{{version}}-{{os}}-{{arch}}

# Generate Windows `.zip` from Unix
[group("Package")]
_zip_windows_from_unix os arch:
    @just _mkdir-{{os()}} {{distDir}}
    @if [ -f {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.zip ]; then rm {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.zip; fi
    @zip -j {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.zip {{buildDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.exe

# Generate Unix `.tar.gz` from Unix
[group("Package")]
_zip_unix_from_unix os arch:
    @just _mkdir-{{os()}} {{distDir}}
    @if [ -f {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.tar.gz ]; then rm {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.tar.gz; fi
    @tar -czf {{distDir}}/{{appName}}-{{version}}-{{os}}-{{arch}}.tar.gz -C {{buildDir}} {{appName}}-{{version}}-{{os}}-{{arch}}

# Generate any compressed archive
[group("Package")]
package os_family os arch:
    @just _print-{{os()}} "Packaging {{os}} binaries..."
    @just _mkdir-{{os()}} {{distDir}}
    @just _zip_{{os_family}}_from_{{os_family()}} {{os}} {{arch}}

# Generate all compressed archives
[group("Package")]
package-all:
    @just clean {{distDir}}
    @just package "windows" "windows" "amd64"
    @just package "windows" "windows" "arm64"
    @just package "unix" "linux" "amd64"
    @just package "unix" "linux" "arm64"
    @just package "unix" "darwin" "amd64"
    @just package "unix" "darwin" "arm64"


# =======
# TESTING
# =======

# Print app help
[group("Test")]
run-check *args:
    @go run ./cmd/grove check {{args}}


# go test
[group("Test")]
test:
    @just _print-{{os()}} "Running Go tests..."
    @go test -v ./...

# Print app version
[group("Test")]
version:
    @go run ./cmd/grove version

# Print app help
[group("Test")]
help:
    @go run ./cmd/grove help
