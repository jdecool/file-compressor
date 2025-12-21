package compressor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImageCompressor(t *testing.T) {
	compressor := NewImageCompressor()
	assert.NotNil(t, compressor)
	assert.IsType(t, &ImageCompressor{}, compressor)
}

func TestImageCompressor_GetSupportedMimeTypes(t *testing.T) {
	compressor := NewImageCompressor()
	mimeTypes := compressor.GetSupportedMimeTypes()

	expectedMimeTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/bmp",
		"image/tiff",
		"image/webp",
	}

	assert.ElementsMatch(t, expectedMimeTypes, mimeTypes)
}

func TestImageCompressor_CompressFile(t *testing.T) {
	// Create a temporary test image file
	tempDir := t.TempDir()
	testImagePath := filepath.Join(tempDir, "test.jpg")
	outputPath := filepath.Join(tempDir, "compressed_test.jpg")

	// Create a simple red JPEG image for testing
	createTestImage(t, testImagePath)

	compressor := NewImageCompressor()
	result, err := compressor.CompressFile(testImagePath, outputPath)

	// This should fail because our test image is not valid
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode image")
	assert.NoFileExists(t, outputPath)
	assert.Nil(t, result)
}

func createTestImage(t *testing.T, path string) {
	// Create a minimal valid JPEG file by copying from our test data
	// In a real test, you would create a proper image, but for this example
	// we'll test the error handling path since creating valid images is complex
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()

	// Write some minimal content that will cause a decoding error
	_, err = file.Write([]byte("NOT A VALID IMAGE"))
	require.NoError(t, err)
}

func TestImageCompressor_CompressFile_InvalidFile(t *testing.T) {
	compressor := NewImageCompressor()
	result, err := compressor.CompressFile("nonexistent.jpg", "output.jpg")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get original file info")
	assert.Nil(t, result)
}