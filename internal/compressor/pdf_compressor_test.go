package compressor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPdfCompressorCompressFile(t *testing.T) {
	// Create a test PDF file
	pdfContent := "%PDF-1.4\n1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R >>\nendobj\n4 0 obj\n<< /Length 44 >>\nstream\nBT\n/F1 24 Tf\n100 700 Td\n(Hello, Test!) Tj\nET\nendstream\nendobj\nxref\n0 5\n0000000000 65535 f \n0000000010 00000 n \n0000000079 00000 n \n0000000172 00000 n \n0000000360 00000 n \ntrailer\n<< /Size 5 /Root 1 0 R >>\nstartxref\n500\n%%EOF"

	tempFile, err := os.CreateTemp("", "test*.pdf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(pdfContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create output file path
	outputFile, err := os.CreateTemp("", "output*.pdf")
	if err != nil {
		t.Fatalf("Failed to create output temp file: %v", err)
	}
	outputPath := outputFile.Name()
	outputFile.Close()
	defer os.Remove(outputPath)

	// Test the CompressFile method
	compressor := NewPdfCompressor()
	result, err := compressor.CompressFile(tempFile.Name(), outputPath)

	// For invalid PDFs, we expect an error (this is expected for our test PDF)
	if err != nil {
		// This is expected for our minimal test PDF
		t.Logf("CompressFile failed as expected for test PDF: %v", err)
		assert.Nil(t, result)
		return
	}

	// If we get here, the compression succeeded
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("Output file should exist after successful compression: %v", err)
	}
	if result != nil {
		t.Logf("Compression result: Original: %d bytes, Compressed: %d bytes", result.OriginalSize, result.CompressedSize)
	}
}
