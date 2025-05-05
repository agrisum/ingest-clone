package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsTextFile checks if a file is likely a text file
func IsTextFile(path string) bool {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))

	// List of common binary file extensions
	binaryExts := map[string]bool{
		".exe": true, ".dll": true, ".so": true, ".dylib": true,
		".obj": true, ".o": true, ".a": true, ".lib": true,
		".bin": true, ".dat": true, ".db": true, ".sqlite": true,
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true,
		".xlsx": true, ".ppt": true, ".pptx": true, ".zip": true,
		".tar": true, ".gz": true, ".rar": true, ".7z": true,
	}

	if binaryExts[ext] {
		return false
	}

	// Check for null bytes in the first 512 bytes
	file, err := os.Open(path)
	if err != nil {
		return false // If we can't read the file, assume it's binary
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return false
	}

	// Look for null bytes, which indicate a binary file
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return false
		}
	}

	return true
}

// FormatSize formats a size in bytes to a human-readable string
func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// FormatTokenCount formats a token count to a human-readable string
func FormatTokenCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	}

	if count < 1000000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}

	return fmt.Sprintf("%.1fM", float64(count)/1000000)
}
