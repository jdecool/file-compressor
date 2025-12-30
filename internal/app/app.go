package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/jdecool/file-compressor/internal/compressor"
	"github.com/jdecool/file-compressor/internal/logger"
	"github.com/jdecool/file-compressor/internal/mime"
	"github.com/jdecool/file-compressor/internal/servicelocator"
)

type Application struct {
	logger          *logger.Logger
	maxWorkers      int
	replaceOriginal bool
	serviceLocator  *servicelocator.ServiceLocator
	mimeDetector    *mime.Detector
	compressionResults []*compressor.CompressionResult
}

func NewApplication() *Application {
	return &Application{
		logger:          logger.NewLogger(false),
		maxWorkers:      GetDefaultWorkersCount(),
		replaceOriginal: false,
		serviceLocator:  servicelocator.NewServiceLocator(),
		mimeDetector:    mime.NewDetector(),
	}
}

func (a *Application) RegisterCompressor(c compressor.Compressor) {
	// Set the logger for the compressor if it supports it
	if setter, ok := c.(interface{ SetLogger(*logger.Logger) }); ok {
		setter.SetLogger(a.logger)
	}
	a.serviceLocator.RegisterCompressor(c)
}

func (a *Application) SetVerboseMode(isVerbose bool) {
	a.logger.SetVerbose(isVerbose)
}

func (a *Application) SetMaxWorkers(count int) error {
	if count < 1 {
		a.maxWorkers = 1

		return fmt.Errorf("worker count must be at least 1, setting to 1")
	}

	if count > runtime.NumCPU() {
		a.maxWorkers = runtime.NumCPU()

		return fmt.Errorf("worker count cannot exceed number of CPU cores (%d), setting to %d", runtime.NumCPU(), runtime.NumCPU())
	}

	a.maxWorkers = count

	return nil
}

func (a *Application) SetReplaceOriginal(replace bool) {
	a.replaceOriginal = replace
}

func (a *Application) replaceOriginalFile(originalPath, compressedPath string) error {
	if err := os.Remove(originalPath); err != nil {
		return fmt.Errorf("failed to remove original file: %v", err)
	}

	if err := os.Rename(compressedPath, originalPath); err != nil {
		return fmt.Errorf("failed to rename compressed file: %v", err)
	}

	return nil
}

func (a *Application) Run(inputPaths []string) {
	a.logger.PrintlnVerbose("Starting file compression process...")
	a.logger.PrintfVerbose("Input paths: %v\n", inputPaths)
	a.logger.PrintfVerbose("Using %d workers for parallel processing\n", a.maxWorkers)

	fileChan := make(chan string, 100)
	doneChan := make(chan bool)
	var wg sync.WaitGroup

	for i := 0; i < a.maxWorkers; i++ {
		wg.Add(1)
		go a.worker(fileChan, &wg)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	go func() {
		for _, path := range inputPaths {
			a.logger.PrintfVerbose("Processing path: %s\n", path)

			fileInfo, err := os.Stat(path)
			if err != nil {
				a.logger.PrintfError("Error accessing path %s: %v\n", path, err)
				continue
			}

			if fileInfo.IsDir() {
				err := a.browseDirectoryAndSendFiles(path, fileChan)
				if err != nil {
					a.logger.PrintfError("Error browsing directory %s: %v\n", path, err)
				}
			} else {
				fileChan <- path
			}
		}
		close(fileChan)
	}()

	<-doneChan

	a.printOperationSummary()
	a.logger.PrintlnVerbose("File compression completed.")
}

func (a *Application) browseDirectoryAndSendFiles(rootPath string, fileChan chan<- string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileChan <- path
		}

		return nil
	})
}

func (a *Application) worker(fileChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range fileChan {
		fileInfo, err := os.Stat(path)
		if err != nil {
			a.logger.PrintfError("Error accessing file %s: %v\n", path, err)
			continue
		}

		if err := a.compressFile(path, fileInfo); err != nil {
			a.logger.PrintfError("Error compressing file %s: %v\n", path, err)
		}
	}
}

func (a *Application) compressFile(path string, info os.FileInfo) error {
	a.logger.PrintfVerbose("Found file: %s (Size: %d bytes)\n", path, info.Size())

	mimeType := a.mimeDetector.DetectMimeType(path)

	outputPath := filepath.Join(filepath.Dir(path), "compressed_"+filepath.Base(path))

	compressor, exists := a.serviceLocator.GetCompressor(mimeType)
	if exists {
		result, err := compressor.CompressFile(path, outputPath)
		if err != nil {
			return fmt.Errorf("failed to compress file %s: %v", path, err)
		}

		a.logger.PrintfVerbose("Compressed file %s using %T compressor. Original: %d bytes, Compressed: %d bytes, Savings: %s\n",
			path, compressor, result.OriginalSize, result.CompressedSize, result.SavingsPercentageAsHumanReadable())

		// Store the compression result for summary
		a.compressionResults = append(a.compressionResults, result)

		if a.replaceOriginal && result.IsPositiveSavings() {
			if err := a.replaceOriginalFile(path, outputPath); err != nil {
				return fmt.Errorf("failed to replace original file %s: %v", path, err)
			}

			a.logger.PrintfVerbose("Replaced original file %s with compressed version\n", path)
		}

		if a.replaceOriginal {
			_ = os.Remove(outputPath)
		}
	} else {
		a.logger.PrintfVerbose("No compressor found for file: %s\n", path)
	}

	return nil
}

func (a *Application) printOperationSummary() {
	totalFiles := len(a.compressionResults)
	if totalFiles == 0 {
		a.logger.Println("No files were compressed.")
		return
	}

	var totalOriginalSize int64
	var totalCompressedSize int64
	var successfulCompressions int
	var totalSavings int64

	for _, result := range a.compressionResults {
		totalOriginalSize += result.OriginalSize
		totalCompressedSize += result.CompressedSize
		if result.IsPositiveSavings() {
			successfulCompressions++
			totalSavings += result.OriginalSize - result.CompressedSize
		}
	}

	fmt.Println("\n=== Operation Summary ===")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Successfully compressed: %d\n", successfulCompressions)
	fmt.Printf("Total original size: %s\n", formatSize(totalOriginalSize))
	fmt.Printf("Total compressed size: %s\n", formatSize(totalCompressedSize))

	if successfulCompressions > 0 {
		savingsPercentage := float64(totalSavings) / float64(totalOriginalSize) * 100
		fmt.Printf("Total savings: %s (%.2f%%)\n", formatSize(totalSavings), savingsPercentage)
	} else {
		fmt.Println("No space savings achieved.")
	}
}

func formatSize(size int64) string {
	if size < 0 {
		return "0 B"
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func GetDefaultWorkersCount() int {
	cpuCount := runtime.NumCPU()
	if cpuCount <= 1 {
		return 1
	}

	return cpuCount / 2
}
