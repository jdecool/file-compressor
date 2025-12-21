package compressor

import "fmt"

type Compressor interface {
	CompressFile(filePath string, outputPath string) (*CompressionResult, error)

	GetSupportedMimeTypes() []string
}

type CompressionResult struct {
	OriginalFile   string
	CompressedFile string
	OriginalSize   int64
	CompressedSize int64
}

func (r *CompressionResult) SavingsPercentage() float64 {
	if r.OriginalSize == 0 {
		return 0.0
	}

	savings := float64(r.OriginalSize - r.CompressedSize)

	return (savings / float64(r.OriginalSize)) * 100
}

func (r *CompressionResult) SavingsPercentageAsHumanReadable() string {
	return fmt.Sprintf("%.2f%%", r.SavingsPercentage())
}

func (r *CompressionResult) SavedSizeAsHumanReadable() string {
	savedSize := r.OriginalSize - r.CompressedSize
	if savedSize < 0 {
		return "0 B"
	}

	const unit = 1024
	if savedSize < unit {
		return fmt.Sprintf("%d B", savedSize)
	}

	div, exp := int64(unit), 0
	for n := savedSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %cB", float64(savedSize)/float64(div), "KMGTPE"[exp])
}

func (r *CompressionResult) IsPositiveSavings() bool {
	return r.CompressedSize < r.OriginalSize
}
