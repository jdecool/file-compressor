package servicelocator

import (
	"testing"

	"github.com/jdecool/file-compressor/internal/compressor"
)

func TestServiceLocator(t *testing.T) {
	sl := NewServiceLocator()

	// Test empty service locator
	if sl.HasCompressor("application/pdf") {
		t.Error("ServiceLocator should not have compressors initially")
	}

	// Create a test compressor
	testCompressor := compressor.NewPdfCompressor()

	// Register the compressor
	sl.RegisterCompressor(testCompressor)

	// Test that compressors are registered
	if !sl.HasCompressor("application/pdf") {
		t.Error("ServiceLocator should have application/pdf compressor")
	}

	// Test getting compressors
	comp, exists := sl.GetCompressor("application/pdf")
	if !exists {
		t.Error("Should be able to get application/pdf compressor")
	}

	if comp == nil {
		t.Error("Retrieved compressor should not be nil")
	}

	// Test getting all compressors
	allCompressors := sl.GetAllCompressors()
	if len(allCompressors) != 1 {
		t.Errorf("Expected 1 compressor, got %d", len(allCompressors))
	}
}
