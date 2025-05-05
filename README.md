# Ingest Clone

A simplified Go implementation inspired by [gitingest](https://github.com/cyclotruc/gitingest). This tool allows you to convert specified files or folders to a standard format file for easier analysis.

## Features

- **File/Directory Analysis**: Scan files and directories, analyze their contents
- **Format Conversion**: Convert files to a consistent text format for easier analysis
- **Pattern Matching**: Include or exclude files based on patterns
- **Simple CLI**: Easy-to-use command-line interface

## Installation

### Prerequisites

- Go 1.22 or higher

### Build from Source

```bash
git clone https://github.com/yourusername/ingest-clone.git
cd ingest-clone
go build -o ingest ./cmd/ingest
```

## Usage

```bash
# Basic usage - scan current directory
./ingest

# Analyze a specific directory
./ingest /path/to/directory

# Specify output file
./ingest -o output.txt /path/to/directory

# Include specific patterns
./ingest -i "*.go,*.md" /path/to/directory

# Exclude specific patterns
./ingest -e "*.tmp,vendor/" /path/to/directory

# Combine include and exclude
./ingest -i "*.go" -e "vendor/" /path/to/directory
```

## Options

- `-o, --output`: Output file (default: digest.txt)
- `-i, --include`: Patterns to include (comma-separated)
- `-e, --exclude`: Patterns to exclude (comma-separated)
- `-h, --help`: Show help
- `-v, --version`: Show version information

## Output Format

The output includes:

1. **Summary**: Information about the analyzed directory or files
2. **Directory Structure**: A tree-like representation of the file structure
3. **File Contents**: Contents of analyzed files with appropriate headers

Example:

```
Directory: myproject

Files analyzed: 15

Estimated tokens: 4.5k

Directory structure:
└── myproject/
    ├── cmd/
    │   └── main.go
    ├── pkg/
    │   ├── analyzer/
    │   │   └── analyzer.go
    │   └── formatter/
    │       └── formatter.go
    └── README.md

================================================
FILE: cmd/main.go
================================================
package main

func main() {
    // Main implementation
}

================================================
FILE: pkg/analyzer/analyzer.go
================================================
package analyzer

// Implementation details
...
```

## License

MIT

## Acknowledgements

This project is inspired by [gitingest](https://github.com/cyclotruc/gitingest) by cyclotruc. 