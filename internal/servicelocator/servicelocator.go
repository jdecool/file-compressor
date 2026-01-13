package servicelocator

import (
	"github.com/jdecool/file-compressor/internal/compressor"
	"github.com/jdecool/file-compressor/internal/logger"
	"sync"
)

type ServiceLocator struct {
	mu          sync.RWMutex
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

	sl.mu.Lock()
	defer sl.mu.Unlock()

	for _, mimeType := range mimeTypes {
		sl.compressors[mimeType] = compressor
		if sl.compressors[mimeType] == nil {
			sl.logger.PrintfVerbose("Registered compressor for mime type: %s\n", mimeType)
		}
	}
}

func (sl *ServiceLocator) GetCompressor(mimeType string) (compressor.Compressor, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	compressor, exists := sl.compressors[mimeType]

	return compressor, exists
}

func (sl *ServiceLocator) GetAllCompressors() map[string]compressor.Compressor {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	copy := make(map[string]compressor.Compressor)
	for k, v := range sl.compressors {
		copy[k] = v
	}

	return copy
}

func (sl *ServiceLocator) HasCompressor(mimeType string) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	_, exists := sl.compressors[mimeType]

	return exists
}
