package compressor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PdfCompressor struct {
	supportedMimeTypes []string
}

func NewPdfCompressor() *PdfCompressor {
	return &PdfCompressor{
		supportedMimeTypes: []string{"application/pdf"},
	}
}

func (pc *PdfCompressor) CompressFile(filePath string, outputPath string) (*CompressionResult, error) {
	fmt.Printf("PDF Compressor: Compressing file %s to %s\n", filepath.Base(filePath), filepath.Base(outputPath))

	// Get original file size
	originalFileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get original file info: %v", err)
	}

	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.OPTIMIZE
	err = api.OptimizeFile(filePath, outputPath, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to optimize PDF file: %v", err)
	}

	// Get compressed file size
	compressedFileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get compressed file info: %v", err)
	}

	fmt.Printf("PDF Compressor: Successfully compressed file to %s\n", outputPath)
	
	return &CompressionResult{
		OriginalFile:   filePath,
		CompressedFile: outputPath,
		OriginalSize:   originalFileInfo.Size(),
		CompressedSize: compressedFileInfo.Size(),
	}, nil
}

func (pc *PdfCompressor) GetSupportedMimeTypes() []string {
	return pc.supportedMimeTypes
}
