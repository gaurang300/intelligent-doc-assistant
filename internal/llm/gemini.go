package llm

import (
	"context"
	"fmt"

	"intelligent-doc-assistant/config"
	"intelligent-doc-assistant/internal/parser"

	"google.golang.org/genai"
)

type Client struct {
	genaiClient *genai.Client
}

func NewClient() *Client {
	cfg := config.GetConfig()
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.GeminiAPIKey,
		// Endpoint:   "generativelanguage.googleapis.com:443",
		// UserAgent:  "intelligent-doc-assistant",
	})
	if err != nil {
		fmt.Printf("Failed to create Gemini client: %v", err)
		return &Client{}
	}

	return &Client{
		genaiClient: client,
	}
}

// GenerateAnswer generates a response to a user's question using relevant code chunks
func (c *Client) GenerateAnswer(ctx context.Context, question string, chunks []parser.CodeChunk) (string, error) {
	if c.genaiClient == nil {
		return "", fmt.Errorf("Gemini client not initialized")
	}

	// Build the prompt by including relevant code chunks
	prompt := buildPrompt(question, chunks)

	// Generate content using the model
	result, err := c.genaiClient.Models.GenerateContent(
		ctx,
		"gemini-pro",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate answer: %w", err)
	}

	return result.Text(), nil

}

func buildPrompt(question string, chunks []parser.CodeChunk) string {
	var prompt string
	prompt = fmt.Sprintf("Question: %s\n\nRelevant code context:\n", question)

	for _, chunk := range chunks {
		prompt += fmt.Sprintf("\nFile: %s (Lines %d-%d)\nFunction: %s\nDescription: %s\n",
			chunk.FilePath, chunk.StartLine, chunk.EndLine, chunk.Name, chunk.Description)
	}

	prompt += "\nBased on the code context above, please provide a clear and concise answer to the question."

	return prompt
}
