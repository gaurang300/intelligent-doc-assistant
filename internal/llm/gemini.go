package llm

import (
	"context"
	"fmt"
	"strings"

	"intelligent-doc-assistant/config"
	"intelligent-doc-assistant/internal/storage"

	genai "cloud.google.com/go/ai/generativelanguage/apiv1"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"google.golang.org/api/option"
)

type Client struct {
	genaiClient *genai.GenerativeClient
}

func NewClient() *Client {
	cfg := config.GetConfig()
	ctx := context.Background()

	client, err := genai.NewGenerativeClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		fmt.Printf("Failed to create Gemini client: %v", err)
		return &Client{}
	}

	return &Client{
		genaiClient: client,
	}
}

// GenerateAnswer generates a response to a user's question using relevant code chunks and their similarity scores
func (c *Client) GenerateAnswer(ctx context.Context, question string, searchResults []storage.SearchResult) (string, error) {
	if c.genaiClient == nil {
		return "", fmt.Errorf("Gemini client not initialized")
	}

	if len(searchResults) == 0 {
		return "", fmt.Errorf("no relevant code found in the codebase for the question. Please make sure to ingest the codebase first using the /ingest endpoint")
	}

	// Build the prompt by including relevant code chunks and their similarity scores
	prompt := buildPromptWithScores(question, searchResults)

	fmt.Printf("Sending prompt to Gemini: %s\n", prompt) // Debug log

	// Create float32 and int32 pointers for config
	temp := float32(0.3)
	topP := float32(0.8)
	topK := int32(40)
	maxTokens := int32(1024)

	// Create the generate content request
	req := &pb.GenerateContentRequest{
		Model: "models/gemini-2.0-flash-001",
		Contents: []*pb.Content{
			{
				Parts: []*pb.Part{
					{
						Data: &pb.Part_Text{
							Text: prompt,
						},
					},
				},
			},
		},
		GenerationConfig: &pb.GenerationConfig{
			Temperature:     &temp,
			TopP:            &topP,
			TopK:            &topK,
			MaxOutputTokens: &maxTokens,
		},
		SafetySettings: []*pb.SafetySetting{
			{
				Category:  pb.HarmCategory_HARM_CATEGORY_HARASSMENT,
				Threshold: pb.SafetySetting_HARM_BLOCK_THRESHOLD_UNSPECIFIED,
			},
			{
				Category:  pb.HarmCategory_HARM_CATEGORY_HATE_SPEECH,
				Threshold: pb.SafetySetting_HARM_BLOCK_THRESHOLD_UNSPECIFIED,
			},
		},
	}

	// Generate response
	resp, err := c.genaiClient.GenerateContent(ctx, req)
	if err != nil {
		return "", fmt.Errorf("Gemini API error: %w", err)
	}

	if resp == nil {
		return "", fmt.Errorf("no response generated from Gemini")
	}

	var answer strings.Builder
	for _, candidate := range resp.Candidates {
		if candidate != nil && len(candidate.Content.Parts) > 0 {
			for _, part := range candidate.Content.Parts {
				if text := part.GetText(); text != "" {
					answer.WriteString(text)
					answer.WriteString(" ")
				}
			}
		}
	}

	result := strings.TrimSpace(answer.String())
	if result == "" {
		return "", fmt.Errorf("no text content in Gemini response")
	}

	fmt.Printf("Generated answer: %s\n", result) // Debug log
	return result, nil
}

func buildPromptWithScores(question string, results []storage.SearchResult) string {
	var prompt string
	prompt = fmt.Sprintf("Question: %s\n\nRelevant code context (sorted by similarity):\n", question)

	for _, result := range results {
		chunk := result.Chunk
		prompt += fmt.Sprintf("\nFile: %s (Lines %d-%d)\nFunction: %s\nDescription: %s\nRelevance Score: %.2f\n",
			chunk.FilePath, chunk.StartLine, chunk.EndLine, chunk.Name, chunk.Description, result.Similarity)
	}

	prompt += "\nBased on the code context above, with consideration for the relevance scores, please provide a clear and concise answer to the question."

	return prompt
}
