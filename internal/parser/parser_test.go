package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCode(t *testing.T) {
	testFilePath := "testdata/test.go"
	want := []CodeChunk{
		{
			Name:        "TestFunction",
			FilePath:    testFilePath,
			StartLine:   3,
			EndLine:     5,
			Parameters:  []Parameter{{Name: "a", Type: "int"}, {Name: "b", Type: "string"}},
			Returns:     "string, error",
			Description: "Test function description",
			Language:    "go",
			Content:     "func TestFunction(a int, b string) (string, error) {\n\t// Test function description\n\treturn \"\", nil\n}",
		},
	}

	p := NewParser()
	got, err := p.parseFile(testFilePath)
	
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

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
	// Create a temporary file for testing
	content := []byte(`func LargeFunction() {
		// First part
		code1
		code2
		// Second part
		code3
		code4
	}`)
	
	tmpFilePath := "testdata/large_test.go"
	err := os.WriteFile(tmpFilePath, content, 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpFilePath)

	p := NewParser()
	chunks, err := p.parseFile(tmpFilePath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(chunks))

	// Verify content is captured
	assert.Contains(t, chunks[0].Content, "LargeFunction")
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
