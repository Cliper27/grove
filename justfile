set windows-shell := ["powershell", "-Command"]

# Project info
appName := "grove"
version := "1.1.0"
distDir := "dist"
buildDir := "build"

# List all commands
default:
    @just --list

# ========
# CLEANING
# ========

# Delete `build` folder
[group("Clean")]
clean-build:
    @if (Test-Path {{buildDir}}) { Remove-Item -Recurse -Force {{buildDir}} }

# Delete `dist` folder
[group("Clean")]
clean-dist:
    @if (Test-Path {{distDir}}) { Remove-Item -Recurse -Force {{distDir}} }

# Delete `build` and `dist` folders
[group("Clean")]
clean: clean-build clean-dist
    @Write-Host "Clean complete."


# ========
# BUILDING
# ========

# Create folder if it doesn't exist
[group("Build")]
_mkdir name:
    @if (-not (Test-Path -Path {{name}})) { New-Item -ItemType Directory -Path {{name}} | Out-Null }

# Build Windows amd64
[group("Build")]
build-windows-amd64:
    @Write-Host "Building Windows AMD64..."
    @just _mkdir {{buildDir}}
    $env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" -o {{buildDir}}/{{appName}}-{{version}}-windows-amd64.exe ./cmd/grove

# Build Windows arm64
[group("Build")]
build-windows-arm64:
    @Write-Host "Building Windows ARM64..."
    @just _mkdir {{buildDir}}
    $env:GOOS="windows"; $env:GOARCH="arm64"; go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" -o {{buildDir}}/{{appName}}-{{version}}-windows-arm64.exe ./cmd/grove

# Build Linux amd64
[group("Build")]
build-linux-amd64:
    @Write-Host "Building Linux AMD64..."
    @just _mkdir {{buildDir}}
    $env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" -o {{buildDir}}/{{appName}}-{{version}}-linux-amd64 ./cmd/grove

# Build Linux arm64
[group("Build")]
build-linux-arm64:
    @Write-Host "Building Linux ARM64..."
    @just _mkdir {{buildDir}}
    $env:GOOS="linux"; $env:GOARCH="arm64"; go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" -o {{buildDir}}/{{appName}}-{{version}}-linux-arm64 ./cmd/grove

# Build macOS amd64
[group("Build")]
build-macos-amd64:
    @Write-Host "Building macOS AMD64..."
    @just _mkdir {{buildDir}}
    $env:GOOS="darwin"; $env:GOARCH="amd64"; go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" -o {{buildDir}}/{{appName}}-{{version}}-macos-amd64 ./cmd/grove

# Build macOS arm64
[group("Build")]
build-macos-arm64:
    @Write-Host "Building macOS ARM64..."
    @just _mkdir {{buildDir}}
    $env:GOOS="darwin"; $env:GOARCH="arm64"; go build -ldflags "-X github.com/Cliper27/grove/internal/version.Version={{version}}" -o {{buildDir}}/{{appName}}-{{version}}-macos-arm64 ./cmd/grove

# Build all platforms
[group("Build")]
build-all: clean-build build-windows-amd64 build-windows-arm64 build-linux-amd64 build-linux-arm64 build-macos-amd64 build-macos-arm64
    @Write-Host "All builds completed."


# =========
# PACKAGING
# =========

# Generate `.zip` for Windows amd64
[group("Package")]
package-windows-amd64:
    @Write-Host "Packaging Windows binaries..."
    @just _mkdir {{distDir}}
    Compress-Archive -Path {{buildDir}}/{{appName}}-{{version}}-windows-amd64.exe -DestinationPath {{distDir}}/{{appName}}-{{version}}-windows-amd64.zip

# Generate `.zip` for Windows arm64
[group("Package")]
package-windows-arm64:
    @Write-Host "Packaging Windows binaries..."
    @just _mkdir {{distDir}}
    Compress-Archive -Path {{buildDir}}/{{appName}}-{{version}}-windows-arm64.exe -DestinationPath {{distDir}}/{{appName}}-{{version}}-windows-arm64.zip

# Generate `.tar.gz` for Linux amd64
[group("Package")]
package-linux-amd64:
    @Write-Host "Packaging Linux binaries..."
    @just _mkdir {{distDir}}
    tar -czf {{buildDir}}/{{appName}}-{{version}}-linux-amd64.tar.gz -C {{distDir}} {{appName}}-{{version}}-linux-amd64

# Generate `.tar.gz` for Linux arm64
[group("Package")]
package-linux-arm64:
    @Write-Host "Packaging Linux binaries..."
    @just _mkdir {{distDir}}
    tar -czf {{buildDir}}/{{appName}}-{{version}}-linux-arm64.tar.gz -C {{distDir}} {{appName}}-{{version}}-linux-arm64

# Generate `.tar.gz` for macOS amd64
[group("Package")]
package-macos-amd64:
    @Write-Host "Packaging macOS binaries..."
    @just _mkdir {{distDir}}
    tar -czf {{buildDir}}/{{appName}}-{{version}}-macos-amd64.tar.gz -C {{distDir}} {{appName}}-{{version}}-macos-amd64

# Generate `.tar.gz` for macOS arm64
[group("Package")]
package-macos-arm64:
    @Write-Host "Packaging macOS binaries..."
    @just _mkdir {{distDir}}
    tar -czf {{buildDir}}/{{appName}}-{{version}}-macos-arm64.tar.gz -C {{distDir}} {{appName}}-{{version}}-macos-arm64

# Generate all compressed archives
[group("Package")]
package-all: clean-dist package-windows-amd64 package-windows-arm64 package-linux-amd64 package-linux-arm64 package-macos-amd64 package-macos-arm64
    @Write-Host "All packages created."


# =======
# TESTING
# =======

# go test
[group("Test")]
test:
    @Write-Host "Running Go tests..."
    @go test -v ./...

# Print app version
[group("Test")]
version:
    @go run ./cmd/grove version

# Print app help
[group("Test")]
help:
    @go run ./cmd/grove help
