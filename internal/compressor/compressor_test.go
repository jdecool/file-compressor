package compressor

import (
	"testing"
)

func TestSavingsPercentage(t *testing.T) {
	tests := []struct {
		name           string
		originalSize   int64
		compressedSize int64
		expected       float64
	}{
		{
			name:           "No savings",
			originalSize:   100,
			compressedSize: 100,
			expected:       0.0,
		},
		{
			name:           "Positive savings",
			originalSize:   100,
			compressedSize: 50,
			expected:       50.0,
		},
		{
			name:           "Negative savings",
			originalSize:   50,
			compressedSize: 100,
			expected:       -100.0,
		},
		{
			name:           "Zero original size",
			originalSize:   0,
			compressedSize: 100,
			expected:       0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &CompressionResult{
				OriginalSize:   tt.originalSize,
				CompressedSize: tt.compressedSize,
			}
			got := result.SavingsPercentage()
			if got != tt.expected {
				t.Errorf("SavingsPercentage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsPositiveSavings(t *testing.T) {
	tests := []struct {
		name           string
		originalSize   int64
		compressedSize int64
		expected       bool
	}{
		{
			name:           "Positive savings",
			originalSize:   100,
			compressedSize: 50,
			expected:       true,
		},
		{
			name:           "No savings",
			originalSize:   100,
			compressedSize: 100,
			expected:       false,
		},
		{
			name:           "Negative savings",
			originalSize:   50,
			compressedSize: 100,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &CompressionResult{
				OriginalSize:   tt.originalSize,
				CompressedSize: tt.compressedSize,
			}
			got := result.IsPositiveSavings()
			if got != tt.expected {
				t.Errorf("IsPositiveSavings() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSavedSizeAsHumanReadable(t *testing.T) {
	tests := []struct {
		name           string
		originalSize   int64
		compressedSize int64
		want           string
	}{
		{
			name:           "Zero savings",
			originalSize:   1000,
			compressedSize: 1000,
			want:           "0 B",
		},
		{
			name:           "Negative savings",
			originalSize:   1000,
			compressedSize: 1500,
			want:           "0 B",
		},
		{
			name:           "Small savings in bytes",
			originalSize:   2000,
			compressedSize: 1500,
			want:           "500 B",
		},
		{
			name:           "Savings in KB",
			originalSize:   5 * 1024,
			compressedSize: 2 * 1024,
			want:           "3.00 KB",
		},
		{
			name:           "Savings in MB",
			originalSize:   5 * 1024 * 1024,
			compressedSize: 2 * 1024 * 1024,
			want:           "3.00 MB",
		},
		{
			name:           "Savings in GB",
			originalSize:   5 * 1024 * 1024 * 1024,
			compressedSize: 2 * 1024 * 1024 * 1024,
			want:           "3.00 GB",
		},
		{
			name:           "Partial KB savings",
			originalSize:   3500,
			compressedSize: 1000,
			want:           "2.44 KB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &CompressionResult{
				OriginalSize:   tt.originalSize,
				CompressedSize: tt.compressedSize,
			}
			got := result.SavedSizeAsHumanReadable()
			if got != tt.want {
				t.Errorf("SavedSizeAsHumanReadable() = %q, want %q", got, tt.want)
			}
		})
	}
}
