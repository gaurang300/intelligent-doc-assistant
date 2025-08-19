package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// CodeChunk represents a chunk of code with its metadata
type CodeChunk struct {
	Name        string
	Description string
	Language    string
	Example     string
	Parameters  []Parameter
	Returns     string
	FilePath    string
	StartLine   int
	EndLine     int
}

// Parameter represents a function parameter
type Parameter struct {
	Name        string
	Type        string
	Description string
}

// Parser handles code analysis and chunking
type Parser struct {
	fset *token.FileSet
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// ParseGoFile parses a single Go file and returns code chunks
func ParseGoFile(filePath string) ([]CodeChunk, error) {
	p := NewParser()
	return p.parseFile(filePath)
}

// Parse analyzes the code at the given path and returns code chunks
func (p *Parser) Parse(path string) ([]CodeChunk, error) {
	var chunks []CodeChunk

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !isGoFile(path) {
			return nil
		}

		fileChunks, err := p.parseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		chunks = append(chunks, fileChunks...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return chunks, nil
}

func (p *Parser) parseFile(path string) ([]CodeChunk, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	file, err := parser.ParseFile(p.fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	var chunks []CodeChunk

	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		chunk := CodeChunk{
			Name:     fn.Name.Name,
			Language: "go",
			FilePath: path,
		}

		// Get position information
		start := p.fset.Position(fn.Pos())
		end := p.fset.Position(fn.End())
		chunk.StartLine = start.Line
		chunk.EndLine = end.Line

		// Extract description and parameters from doc comments
		if fn.Doc != nil {
			chunk.Description = strings.TrimSpace(fn.Doc.Text())
		}

		// Extract parameters
		if fn.Type.Params != nil {
			for _, param := range fn.Type.Params.List {
				for _, name := range param.Names {
					chunk.Parameters = append(chunk.Parameters, Parameter{
						Name: name.Name,
						Type: typeToString(param.Type),
					})
				}
			}
		}

		// Extract return type
		if fn.Type.Results != nil {
			var returns []string
			for _, result := range fn.Type.Results.List {
				returns = append(returns, typeToString(result.Type))
			}
			chunk.Returns = strings.Join(returns, ", ")
		}

		chunks = append(chunks, chunk)
		return true
	})

	return chunks, nil
}

func isGoFile(path string) bool {
	return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	default:
		return fmt.Sprintf("%T", expr)
	}
}
