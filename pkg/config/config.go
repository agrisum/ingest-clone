package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Constants for default values
const (
	DefaultMaxFileSize  = 10 * 1024 * 1024 // 10 MB
	DefaultOutputFile   = "digest.txt"
	DefaultDirDepth     = 20
	DefaultMaxFiles     = 10000
	DefaultMaxTotalSize = 500 * 1024 * 1024 // 500 MB
	Separator           = "================================================"
)

// Config holds the application configuration
type Config struct {
	// Source directory or file to analyze
	Source string

	// Output file path
	OutputFile string

	// Maximum file size to process in bytes
	MaxFileSize int64

	// Patterns to include (comma-separated)
	IncludePatterns []string

	// Patterns to exclude (comma-separated)
	ExcludePatterns []string

	// Maximum directory depth to traverse
	MaxDirDepth int

	// Maximum number of files to process
	MaxFiles int

	// Maximum total size in bytes
	MaxTotalSize int64
}

// Stats tracks statistics during file processing
type Stats struct {
	TotalFiles int
	TotalSize  int64
	FileCount  int
	DirCount   int
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		Source:          ".",
		OutputFile:      DefaultOutputFile,
		MaxFileSize:     DefaultMaxFileSize,
		IncludePatterns: []string{},
		ExcludePatterns: getDefaultExcludePatterns(),
		MaxDirDepth:     DefaultDirDepth,
		MaxFiles:        DefaultMaxFiles,
		MaxTotalSize:    DefaultMaxTotalSize,
	}
}

// ShouldInclude determines if the given path should be included based on patterns
func (c *Config) ShouldInclude(path string) bool {
	// If no include patterns are specified, include everything by default
	if len(c.IncludePatterns) == 0 {
		return !c.ShouldExclude(path)
	}

	// Check if the path matches any include pattern
	for _, pattern := range c.IncludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}

		// Check for directory patterns like "dir/"
		if strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, strings.TrimSuffix(pattern, "/")) {
			return true
		}
	}

	// If include patterns are specified but none matched, exclude the path
	return false
}

// ShouldExclude determines if the given path should be excluded based on patterns
func (c *Config) ShouldExclude(path string) bool {
	// Check if the path matches any exclude pattern
	for _, pattern := range c.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}

		// Check for directory patterns like "dir/"
		if strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, strings.TrimSuffix(pattern, "/")) {
			return true
		}
	}

	return false
}

// ParsePatterns splits a comma-separated string into a slice of patterns
func ParsePatterns(patterns string) []string {
	if patterns == "" {
		return []string{}
	}

	result := []string{}
	for _, p := range strings.Split(patterns, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

// AbsPath returns the absolute path of a given path
func AbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absPath
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// getDefaultExcludePatterns returns the default patterns to exclude
func getDefaultExcludePatterns() []string {
	return []string{
		// Version control
		".git", ".git/", ".svn", ".svn/", ".hg", ".hg/",
		".gitignore", ".gitattributes", ".gitmodules",

		// Build directories and artifacts
		"node_modules", "node_modules/",
		"vendor", "vendor/",
		"dist", "dist/",
		"build", "build/",

		// Binary and package files
		"*.exe", "*.dll", "*.so", "*.dylib",
		"*.o", "*.obj", "*.a", "*.lib",
		"*.jar", "*.war", "*.ear", "*.zip",
		"*.tar.gz", "*.rar",

		// IDE files
		".idea/", ".vscode/", ".vs/",
		"*.swp", "*.swo",

		// Temporary and cache files
		"*.tmp", "*.temp",
		".cache/", ".sass-cache/",
		".DS_Store", "Thumbs.db",

		// Log files
		"*.log",
	}
}
