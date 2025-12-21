package mime

import (
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) DetectMimeType(filePath string) string {
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		// Fallback to extension-based detection if content detection fails
		return d.detectMimeTypeFromExtension(filePath)
	}

	return mtype.String()
}

func (d *Detector) detectMimeTypeFromExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".txt", ".log":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "text/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	default:
		return "application/octet-stream"
	}
}
