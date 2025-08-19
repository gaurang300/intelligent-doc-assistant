package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"intelligent-doc-assistant/config"
	"intelligent-doc-assistant/internal/embeddings"
	"intelligent-doc-assistant/internal/parser"

	_ "github.com/lib/pq"
)

type Store struct {
	db       *sql.DB
	embedder *embeddings.GeminiClient
}

func NewStore() *Store {
	cfg := config.GetConfig()
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return nil
	}

	if err := initSchema(db); err != nil {
		fmt.Printf("Failed to initialize schema: %v\n", err)
		return nil
	}

	embedder, err := embeddings.NewGeminiClient(cfg.GeminiAPIKey)
	if err != nil {
		fmt.Printf("Failed to create embedder: %v\n", err)
		return nil
	}
	return &Store{
		db:       db,
		embedder: embedder,
	}
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(CREATE_TABLE_CODE_CHUNKS)
	return err
}

func (s *Store) StoreChunks(ctx context.Context, chunks []parser.CodeChunk) error {
	for _, chunk := range chunks {
		// Generate embeddings for the chunk
		text := fmt.Sprintf("%s\n%s", chunk.Name, chunk.Description)
		embeddings, err := s.embedder.CreateEmbeddings([]string{text})
		if err != nil {
			return fmt.Errorf("failed to generate embeddings: %w", err)
		}

		if len(embeddings) == 0 {
			return fmt.Errorf("no embeddings generated for chunk")
		}

		// Serialize chunk data
		chunkData, err := json.Marshal(chunk)
		if err != nil {
			return fmt.Errorf("failed to marshal chunk: %w", err)
		}

		// Store chunk with its embedding
		_, err = s.db.Exec(INSERT_CODE_CHUNK,
			chunk.FilePath,
			string(chunkData),
			embeddings[0], // Use the first embedding since we only send one text
		)
		if err != nil {
			return fmt.Errorf("failed to insert chunk: %w", err)
		}
	}
	return nil
}

func (s *Store) SearchChunks(ctx context.Context, query string) ([]parser.CodeChunk, error) {
	// Generate embedding for the query
	embeddings, err := s.embedder.CreateEmbeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding generated for query")
	}

	// Search for similar chunks
	rows, err := s.db.Query(SEARCH_SIMILAR_CHUNKS, embeddings[0], 5) // Limit to top 5 most relevant chunks
	if err != nil {
		return nil, fmt.Errorf("failed to search chunks: %w", err)
	}
	defer rows.Close()

	var chunks []parser.CodeChunk
	for rows.Next() {
		var filePath, chunkData string
		if err := rows.Scan(&filePath, &chunkData); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var chunk parser.CodeChunk
		if err := json.Unmarshal([]byte(chunkData), &chunk); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
		}

		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
