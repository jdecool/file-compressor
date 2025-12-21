package mime

import (
	"os"
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Error("NewDetector() should not return nil")
	}
}

func TestDetectMimeType(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "Text file",
			filePath: "testdata/test.txt",
			expected: "text/plain",
		},
		{
			name:     "JSON file",
			filePath: "testdata/test.json",
			expected: "application/json",
		},
		{
			name:     "HTML file",
			filePath: "testdata/test.html",
			expected: "text/html",
		},
		{
			name:     "CSV file",
			filePath: "testdata/test.csv",
			expected: "text/csv",
		},
		{
			name:     "XML file",
			filePath: "testdata/test.xml",
			expected: "text/xml",
		},
		{
			name:     "Fake PDF (actually text)",
			filePath: "testdata/fake.pdf",
			expected: "text/plain",
		},
		{
			name:     "Unknown extension",
			filePath: "testdata/unknown.ext",
			expected: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the full path to the test file
			fullPath := tt.filePath
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Skipf("Test file %s does not exist", fullPath)
			}

			mimeType := detector.DetectMimeType(fullPath)
			
			// Check if the detected MIME type contains the expected type
			// (we check for substring because mimetype library may add charset info)
			if !containsMIMEType(mimeType, tt.expected) {
				t.Errorf("DetectMimeType(%s) = %s, expected to contain %s", tt.filePath, mimeType, tt.expected)
			}
		})
	}
}

func TestDetectMimeTypeFromExtension(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "PDF extension",
			filePath: "test.pdf",
			expected: "application/pdf",
		},
		{
			name:     "Text extension",
			filePath: "test.txt",
			expected: "text/plain",
		},
		{
			name:     "JSON extension",
			filePath: "test.json",
			expected: "application/json",
		},
		{
			name:     "HTML extension",
			filePath: "test.html",
			expected: "text/html",
		},
		{
			name:     "CSV extension",
			filePath: "test.csv",
			expected: "text/csv",
		},
		{
			name:     "XML extension",
			filePath: "test.xml",
			expected: "application/xml",
		},
		{
			name:     "Unknown extension",
			filePath: "test.unknown",
			expected: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mimeType := detector.detectMimeTypeFromExtension(tt.filePath)
			if mimeType != tt.expected {
				t.Errorf("detectMimeTypeFromExtension(%s) = %s, want %s", tt.filePath, mimeType, tt.expected)
			}
		})
	}
}

func TestContentBasedDetectionPriority(t *testing.T) {
	detector := NewDetector()

	// Test that content-based detection takes priority over extension
	// fake.pdf contains text content but has .pdf extension
	mimeType := detector.DetectMimeType("testdata/fake.pdf")
	
	// Should detect as text/plain (content) not application/pdf (extension)
	if !containsMIMEType(mimeType, "text/plain") {
		t.Errorf("Content-based detection failed: got %s, expected text/plain for fake.pdf", mimeType)
	}
}

// Helper function to check if a MIME type string contains the expected type
// This handles cases where the library adds additional info like charset
func containsMIMEType(fullMIME, expectedType string) bool {
	// Handle cases where expectedType might have charset info
	if fullMIME == expectedType {
		return true
	}
	
	// Check if the expected type is the main part (before ; or space)
	mainPart := fullMIME
	if idx := indexOfSeparator(fullMIME); idx != -1 {
		mainPart = fullMIME[:idx]
	}
	
	return mainPart == expectedType
}

// Helper function to find the separator in MIME type strings
func indexOfSeparator(s string) int {
	for i, c := range s {
		if c == ';' || c == ' ' {
			return i
		}
	}
	return -1
}