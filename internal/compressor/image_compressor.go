package compressor

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageCompressor struct {
	supportedMimeTypes []string
}

func NewImageCompressor() *ImageCompressor {
	return &ImageCompressor{
		supportedMimeTypes: []string{
			"image/jpeg",
			"image/jpg",
			"image/png",
			"image/gif",
			"image/bmp",
			"image/tiff",
			"image/webp",
		},
	}
}

func (ic *ImageCompressor) CompressFile(filePath string, outputPath string) (*CompressionResult, error) {
	fmt.Printf("Image Compressor: Compressing file %s to %s\n", filepath.Base(filePath), filepath.Base(outputPath))

	// Get original file size
	originalFileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get original file info: %v", err)
	}

	srcFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer srcFile.Close()

	srcImage, format, err := image.Decode(srcFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := srcImage.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width > 2000 || height > 2000 {
		srcImage = imaging.Resize(srcImage, 2000, 0, imaging.Lanczos) // 0 means maintain aspect ratio
	}

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") && !strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
		}
	case "png", "gif", "bmp", "tiff", "webp":
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") && !strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
		}
	default:
		// For unknown formats, keep original extension but still compress
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	options := imaging.JPEGQuality(85)
	err = imaging.Encode(outFile, srcImage, imaging.JPEG, options)
	if err != nil {
		return nil, fmt.Errorf("failed to encode compressed image: %v", err)
	}

	// Get compressed file size
	compressedFileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get compressed file info: %v", err)
	}

	fmt.Printf("Image Compressor: Successfully compressed file to %s\n", outputPath)

	return &CompressionResult{
		OriginalFile:   filePath,
		CompressedFile: outputPath,
		OriginalSize:   originalFileInfo.Size(),
		CompressedSize: compressedFileInfo.Size(),
	}, nil
}

func (ic *ImageCompressor) GetSupportedMimeTypes() []string {
	return ic.supportedMimeTypes
}
