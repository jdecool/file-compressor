package compressor

import (
	"fmt"
	"image"
	"image/png"
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

	// Determine the output format and ensure correct file extension
	var imgFormat imaging.Format
	var encodeOptions []imaging.EncodeOption

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		imgFormat = imaging.JPEG
		encodeOptions = []imaging.EncodeOption{imaging.JPEGQuality(85)}
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") && !strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
		}
	case "png":
		imgFormat = imaging.PNG
		encodeOptions = []imaging.EncodeOption{imaging.PNGCompressionLevel(png.BestCompression)}
		if !strings.HasSuffix(strings.ToLower(outputPath), ".png") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".png"
		}
	case "gif":
		imgFormat = imaging.GIF
		encodeOptions = []imaging.EncodeOption{imaging.GIFNumColors(256)}
		if !strings.HasSuffix(strings.ToLower(outputPath), ".gif") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".gif"
		}
	case "bmp":
		imgFormat = imaging.BMP
		if !strings.HasSuffix(strings.ToLower(outputPath), ".bmp") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".bmp"
		}
	case "tiff":
		imgFormat = imaging.TIFF
		if !strings.HasSuffix(strings.ToLower(outputPath), ".tiff") && !strings.HasSuffix(strings.ToLower(outputPath), ".tif") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".tiff"
		}
	default:
		// For unknown formats, default to JPEG with quality compression
		imgFormat = imaging.JPEG
		encodeOptions = []imaging.EncodeOption{imaging.JPEGQuality(85)}
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") && !strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".jpg"
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = imaging.Encode(outFile, srcImage, imgFormat, encodeOptions...)
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
