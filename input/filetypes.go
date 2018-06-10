package input

import (
	"path/filepath"
	"strings"
)

// FileType is an enum/constant helper type
type FileType string

const (
	// FileTypeJSON represents JSON files
	FileTypeJSON = "JSON"

	// FileTypeYAML represents YAML files
	FileTypeYAML = "YAML"

	//FileTypeUnknown represents we don't now (yet)
	FileTypeUnknown = ""
)

// DetectFileType tries to detect the corresponding filetype for a file extension
func DetectFileType(path string) FileType {
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)

	switch ext {
	case ".json":
		return FileTypeJSON

	case ".yaml", ".yml":
		return FileTypeYAML

	default:
		return FileTypeUnknown
	}
}
