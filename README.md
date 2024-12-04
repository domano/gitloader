# GitHub Repository Content Downloader

Downloads GitHub repository contents into a single text file, skipping build artifacts and system files. Perfect for creating context files for AI assistants.

## Features

- Downloads all source files from a GitHub repository
- Skips common build artifacts, dependencies, and system files
- Preserves file paths in the output
- Uses go-git for authentication and private repos

## Installation

### Option 1: Direct Install

```bash
go install github.com/domano/gitloader@latest
```

### Option 2: From Source

```bash
git clone https://github.com/domano/gitloader
cd reponame
go mod init github-content
go get github.com/go-git/go-git/v5
go install
```

## Usage

```bash
go run main.go https://github.com/user/repo
```

Output will be saved as `user-repo.txt` in the current directory.

## Skipped Content

- Hidden files/directories (starting with .)
- Build directories (dist, build)
- Dependencies (node_modules, vendor)
- Lock files (package-lock.json, go.sum, etc.)

## Authentication

Uses your system's Git credentials for private repositories.