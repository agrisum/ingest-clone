package formatter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/agris/ingest-clone/pkg/analyzer"
	"github.com/agris/ingest-clone/pkg/config"
)

// AnalysisResult holds the formatted analysis results
type AnalysisResult struct {
	Summary            string // Summary of the analysis
	DirectoryStructure string // Tree-like representation of the directory structure
	FileContents       string // Contents of the files
}

// FormatResults formats the analysis results
func FormatResults(root *analyzer.FileSystemNode, cfg *config.Config) *AnalysisResult {
	result := &AnalysisResult{}

	// Generate summary
	result.Summary = formatSummary(root, cfg)

	// Generate directory structure
	result.DirectoryStructure = formatDirectoryStructure(root)

	// Generate file contents
	result.FileContents = formatFileContents(root)

	return result
}

// formatSummary generates a summary of the analysis
func formatSummary(node *analyzer.FileSystemNode, cfg *config.Config) string {
	var summary strings.Builder

	if node.IsDir {
		summary.WriteString(fmt.Sprintf("Directory: %s\n\n", node.Name))
		summary.WriteString(fmt.Sprintf("Files analyzed: %d\n", node.FileCount))
		summary.WriteString(fmt.Sprintf("Total size: %s\n", formatSize(node.Size)))
	} else {
		summary.WriteString(fmt.Sprintf("File: %s\n\n", node.Name))
		summary.WriteString(fmt.Sprintf("Size: %s\n", formatSize(node.Size)))
		summary.WriteString(fmt.Sprintf("Lines: %d\n", strings.Count(node.Content, "\n")+1))
	}

	// Add token count estimation (simplified)
	tokenCount := estimateTokens(node)
	if tokenCount > 0 {
		summary.WriteString(fmt.Sprintf("\nEstimated tokens: %s\n", formatTokenCount(tokenCount)))
	}

	return summary.String()
}

// formatDirectoryStructure generates a tree-like representation of the directory structure
func formatDirectoryStructure(node *analyzer.FileSystemNode) string {
	var builder strings.Builder
	builder.WriteString("Directory structure:\n")

	if node.IsDir {
		prefix := ""
		isLast := true
		buildTree(node, prefix, isLast, &builder)
	} else {
		builder.WriteString(fmt.Sprintf("└── %s\n", node.Name))
	}

	return builder.String()
}

// buildTree recursively builds a tree representation
func buildTree(node *analyzer.FileSystemNode, prefix string, isLast bool, builder *strings.Builder) {
	// Add the current node to the tree
	currentPrefix := "└── "
	if !isLast {
		currentPrefix = "├── "
	}

	// Add trailing slash for directories
	name := node.Name
	if node.IsDir {
		name += "/"
	}

	builder.WriteString(fmt.Sprintf("%s%s%s\n", prefix, currentPrefix, name))

	// If this is not a directory or has no children, return
	if !node.IsDir || len(node.Children) == 0 {
		return
	}

	// Prepare the prefix for children
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Process children
	for i, child := range node.Children {
		isChildLast := i == len(node.Children)-1
		buildTree(child, childPrefix, isChildLast, builder)
	}
}

// formatFileContents formats the contents of all files
func formatFileContents(node *analyzer.FileSystemNode) string {
	var builder strings.Builder

	if !node.IsDir {
		// For a single file, just return its content with a header
		builder.WriteString(formatFileContent(node))
		return builder.String()
	}

	// For a directory, recursively format all files
	formatDirectoryContent(node, &builder)

	return builder.String()
}

// formatDirectoryContent recursively formats the contents of a directory
func formatDirectoryContent(node *analyzer.FileSystemNode, builder *strings.Builder) {
	if !node.IsDir {
		builder.WriteString(formatFileContent(node))
		return
	}

	// Recursively process children
	for _, child := range node.Children {
		formatDirectoryContent(child, builder)
	}
}

// formatFileContent formats the content of a file
func formatFileContent(node *analyzer.FileSystemNode) string {
	if node.IsDir {
		return ""
	}

	var builder strings.Builder

	// Add file header
	relPath := filepath.Base(filepath.Dir(node.Path))
	if relPath == "." {
		relPath = ""
	} else {
		relPath += "/"
	}

	builder.WriteString(fmt.Sprintf("%s\nFILE: %s%s\n%s\n",
		config.Separator, relPath, node.Name, config.Separator))

	// Add file content
	builder.WriteString(node.Content)
	builder.WriteString("\n\n")

	return builder.String()
}

// formatSize formats a size in bytes to a human-readable string
func formatSize(size int64) string {
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

// estimateTokens estimates the number of tokens in the node
func estimateTokens(node *analyzer.FileSystemNode) int {
	// Simple estimation: 1 token ≈ 4 characters
	if node.IsDir {
		return countAllCharacters(node) / 4
	}

	return len(node.Content) / 4
}

// countAllCharacters counts the total number of characters in all files
func countAllCharacters(node *analyzer.FileSystemNode) int {
	if !node.IsDir {
		return len(node.Content)
	}

	total := 0
	for _, child := range node.Children {
		if !child.IsDir {
			total += len(child.Content)
		} else {
			total += countAllCharacters(child)
		}
	}

	return total
}

// formatTokenCount formats a token count to a human-readable string
func formatTokenCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	}

	if count < 1000000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}

	return fmt.Sprintf("%.1fM", float64(count)/1000000)
}
