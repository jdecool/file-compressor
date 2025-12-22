package compressor

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
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
	createTestImage(t, testImagePath, "jpeg")

	compressor := NewImageCompressor()
	result, err := compressor.CompressFile(testImagePath, outputPath)

	// Should succeed now with valid image
	assert.NoError(t, err)
	assert.FileExists(t, outputPath)
	assert.NotNil(t, result)
	assert.Equal(t, testImagePath, result.OriginalFile)
	assert.Equal(t, outputPath, result.CompressedFile)
	assert.True(t, result.OriginalSize > 0)
	assert.True(t, result.CompressedSize > 0)
}

// createTestImage creates a simple 10x10 red image in the specified format
func createTestImage(t *testing.T, path string, format string) {
	// Create a simple 10x10 red image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
		}
	}

	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(file, img)
	case "gif":
		// For GIF, we need to use a different approach since standard lib doesn't support GIF encoding
		// We'll fall back to JPEG for this test
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "bmp":
		// For BMP, we'll fall back to JPEG for this test
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "tiff":
		// For TIFF, we'll fall back to JPEG for this test
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	default:
		// For other formats, we'll just use JPEG as a fallback for testing
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	}
	require.NoError(t, err)
}

func TestImageCompressor_CompressFile_InvalidFile(t *testing.T) {
	compressor := NewImageCompressor()
	result, err := compressor.CompressFile("nonexistent.jpg", "output.jpg")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get original file info")
	assert.Nil(t, result)
}

func TestImageCompressor_FormatPreservation(t *testing.T) {
	compressor := NewImageCompressor()
	tempDir := t.TempDir()

	testCases := []struct {
		name           string
		inputFormat    string
		inputExtension string
		expectedExt    string
	}{
		{"JPEG format preservation", "jpeg", ".jpg", ".jpg"},
		{"JPG format preservation", "jpeg", ".jpeg", ".jpg"},
		{"PNG format preservation", "png", ".png", ".png"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test image with specific extension
			inputPath := filepath.Join(tempDir, "test"+tc.inputExtension)
			outputPath := filepath.Join(tempDir, "compressed_test") // No extension - should be added

			// Create test image in the specified format
			createTestImage(t, inputPath, tc.inputFormat)

			// Compress the file
			result, err := compressor.CompressFile(inputPath, outputPath)

			// Should succeed
			assert.NoError(t, err, "Compression should succeed for %s", tc.name)
			assert.NotNil(t, result, "Result should not be nil for %s", tc.name)

			// Check that the output file has the correct extension
			actualOutputPath := result.CompressedFile
			actualExt := filepath.Ext(actualOutputPath)
			assert.Equal(t, tc.expectedExt, actualExt, "Output should have %s extension for %s", tc.expectedExt, tc.name)

			// Verify the file exists
			assert.FileExists(t, actualOutputPath, "Compressed file should exist for %s", tc.name)

			// Verify sizes are reasonable
			assert.True(t, result.OriginalSize > 0, "Original size should be > 0 for %s", tc.name)
			assert.True(t, result.CompressedSize > 0, "Compressed size should be > 0 for %s", tc.name)
		})
	}
}

func TestImageCompressor_FormatPreservation_WithOutputExtension(t *testing.T) {
	compressor := NewImageCompressor()
	tempDir := t.TempDir()

	// Test that when output path already has an extension, it's preserved
	inputPath := filepath.Join(tempDir, "test.png")
	outputPath := filepath.Join(tempDir, "compressed.png") // Already has .png extension

	// Create a PNG test image
	createTestImage(t, inputPath, "png")

	result, err := compressor.CompressFile(inputPath, outputPath)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, outputPath, result.CompressedFile, "Output path should be preserved when extension is already correct")
	assert.FileExists(t, outputPath)
}