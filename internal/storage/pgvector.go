package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"intelligent-doc-assistant/config"
	"intelligent-doc-assistant/internal/embeddings"
	"intelligent-doc-assistant/internal/parser"

	_ "github.com/lib/pq"
)

// joinFloat32s converts a slice of float32 to a comma-separated string
func joinFloat32s(v []float32) string {
	s := make([]string, len(v))
	for i, f := range v {
		s[i] = fmt.Sprintf("%.6f", f)
	}
	return strings.Join(s, ",")
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Chunk      parser.CodeChunk
	Similarity float64
}

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
	// Drop existing table
	if _, err := db.Exec(DROP_TABLE_CODE_CHUNKS); err != nil {
		return fmt.Errorf("failed to drop existing table: %w", err)
	}

	// Create new table with proper vector support
	if _, err := db.Exec(CREATE_TABLE_CODE_CHUNKS); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func (s *Store) StoreChunks(ctx context.Context, chunks []parser.CodeChunk) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // will be ignored if tx.Commit() is called

	stmt, err := tx.PrepareContext(ctx, INSERT_CODE_CHUNK)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, chunk := range chunks {
		// Generate embeddings for the chunk
		text := fmt.Sprintf("%s\n%s", chunk.Name, chunk.Description)
		fmt.Printf("Generating embedding for chunk: %s\n", text) // Debug log

		embeddings, err := s.embedder.CreateEmbeddings([]string{text})
		if err != nil {
			return fmt.Errorf("failed to generate embeddings for %s: %w", chunk.FilePath, err)
		}

		if len(embeddings) == 0 {
			return fmt.Errorf("no embeddings generated for chunk in file %s", chunk.FilePath)
		}

		if len(embeddings[0]) != 768 {
			return fmt.Errorf("unexpected embedding dimension %d for file %s", len(embeddings[0]), chunk.FilePath)
		}

		// Serialize chunk data as JSONB
		chunkData, err := json.Marshal(chunk)
		if err != nil {
			return fmt.Errorf("failed to marshal chunk for file %s: %w", chunk.FilePath, err)
		}

		// Format embedding vector
		embedding := Vector(embeddings[0])
		encodedEmbedding := fmt.Sprintf("[%s]", joinFloat32s(embedding))

		fmt.Printf("Inserting chunk for file: %s\n", chunk.FilePath) // Debug log
		fmt.Printf("Chunk data length: %d bytes\n", len(chunkData))  // Debug log
		fmt.Printf("Embedding length: %d values\n", len(embedding))  // Debug log

		// Execute prepared statement
		_, err = stmt.ExecContext(ctx,
			chunk.FilePath,
			chunkData, // Will be automatically cast to JSONB
			encodedEmbedding,
		)
		if err != nil {
			return fmt.Errorf("failed to insert chunk for file %s: %w", chunk.FilePath, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully stored %d chunks\n", len(chunks)) // Debug log
	return nil
}

func (s *Store) SearchChunks(ctx context.Context, query string) ([]SearchResult, error) {
	// Generate embedding for the query
	embeddings, err := s.embedder.CreateEmbeddings([]string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding generated for query")
	}

	// Search for similar chunks
	embedding := Vector(embeddings[0])
	encodedEmbedding := fmt.Sprintf("[%s]", joinFloat32s(embedding))

	// Use prepared statement for better performance
	stmt, err := s.db.PrepareContext(ctx, SEARCH_SIMILAR_CHUNKS)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare search statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, encodedEmbedding, 5) // Top 5 most relevant chunks
	if err != nil {
		return nil, fmt.Errorf("failed to search chunks: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var (
			filePath   string
			chunkData  []byte
			similarity float64
		)

		if err := rows.Scan(&filePath, &chunkData, &similarity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var chunk parser.CodeChunk
		if err := json.Unmarshal(chunkData, &chunk); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
		}

		results = append(results, SearchResult{
			Chunk:      chunk,
			Similarity: similarity,
		})
	}

	return results, nil
}
