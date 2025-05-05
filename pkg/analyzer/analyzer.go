package analyzer

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/agris/ingest-clone/pkg/config"
)

// FileSystemNode represents a node in the file system tree
type FileSystemNode struct {
	Name      string            // Name of the file or directory
	Path      string            // Full path to the file or directory
	IsDir     bool              // Whether the node is a directory
	Size      int64             // Size of the file in bytes
	Depth     int               // Depth in the directory tree
	Content   string            // File content (if it's a file)
	Children  []*FileSystemNode // Child nodes (if it's a directory)
	FileCount int               // Number of files in this directory and subdirectories
	DirCount  int               // Number of directories in this directory and subdirectories
}

// NewFileSystemNode creates a new FileSystemNode
func NewFileSystemNode(path string, info fs.FileInfo, depth int) *FileSystemNode {
	return &FileSystemNode{
		Name:      info.Name(),
		Path:      path,
		IsDir:     info.IsDir(),
		Size:      info.Size(),
		Depth:     depth,
		Children:  []*FileSystemNode{},
		FileCount: 0,
		DirCount:  0,
	}
}

// ProcessPath analyzes a file or directory and returns a FileSystemNode
func ProcessPath(path string, cfg *config.Config) (*FileSystemNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Create root node
	root := NewFileSystemNode(absPath, info, 0)

	// Create stats object to track file processing stats
	stats := &config.Stats{}

	// Process the node
	if info.IsDir() {
		err = processDirectory(root, cfg, stats)
	} else {
		err = processFile(root, cfg)
	}

	return root, err
}

// processDirectory processes a directory and its contents
func processDirectory(node *FileSystemNode, cfg *config.Config, stats *config.Stats) error {
	// Check if max depth is reached
	if node.Depth >= cfg.MaxDirDepth {
		return nil
	}

	// Read directory entries
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return err
	}

	// Process each entry
	for _, entry := range entries {
		entryPath := filepath.Join(node.Path, entry.Name())

		// Check if we should include this path
		if !cfg.ShouldInclude(entryPath) || cfg.ShouldExclude(entryPath) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue // Skip entries that can't be accessed
		}

		child := NewFileSystemNode(entryPath, info, node.Depth+1)

		if entry.IsDir() {
			// Process subdirectory
			err = processDirectory(child, cfg, stats)
			if err != nil {
				// Log error but continue processing
				continue
			}
			node.DirCount += child.DirCount + 1
			node.FileCount += child.FileCount
			node.Size += child.Size
		} else {
			// Process file
			if stats.TotalFiles >= cfg.MaxFiles {
				continue // Skip if max files limit reached
			}

			if stats.TotalSize+info.Size() > cfg.MaxTotalSize {
				continue // Skip if max total size limit reached
			}

			if info.Size() > cfg.MaxFileSize {
				continue // Skip if file size exceeds limit
			}

			err = processFile(child, cfg)
			if err != nil {
				// Log error but continue processing
				continue
			}

			node.FileCount++
			node.Size += child.Size
			stats.TotalFiles++
			stats.TotalSize += child.Size
		}

		// Add child to node
		node.Children = append(node.Children, child)
	}

	// Sort children for consistent output
	sortChildren(node)

	return nil
}

// processFile reads and processes a file
func processFile(node *FileSystemNode, cfg *config.Config) error {
	// Skip if file is too large
	if node.Size > cfg.MaxFileSize {
		node.Content = "[File too large]"
		return nil
	}

	// Check if file is binary
	if isBinaryFile(node.Path) {
		node.Content = "[Binary file]"
		return nil
	}

	// Read file content
	content, err := os.ReadFile(node.Path)
	if err != nil {
		node.Content = "[Error reading file]"
		return err
	}

	node.Content = string(content)
	return nil
}

// isBinaryFile checks if a file is likely binary
func isBinaryFile(path string) bool {
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
		return true
	}

	// Check for null bytes in the first 512 bytes
	file, err := os.Open(path)
	if err != nil {
		return true // If we can't read the file, assume it's binary
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return true
	}

	// Look for null bytes, which indicate a binary file
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}

// sortChildren sorts the children of a node
func sortChildren(node *FileSystemNode) {
	// Sort children: first directories (alphabetically), then files (alphabetically)
	if len(node.Children) <= 1 {
		return
	}

	// Custom sorting based on type (directory first) and name
	for i := 0; i < len(node.Children)-1; i++ {
		for j := i + 1; j < len(node.Children); j++ {
			// If i is a file and j is a directory, swap them
			if !node.Children[i].IsDir && node.Children[j].IsDir {
				node.Children[i], node.Children[j] = node.Children[j], node.Children[i]
			} else if (node.Children[i].IsDir == node.Children[j].IsDir) &&
				(strings.ToLower(node.Children[i].Name) > strings.ToLower(node.Children[j].Name)) {
				// If both are of the same type, sort alphabetically
				node.Children[i], node.Children[j] = node.Children[j], node.Children[i]
			}
		}
	}
}
