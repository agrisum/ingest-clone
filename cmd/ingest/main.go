package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/agris/ingest-clone/pkg/analyzer"
	"github.com/agris/ingest-clone/pkg/config"
	"github.com/agris/ingest-clone/pkg/formatter"
)

const (
	appName    = "ingest"
	appVersion = "0.1.0"
)

func main() {
	// Parse command line flags
	outputFile := flag.String("o", config.DefaultOutputFile, "Output file")
	includePatterns := flag.String("i", "", "Patterns to include (comma-separated)")
	excludePatterns := flag.String("e", "", "Patterns to exclude (comma-separated)")
	filesList := flag.String("f", "", "Specific files to analyze (comma-separated)")
	maxFileSize := flag.Int64("s", config.DefaultMaxFileSize, "Maximum file size to process in bytes")
	showVersion := flag.Bool("v", false, "Show version information")
	showHelp := flag.Bool("h", false, "Show help")

	// Create aliases for flags
	flag.String("output", config.DefaultOutputFile, "Output file (alias for -o)")
	flag.String("include", "", "Patterns to include (alias for -i)")
	flag.String("exclude", "", "Patterns to exclude (alias for -e)")
	flag.String("files", "", "Specific files to analyze (comma-separated) (alias for -f)")
	flag.Int64("size", config.DefaultMaxFileSize, "Maximum file size to process in bytes (alias for -s)")
	flag.Bool("version", false, "Show version information (alias for -v)")
	flag.Bool("help", false, "Show help (alias for -h)")

	flag.Parse()

	// Show version if requested
	if *showVersion {
		fmt.Printf("%s version %s\n", appName, appVersion)
		return
	}

	// Show help if requested
	if *showHelp {
		printUsage()
		return
	}

	// Create configuration
	cfg := config.NewConfig()
	cfg.MaxFileSize = *maxFileSize
	cfg.OutputFile = *outputFile

	// Parse include/exclude patterns
	if *includePatterns != "" {
		cfg.IncludePatterns = config.ParsePatterns(*includePatterns)
	}

	if *excludePatterns != "" {
		cfg.ExcludePatterns = append(cfg.ExcludePatterns, config.ParsePatterns(*excludePatterns)...)
	}

	// Get source directory/file from args or use current directory as default
	args := flag.Args()
	if len(args) > 0 {
		cfg.Source = args[0]
	}

	// Process based on input type
	var allNodes []*analyzer.FileSystemNode

	// If specific files are provided via -f flag, process them
	if *filesList != "" {
		files := config.ParsePatterns(*filesList)
		for _, file := range files {
			// Verify that each file exists
			if !config.FileExists(file) {
				fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", file)
				continue
			}

			// Process the file
			node, err := analyzer.ProcessPath(file, cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to process '%s': %v\n", file, err)
				continue
			}

			allNodes = append(allNodes, node)
		}

		if len(allNodes) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No valid files were found to process\n")
			os.Exit(1)
		}
	} else {
		// Process the source directory/file specified as positional argument
		if !config.FileExists(cfg.Source) && !config.DirExists(cfg.Source) {
			fmt.Fprintf(os.Stderr, "Error: Source '%s' does not exist\n", cfg.Source)
			os.Exit(1)
		}

		node, err := analyzer.ProcessPath(cfg.Source, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to process '%s': %v\n", cfg.Source, err)
			os.Exit(1)
		}

		allNodes = append(allNodes, node)
	}

	// Prepare output
	output := ""

	// Process each node and add to output
	for i, node := range allNodes {
		// Format the results
		result := formatter.FormatResults(node, cfg)

		// Add separator between multiple files
		if i > 0 {
			output += "\n" + config.Separator + "\n\n"
		}

		// Add formatted content
		output += result.Summary + "\n"
		output += result.DirectoryStructure + "\n"
		output += result.FileContents
	}

	// Write the output to a file
	outputDir := filepath.Dir(cfg.OutputFile)
	if outputDir != "" && outputDir != "." {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create output directory: %v\n", err)
			os.Exit(1)
		}
	}

	err := os.WriteFile(cfg.OutputFile, []byte(output), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Analysis complete! Output written to: %s\n", cfg.OutputFile)
}

// printUsage prints the usage information
func printUsage() {
	fmt.Printf("Usage: %s [options] [source]\n\n", appName)
	fmt.Println("Options:")
	fmt.Println("  -o, --output FILE    Output file (default: digest.txt)")
	fmt.Println("  -i, --include PATTERN Patterns to include (comma-separated)")
	fmt.Println("  -e, --exclude PATTERN Patterns to exclude (comma-separated)")
	fmt.Println("  -f, --files FILES    Specific files to analyze (comma-separated)")
	fmt.Println("  -s, --size SIZE      Maximum file size to process in bytes (default: 10MB)")
	fmt.Println("  -v, --version        Show version information")
	fmt.Println("  -h, --help           Show help")
	fmt.Println("\nExamples:")
	fmt.Println("  ingest                           # Analyze current directory")
	fmt.Println("  ingest /path/to/directory        # Analyze specific directory")
	fmt.Println("  ingest -o output.txt /path/to/dir # Specify output file")
	fmt.Println("  ingest -i \"*.go,*.md\" /path/to/dir # Include specific patterns")
	fmt.Println("  ingest -e \"vendor/,*.tmp\" /path/to/dir # Exclude specific patterns")
	fmt.Println("  ingest -f \"file1.go,file2.go,README.md\" # Analyze specific files")
}
