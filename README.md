# Optix File Processor

A powerful Go-based CLI tool for processing text, CSV, and JSON files with advanced features like text processing, batch operations, concurrency, and data transformation.

## 🚀 Features

### ✅ Completed Features

- **📄 File Reading & Display**: Read and display contents of text, CSV, and JSON files
- **📊 File Statistics**: Show detailed file information (size, lines, words)
- **🔍 Text Search**: Advanced pattern searching with regex support
- **🔄 Text Replace**: Search and replace operations with automatic backups
- **📋 Text Filtering**: Extract lines matching specific criteria
- **🔧 Text Transformations**: Case conversion and whitespace cleanup
- **✅ File Validation**: Built-in file existence and readability checks
- **🏗️ Strategy Pattern Architecture**: Extensible design for easy feature additions
- **🧪 Dry Run Mode**: Preview changes before applying them
- **💾 Automatic Backups**: Safe file modifications with backup creation

### 🔮 Planned Features

- **CSV Processing**: Data manipulation, filtering, and aggregation
- **JSON Processing**: Data extraction, validation, and transformation
- **Batch Processing**: Concurrent processing of multiple files
- **Configuration Management**: Settings and preferences
- **Report Generation**: Processing summaries and analytics

## 📦 Installation

### Prerequisites
- Go 1.23.4 or later

### Build from Source
```bash
git clone https://github.com/kcansari/optix.git
cd optix
make build
```

## 🎯 Usage

### Basic Commands

```bash
# Display help
./optix --help

# Show version information
./optix version

# Display file contents
./optix show myfile.txt

# Show file statistics
./optix stats data.csv
```

### 🔍 Text Search Operations

```bash
# Search for patterns in files
./optix search --pattern "error" --files "*.log"

# Use regular expressions
./optix search --pattern "user\d+" --regex --files "data.txt"

# Case-sensitive search with context
./optix search --pattern "TODO" --case-sensitive --context 2 --files "*.go"

# Whole word matching
./optix search --pattern "config" --whole-word --files "*.json"
```

### 🔄 Text Replace Operations

```bash
# Simple text replacement
./optix replace --find "old_url" --replace "new_url" --file config.txt

# Regex replacement with dry run
./optix replace --find "user\d+" --replace "customer$0" --regex --file data.txt --dry-run

# Replace with custom backup directory
./optix replace --find "localhost" --replace "production.com" --file config.txt --backup --backup-dir ./backups
```

### 📋 Text Filtering Operations

```bash
# Filter lines containing specific text
./optix filter --contains "WARNING" --input app.log --output warnings.log

# Use regex patterns
./optix filter --pattern "error\d+" --regex --input system.log


# Extract only matching parts
./optix filter --pattern "\b\w+@\w+\.\w+\b" --regex --only-matching --input emails.txt
```

### 🔧 Text Transformation Operations

```bash
# Convert text to uppercase
./optix transform --type upper --file document.txt

# Convert to lowercase with output to new file
./optix transform --type lower --file README.md --output readme.md

# Trim whitespace from all lines
./optix transform --type trim --file data.csv

# Preview transformation without changes
./optix transform --type title --file notes.txt --dry-run
```

## 🏗️ Architecture

Optix follows a **Strategy Pattern** design that makes it highly extensible and maintainable:

### Core Components

- **CLI Interface**: Command parsing and user interaction (Cobra framework)
- **File Reader Engine**: Multi-format file reading with strategy pattern
- **Text Processing Engine**: Search, replace, filter, and transform operations
- **File Validator**: File existence and readability validation
- **Configuration System**: Settings and preferences management

### Strategy Pattern Implementation

```
TextProcessor Interface
├── SearchProcessor     (Pattern matching with regex)
├── ReplaceProcessor   (Text replacement with backups)
├── FilterProcessor    (Line filtering and extraction)
└── TransformProcessor (Case conversion and cleanup)
```

Each processor is independent, testable, and can be easily extended without modifying existing code.

## 🧪 Development

### Build Commands

```bash
# Complete development workflow
make                    # Format, vet, build, and execute

# Individual commands
make build             # Build the binary
make fmt              # Format Go code
make vet              # Run Go vet analysis
make clean            # Clean build artifacts
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage and HTML report
make test-coverage

# Run tests with verbose output
make test-verbose
```

### Project Structure

```
optix/
├── cmd/optix/           # CLI commands
│   ├── root.go         # Root command setup
│   ├── version.go      # Version command
│   ├── show.go         # File display command
│   ├── stats.go        # File statistics command
│   ├── search.go       # Text search command
│   ├── replace.go      # Text replace command
│   ├── filter.go       # Text filter command
│   └── transform.go    # Text transform command
├── internal/
│   ├── reader/         # File reading strategies
│   ├── processor/      # Text processing strategies
│   ├── validator/      # File validation
│   ├── logger/         # Structured logging
│   └── version/        # Version information
├── test_data/          # Test files for debugging
├── main.go             # Application entry point
├── Makefile           # Build automation
```




## 📊 Checkpoint Progress

- ✅ **Checkpoint 1**: Project Foundation
- ✅ **Checkpoint 2**: Basic File Reading
- ✅ **Checkpoint 3**: Text Processing Engine
- 🔄 **Checkpoint 4**: CSV Processing (Planned)
- 🔄 **Checkpoint 5**: JSON Processing (Planned)
- 🔄 **Checkpoint 6**: Batch Processing (Planned)
- 🔄 **Checkpoint 7**: Advanced Features (Planned)
- 🔄 **Checkpoint 8**: Testing & Documentation (Planned)

## 🔧 Dependencies

- [Cobra](https://github.com/spf13/cobra) v1.8.1 - CLI framework
- Go standard library (regexp, strings, os, time, etc.)

## 🤝 Contributing

This is a learning project following a structured development plan. Each checkpoint builds upon the previous one, demonstrating Go best practices and design patterns.

### Development Principles

- **Strategy Pattern**: For extensible processing operations
- **Test-Driven Development**: Comprehensive test coverage
- **Clean Architecture**: Modular, maintainable code structure
- **User Experience**: Clear error messages and helpful CLI interface

## 📞 Contact

Created by [kcansari](https://github.com/kcansari) - feel free to contact me!

---

*This project demonstrates Go fundamentals through practical implementation, following modern software development practices and design patterns.*

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
