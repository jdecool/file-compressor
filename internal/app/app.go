package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/jdecool/file-compressor/internal/compressor"
	"github.com/jdecool/file-compressor/internal/mime"
	"github.com/jdecool/file-compressor/internal/servicelocator"
)

type Application struct {
	isVerbose       bool
	maxWorkers      int
	replaceOriginal bool
	serviceLocator  *servicelocator.ServiceLocator
	mimeDetector    *mime.Detector
}

func NewApplication() *Application {
	return &Application{
		isVerbose:       false,
		maxWorkers:      GetDefaultWorkersCount(),
		replaceOriginal: false,
		serviceLocator:  servicelocator.NewServiceLocator(),
		mimeDetector:    mime.NewDetector(),
	}
}

func (a *Application) RegisterCompressor(c compressor.Compressor) {
	a.serviceLocator.RegisterCompressor(c)
}

func (a *Application) SetVerboseMode(isVerbose bool) {
	a.isVerbose = isVerbose
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
	if a.isVerbose {
		fmt.Println("Starting file compression process...")
		fmt.Println("Input paths:", inputPaths)
		fmt.Printf("Using %d workers for parallel processing\n", a.maxWorkers)
	}

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
			if a.isVerbose {
				fmt.Printf("Processing path: %s\n", path)
			}

			fileInfo, err := os.Stat(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", path, err)
				continue
			}

			if fileInfo.IsDir() {
				err := a.browseDirectoryAndSendFiles(path, fileChan)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error browsing directory %s: %v\n", path, err)
				}
			} else {
				fileChan <- path
			}
		}
		close(fileChan)
	}()

	<-doneChan

	if a.isVerbose {
		fmt.Println("File compression completed.")
	}
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
			fmt.Fprintf(os.Stderr, "Error accessing file %s: %v\n", path, err)
			continue
		}

		a.compressFile(path, fileInfo)
	}
}

func (a *Application) compressFile(path string, info os.FileInfo) error {
	fmt.Printf("Found file: %s (Size: %d bytes)\n", path, info.Size())

	mimeType := a.mimeDetector.DetectMimeType(path)

	outputPath := filepath.Join(filepath.Dir(path), "compressed_"+filepath.Base(path))

	compressor, exists := a.serviceLocator.GetCompressor(mimeType)
	if exists {
		result, err := compressor.CompressFile(path, outputPath)
		if err != nil {
			return fmt.Errorf("failed to compress file %s: %v", path, err)
		}

		fmt.Printf("Compressed file %s using %T compressor. Original: %d bytes, Compressed: %d bytes, Savings: %s\n",
			path, compressor, result.OriginalSize, result.CompressedSize, result.SavingsPercentageAsHumanReadable())

		if a.replaceOriginal && result.IsPositiveSavings() {
			if err := a.replaceOriginalFile(path, outputPath); err != nil {
				return fmt.Errorf("failed to replace original file %s: %v", path, err)
			}

			fmt.Printf("Replaced original file %s with compressed version\n", path)
		}

		if a.replaceOriginal {
			_ = os.Remove(outputPath)
		}
	} else {
		fmt.Printf("No compressor found for file: %s\n", path)
	}

	return nil
}

func GetDefaultWorkersCount() int {
	cpuCount := runtime.NumCPU()
	if cpuCount <= 1 {
		return 1
	}

	return cpuCount / 2
}
