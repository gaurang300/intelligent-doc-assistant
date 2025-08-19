package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCode(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		language string
		want     []CodeChunk
		wantErr  bool
	}{
		{
			name: "parse go function",
			content: `package main

func TestFunction(a int, b string) (string, error) {
	// Test function description
	return "", nil
}`,
			language: "go",
			want: []CodeChunk{
				{
					Name:        "TestFunction",
					FilePath:    "test.go",
					StartLine:   3,
					EndLine:     6,
					Parameters:  []string{"a int", "b string"},
					Returns:     []string{"string", "error"},
					Description: "Test function description",
					Language:    "go",
				},
			},
			wantErr: false,
		},
		{
			name: "parse empty file",
			content: `package main

// Empty file for testing`,
			language: "go",
			want:     []CodeChunk{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCode("test.go", tt.content, tt.language)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestChunkCode(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		language string
		maxSize  int
		want     []CodeChunk
		wantErr  bool
	}{
		{
			name: "chunk large function",
			content: `func LargeFunction() {
				// First part
				code1
				code2
				// Second part
				code3
				code4
			}`,
			language: "go",
			maxSize:  100,
			want: []CodeChunk{
				{
					Name:        "LargeFunction part 1",
					FilePath:    "test.go",
					StartLine:   1,
					EndLine:     4,
					Description: "First part",
					Language:    "go",
				},
				{
					Name:        "LargeFunction part 2",
					FilePath:    "test.go",
					StartLine:   4,
					EndLine:     7,
					Description: "Second part",
					Language:    "go",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChunkCode("test.go", tt.content, tt.language, tt.maxSize)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
