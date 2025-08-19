package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"intelligent-doc-assistant/internal/parser"
)

type Generator struct {
	templatesDir string
}

func NewGenerator(templatesDir string) *Generator {
	return &Generator{
		templatesDir: templatesDir,
	}
}

// GenerateDocumentation generates documentation for the given code chunks
func (g *Generator) GenerateDocumentation(chunks []parser.CodeChunk, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Load documentation templates
	funcTemplate, err := template.ParseFiles(filepath.Join(g.templatesDir, "function.md"))
	if err != nil {
		return fmt.Errorf("failed to parse function template: %w", err)
	}

	// Generate documentation for each chunk
	for _, chunk := range chunks {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.md", chunk.Name))
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()

		if err := funcTemplate.Execute(file, chunk); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
	}

	// Generate index file
	indexPath := filepath.Join(outputDir, "index.md")
	index, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer index.Close()

	indexTemplate, err := template.ParseFiles(filepath.Join(g.templatesDir, "index.md"))
	if err != nil {
		return fmt.Errorf("failed to parse index template: %w", err)
	}

	if err := indexTemplate.Execute(index, map[string]interface{}{
		"Chunks": chunks,
	}); err != nil {
		return fmt.Errorf("failed to execute index template: %w", err)
	}

	return nil
}
