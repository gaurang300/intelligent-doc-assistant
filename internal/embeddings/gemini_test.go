package embeddings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGeminiAPI mocks the Gemini API client
type MockGeminiAPI struct {
	mock.Mock
}

func (m *MockGeminiAPI) EmbedContent(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	return args.Get(0).([]float32), args.Error(1)
}

func TestCreateEmbeddings(t *testing.T) {
	tests := []struct {
		name    string
		texts   []string
		setup   func(*MockGeminiAPI)
		want    [][]float32
		wantErr bool
	}{
		{
			name:  "successful embedding generation",
			texts: []string{"test code", "another test"},
			setup: func(m *MockGeminiAPI) {
				m.On("EmbedContent", mock.Anything, "test code").
					Return([]float32{0.1, 0.2, 0.3}, nil)
				m.On("EmbedContent", mock.Anything, "another test").
					Return([]float32{0.4, 0.5, 0.6}, nil)
			},
			want: [][]float32{
				{0.1, 0.2, 0.3},
				{0.4, 0.5, 0.6},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			texts:   []string{},
			setup:   func(m *MockGeminiAPI) {},
			want:    [][]float32{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := new(MockGeminiAPI)
			tt.setup(mockAPI)

			client := &GeminiClient{
				client: mockAPI,
			}

			got, err := client.CreateEmbeddings(tt.texts)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestGetEmbedding(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		setup   func(*MockGeminiAPI)
		want    []float32
		wantErr bool
	}{
		{
			name: "successful single embedding",
			text: "test code",
			setup: func(m *MockGeminiAPI) {
				m.On("EmbedContent", mock.Anything, "test code").
					Return([]float32{0.1, 0.2, 0.3}, nil)
			},
			want:    []float32{0.1, 0.2, 0.3},
			wantErr: false,
		},
		{
			name: "empty text",
			text: "",
			setup: func(m *MockGeminiAPI) {
				m.On("EmbedContent", mock.Anything, "").
					Return([]float32{}, nil)
			},
			want:    []float32{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := new(MockGeminiAPI)
			tt.setup(mockAPI)

			ctx := context.Background()
			got, err := GetEmbedding(ctx, tt.text, "test-api-key")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockAPI.AssertExpectations(t)
		})
	}
}
