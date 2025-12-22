# File Compressor

A Go-based file compression utility that supports various file types and compression algorithms.

## Features

- PDF compression
- Image compression
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

- `internal/` - Internal package containing core functionality
  - `app/` - Application logic
    - `app.go` - Main application interface
    - `app_test.go` - Application tests
  - `compressor/` - Core compression logic
    - `compressor.go` - Main compression interface
    - `compressor_test.go` - Compression tests
    - `image_compressor.go` - Image-specific compression
    - `image_compressor_test.go` - Image compression tests
    - `pdf_compressor.go` - PDF-specific compression
    - `pdf_compressor_test.go` - PDF compression tests
  - `mime/` - MIME type detection
    - `detector.go` - MIME type detection logic
    - `detector_test.go` - MIME detection tests
    - `testdata/` - Test files for MIME detection
  - `servicelocator/` - Service locator pattern implementation
    - `servicelocator.go` - Service registration and location
    - `servicelocator_test.go` - Service locator tests
- `main.go` - CLI entry point
- `Makefile` - Build automation
- `go.mod` - Go module definition
- `go.sum` - Go dependency checksums

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test ./compressor
```

## License

[MIT License](LICENSE)
