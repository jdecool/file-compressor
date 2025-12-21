package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jdecool/file-compressor/internal/app"
	"github.com/jdecool/file-compressor/internal/compressor"
)

func main() {
	var displayHelp bool
	var isVerbose bool
	var maxWorkers int
	var replaceOriginal bool

	flag.BoolVar(&displayHelp, "help", false, "Show help message")
	flag.BoolVar(&isVerbose, "verbose", false, "Enable verbose output")
	flag.IntVar(&maxWorkers, "workers", app.GetDefaultWorkersCount(), "Set maximum number of workers")
	flag.BoolVar(&replaceOriginal, "replace", false, "Replace original file if compression results in savings")
	flag.Parse()

	var inputPaths = flag.Args()

	if displayHelp {
		printUsage()
		os.Exit(0)
	}

	if len(inputPaths) == 0 {
		printUsage()
		os.Exit(0)
	}

	app := app.NewApplication()
	app.SetVerboseMode(isVerbose)
	app.SetMaxWorkers(maxWorkers)
	app.SetReplaceOriginal(replaceOriginal)
	app.RegisterCompressor(compressor.NewPdfCompressor())
	app.RegisterCompressor(compressor.NewImageCompressor())
	app.Run(inputPaths)
}

func printUsage() {
	fmt.Println("File Compressor CLI")
	fmt.Println("Usage: file-compressor [--help] [--verbose] [--replace] <path1> [<path2> ...]")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  file-compressor file.txt                    # Compress a single file")
	fmt.Println("  file-compressor dir/                       # Compress files in directory")
	fmt.Println("  file-compressor file1.txt dir/             # Multiple paths")
	fmt.Println("  file-compressor --verbose file.txt         # Verbose output")
	fmt.Println("  file-compressor --replace file.txt         # Replace original if savings achieved")
	fmt.Println("  file-compressor --help                    # Show this help message")
}
