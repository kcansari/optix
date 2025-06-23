# Optix File Processor

A powerful Go-based CLI tool for processing text, CSV, and JSON files with advanced features like batch processing, concurrency, and data transformation.

## Features

- **Multi-format Support**: Handle text, CSV, and JSON files
- **Batch Processing**: Process multiple files concurrently
- **File Validation**: Built-in file existence and readability checks
- **Extensible Architecture**: Strategy pattern for easy feature additions
- **High Performance**: Concurrent processing with goroutines
- **Professional CLI**: Built with Cobra framework

## Installation

### Prerequisites
- Go 1.23.4 or later

### Build from Source
```bash
git clone https://github.com/kcansari/optix.git
cd optix
make build
```

## Usage

### Basic Commands

```bash
# Display help
./optix --help

# Show version information
./optix version
```

### Dependencies

- [Cobra](https://github.com/spf13/cobra) v1.8.1 - CLI framework
- Go standard library


## Roadmap

### Current Status: Checkpoint 1
-  ✅  Go module initialization
-  ✅  Cobra CLI framework integration
-  ✅  Version command implementation
-  ✅  File validation with strategy pattern
-  ✅  Clean project structure
-  ✅  Basic testing setup

### Next: Checkpoint 2 🚧
- [ ] File reading functionality (txt, csv, json)
- [ ] Content display commands
- [ ] File statistics (line count, word count, size)
- [ ] Enhanced error handling

## 🤝 Contributing

This is a learning project, but feel free to suggest improvements or report issues.

## 📄 License

This project is licensed under the MIT License.

## 📞 Contact

Created by [kcansari](https://github.com/kcansari) - feel free to contact me!

---

*This project is part of a Go learning journey. Each commit represents a step in understanding Go fundamentals through practical implementation.*

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.