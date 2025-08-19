package storage

import (
	"context"
	"encoding/json"
	"testing"

	"intelligent-doc-assistant/internal/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmbedder is a mock implementation of the embedding client
type MockEmbedder struct {
	mock.Mock
}

func (m *MockEmbedder) CreateEmbeddings(texts []string) ([][]float32, error) {
	args := m.Called(texts)
	return args.Get(0).([][]float32), args.Error(1)
}

// MockDB is a mock implementation of the database
type MockDB struct {
	mock.Mock
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	mockArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	return mockArgs.Get(0), mockArgs.Error(1)
}

func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	mockArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	return mockArgs.Get(0), mockArgs.Error(1)
}

func TestStoreChunks(t *testing.T) {
	tests := []struct {
		name    string
		chunks  []parser.CodeChunk
		setup   func(*MockEmbedder, *MockDB)
		wantErr bool
	}{
		{
			name: "successful storage",
			chunks: []parser.CodeChunk{
				{
					Name:        "TestFunction",
					FilePath:    "/test/path.go",
					StartLine:   1,
					EndLine:     10,
					Description: "Test function description",
				},
			},
			setup: func(me *MockEmbedder, md *MockDB) {
				embedding := []float32{0.1, 0.2, 0.3}
				me.On("CreateEmbeddings", []string{"TestFunction\nTest function description"}).
					Return([][]float32{embedding}, nil)

				expectedData, _ := json.Marshal(parser.CodeChunk{
					Name:        "TestFunction",
					FilePath:    "/test/path.go",
					StartLine:   1,
					EndLine:     10,
					Description: "Test function description",
				})

				md.On("ExecContext",
					mock.Anything,
					INSERT_CODE_CHUNK,
					"/test/path.go",
					expectedData,
					"[0.100000,0.200000,0.300000]",
				).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:    "empty chunks",
			chunks:  []parser.CodeChunk{},
			setup:   func(me *MockEmbedder, md *MockDB) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmbedder := new(MockEmbedder)
			mockDB := new(MockDB)
			tt.setup(mockEmbedder, mockDB)

			store := &Store{
				db:       mockDB,
				embedder: mockEmbedder,
			}

			ctx := context.Background()
			err := store.StoreChunks(ctx, tt.chunks)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockEmbedder.AssertExpectations(t)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestSearchChunks(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		setup   func(*MockEmbedder, *MockDB)
		want    []SearchResult
		wantErr bool
	}{
		{
			name:  "successful search",
			query: "test query",
			setup: func(me *MockEmbedder, md *MockDB) {
				embedding := []float32{0.1, 0.2, 0.3}
				me.On("CreateEmbeddings", []string{"test query"}).
					Return([][]float32{embedding}, nil)

				chunk := parser.CodeChunk{
					Name:        "TestFunction",
					FilePath:    "/test/path.go",
					StartLine:   1,
					EndLine:     10,
					Description: "Test function description",
				}
				chunkData, _ := json.Marshal(chunk)

				md.On("QueryContext",
					mock.Anything,
					SEARCH_SIMILAR_CHUNKS,
					"[0.100000,0.200000,0.300000]",
					5,
				).Return([]struct {
					FilePath   string
					ChunkData  []byte
					Similarity float64
				}{
					{
						FilePath:   "/test/path.go",
						ChunkData:  chunkData,
						Similarity: 0.95,
					},
				}, nil)
			},
			want: []SearchResult{
				{
					Chunk: parser.CodeChunk{
						Name:        "TestFunction",
						FilePath:    "/test/path.go",
						StartLine:   1,
						EndLine:     10,
						Description: "Test function description",
					},
					Similarity: 0.95,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEmbedder := new(MockEmbedder)
			mockDB := new(MockDB)
			tt.setup(mockEmbedder, mockDB)

			store := &Store{
				db:       mockDB,
				embedder: mockEmbedder,
			}

			ctx := context.Background()
			got, err := store.SearchChunks(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockEmbedder.AssertExpectations(t)
			mockDB.AssertExpectations(t)
		})
	}
}
