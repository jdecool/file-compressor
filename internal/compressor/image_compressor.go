package compressor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/jdecool/file-compressor/internal/logger"
	"github.com/rwcarlsen/goexif/exif"
)

type ImageCompressor struct {
	supportedMimeTypes []string
	logger             *logger.Logger
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
		logger: logger.NewLogger(false),
	}
}

func (ic *ImageCompressor) CompressFile(filePath string, outputPath string) (*CompressionResult, error) {
	ic.logger.PrintfVerbose("Image Compressor: Compressing file %s to %s\n", filepath.Base(filePath), filepath.Base(outputPath))

	// Get original file size
	originalFileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get original file info: %v", err)
	}

	// Extract EXIF data before processing (if available)
	var exifData *exif.Exif
	srcFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}

	exifData, err = exif.Decode(srcFile)
	if err != nil {
		ic.logger.PrintfVerbose("Image Compressor: No EXIF data found or unable to decode EXIF (this is normal for some images)\n")
		exifData = nil // No EXIF data available, continue without it
	}
	srcFile.Close()

	// Reopen file for image decoding
	srcFile, err = os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to reopen image file: %v", err)
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

	// For JPEG files with EXIF data, use special handling to preserve metadata
	if (strings.ToLower(format) == "jpeg" || strings.ToLower(format) == "jpg") && exifData != nil {
		err = ic.compressJPEGWithEXIF(srcImage, outputPath, exifData, filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to compress JPEG with EXIF: %v", err)
		}
	} else {
		// Standard compression without EXIF preservation
		outFile, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %v", err)
		}
		defer outFile.Close()

		err = imaging.Encode(outFile, srcImage, imgFormat, encodeOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to encode compressed image: %v", err)
		}
	}

	// Get compressed file size
	compressedFileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get compressed file info: %v", err)
	}

	ic.logger.PrintfVerbose("Image Compressor: Successfully compressed file to %s\n", outputPath)

	return &CompressionResult{
		OriginalFile:   filePath,
		CompressedFile: outputPath,
		OriginalSize:   originalFileInfo.Size(),
		CompressedSize: compressedFileInfo.Size(),
	}, nil
}

// compressJPEGWithEXIF compresses a JPEG image while preserving EXIF metadata
func (ic *ImageCompressor) compressJPEGWithEXIF(img image.Image, outputPath string, exifData *exif.Exif, originalPath string) error {
	// Encode the processed image to a buffer
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return fmt.Errorf("failed to encode JPEG: %v", err)
	}

	// Read the original file to get the full EXIF segment
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return fmt.Errorf("failed to open original file: %v", err)
	}
	defer originalFile.Close()

	originalData, err := io.ReadAll(originalFile)
	if err != nil {
		return fmt.Errorf("failed to read original file: %v", err)
	}

	// Parse the original JPEG to extract EXIF segment
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(originalData)
	if err != nil {
		// If we can't parse with the JPEG structure library, fall back to basic encoding
		ic.logger.PrintfVerbose("Image Compressor: Could not parse JPEG structure, saving without EXIF preservation\n")
		return os.WriteFile(outputPath, buf.Bytes(), 0644)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Find the EXIF segment from the original
	var exifSegmentData []byte
	segments := sl.Segments()
	for _, segment := range segments {
		if segment.MarkerName == "APP1" {
			// Check if this is an EXIF segment (starts with "Exif\x00\x00")
			if len(segment.Data) >= 6 {
				if string(segment.Data[0:4]) == "Exif" && segment.Data[4] == 0 && segment.Data[5] == 0 {
					exifSegmentData = segment.Data
					break
				}
			}
		}
	}

	if exifSegmentData == nil {
		// No EXIF segment found in original
		ic.logger.PrintfVerbose("Image Compressor: No EXIF APP1 segment found in original\n")
		return os.WriteFile(outputPath, buf.Bytes(), 0644)
	}

	// Parse the newly compressed JPEG
	intfc2, err := jmp.ParseBytes(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to parse compressed JPEG: %v", err)
	}

	newSl := intfc2.(*jpegstructure.SegmentList)

	// Insert the EXIF segment into the compressed image
	// Find the SOI (Start of Image) marker and insert after it
	segments2 := newSl.Segments()
	if len(segments2) < 2 {
		ic.logger.PrintfVerbose("Image Compressor: Compressed JPEG has too few segments\n")
		return os.WriteFile(outputPath, buf.Bytes(), 0644)
	}

	// Create a new segment list with the EXIF segment inserted after SOI
	newSegments := make([]*jpegstructure.Segment, 0, len(segments2)+1)
	newSegments = append(newSegments, segments2[0]) // SOI

	// Create new EXIF segment
	exifSeg := &jpegstructure.Segment{
		MarkerId:   0xe1, // APP1
		MarkerName: "APP1",
		Offset:     0,
		Data:       exifSegmentData,
	}
	newSegments = append(newSegments, exifSeg)

	// Add remaining segments (skip any existing EXIF segments)
	for i := 1; i < len(segments2); i++ {
		seg := segments2[i]
		// Skip existing APP1/EXIF segments to avoid duplicates
		if seg.MarkerName == "APP1" {
			if len(seg.Data) >= 6 && string(seg.Data[0:4]) == "Exif" {
				continue // Skip existing EXIF
			}
		}
		newSegments = append(newSegments, seg)
	}

	// Build new segment list
	newSl2 := jpegstructure.NewSegmentList(newSegments)

	// Write the final JPEG with EXIF to file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	err = newSl2.Write(outputFile)
	if err != nil {
		return fmt.Errorf("failed to write JPEG with EXIF: %v", err)
	}

	ic.logger.PrintfVerbose("Image Compressor: EXIF metadata preserved\n")

	return nil
}

func (ic *ImageCompressor) GetSupportedMimeTypes() []string {
	return ic.supportedMimeTypes
}

func (ic *ImageCompressor) SetLogger(logger *logger.Logger) {
	ic.logger = logger
}
