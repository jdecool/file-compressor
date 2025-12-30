package servicelocator

import (
	"github.com/jdecool/file-compressor/internal/compressor"
	"github.com/jdecool/file-compressor/internal/logger"
)

type ServiceLocator struct {
	compressors map[string]compressor.Compressor
	logger      *logger.Logger
}

func NewServiceLocator() *ServiceLocator {
	return &ServiceLocator{
		compressors: make(map[string]compressor.Compressor),
		logger:      logger.NewLogger(false),
	}
}

func (sl *ServiceLocator) RegisterCompressor(compressor compressor.Compressor) {
	mimeTypes := compressor.GetSupportedMimeTypes()
	for _, mimeType := range mimeTypes {
		sl.compressors[mimeType] = compressor
		if sl.compressors[mimeType] == nil {
			sl.logger.PrintfVerbose("Registered compressor for mime type: %s\n", mimeType)
		}
	}
}

func (sl *ServiceLocator) GetCompressor(mimeType string) (compressor.Compressor, bool) {
	compressor, exists := sl.compressors[mimeType]

	return compressor, exists
}

func (sl *ServiceLocator) GetAllCompressors() map[string]compressor.Compressor {
	copy := make(map[string]compressor.Compressor)
	for k, v := range sl.compressors {
		copy[k] = v
	}

	return copy
}

func (sl *ServiceLocator) HasCompressor(mimeType string) bool {
	_, exists := sl.compressors[mimeType]

	return exists
}
