package embeddings

import (
	"context"
	"fmt"

	genai "cloud.google.com/go/ai/generativelanguage/apiv1"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"google.golang.org/api/option"
)

// GeminiClient handles communication with Gemini's embedding API.
type GeminiClient struct {
	client *genai.GenerativeClient
	model  string
}

// NewGeminiClient initializes a new GeminiClient.
func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewGenerativeClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{
		client: client,
		model:  "models/embedding-001",
	}, nil
}

// CreateEmbeddings generates embeddings for the given input text.
func (c *GeminiClient) CreateEmbeddings(input []string) ([][]float32, error) {
	ctx := context.Background()
	var embeddings [][]float32

	for _, text := range input {
		request := &pb.EmbedContentRequest{
			Model: c.model,
			Content: &pb.Content{
				Parts: []*pb.Part{
					{Data: &pb.Part_Text{Text: text}},
				},
			},
		}
		response, err := c.client.EmbedContent(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("embedding request failed: %w", err)
		}
		embeddings = append(embeddings, response.GetEmbedding().GetValues())
	}

	return embeddings, nil
}

// GetEmbedding generates an embedding for a single text
func GetEmbedding(ctx context.Context, text string, apiKey string) ([]float32, error) {
	client, err := genai.NewGenerativeClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	defer client.Close()

	request := &pb.EmbedContentRequest{
		Model: "models/embedding-001",
		Content: &pb.Content{
			Parts: []*pb.Part{
				{Data: &pb.Part_Text{Text: text}},
			},
		},
	}
	response, err := client.EmbedContent(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}

	return response.GetEmbedding().GetValues(), nil
}
