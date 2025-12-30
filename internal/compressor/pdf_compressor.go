package compressor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jdecool/file-compressor/internal/logger"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PdfCompressor struct {
	supportedMimeTypes []string
	logger             *logger.Logger
}

func NewPdfCompressor() *PdfCompressor {
	return &PdfCompressor{
		supportedMimeTypes: []string{"application/pdf"},
		logger:             logger.NewLogger(false),
	}
}

func (pc *PdfCompressor) CompressFile(filePath string, outputPath string) (*CompressionResult, error) {
	pc.logger.PrintfVerbose("PDF Compressor: Compressing file %s to %s\n", filepath.Base(filePath), filepath.Base(outputPath))

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

	pc.logger.PrintfVerbose("PDF Compressor: Successfully compressed file to %s\n", outputPath)

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

func (pc *PdfCompressor) SetLogger(logger *logger.Logger) {
	pc.logger = logger
}
