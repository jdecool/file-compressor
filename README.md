# File Compressor

A Go-based file compression utility that supports various file types and compression algorithms.

## Features

- PDF compression
- Multiple compression algorithms
- MIME type detection
- Service locator pattern for extensibility
- Easy-to-use CLI interface

## Installation

```bash
# Clone the repository
git clone https://github.com/jdecool/file-compressor.git
cd file-compressor

# Build the project
go build -o file-compressor
```

## Usage

```bash
# Basic usage
./file-compressor input.pdf output.pdf

# Help
./file-compressor --help
```

## Project Structure

- `compressor/` - Core compression logic
  - `compressor.go` - Main compression interface
  - `pdf_compressor.go` - PDF-specific compression
- `mime/` - MIME type detection
  - `detector.go` - MIME type detection logic
- `servicelocator/` - Service locator pattern implementation
  - `servicelocator.go` - Service registration and location
- `compression.go` - Compression algorithms
- `main.go` - CLI entry point

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test ./compressor
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

[MIT License](LICENSE)
