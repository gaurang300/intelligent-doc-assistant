package storage

const (
	// First, drop the existing table if it exists
	DROP_TABLE_CODE_CHUNKS = `
	DROP TABLE IF EXISTS code_chunks;`

	// Create table with proper vector handling
	CREATE_TABLE_CODE_CHUNKS = `
	CREATE EXTENSION IF NOT EXISTS vector;
	
	CREATE TABLE IF NOT EXISTS code_chunks (
		id SERIAL PRIMARY KEY,
		file_path TEXT NOT NULL,
		chunk_text JSONB NOT NULL,
		embedding vector(768) NOT NULL,  -- Gemini embeddings size
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	
	-- Create an index for vector similarity search
	CREATE INDEX IF NOT EXISTS code_chunks_embedding_idx ON code_chunks 
	USING ivfflat (embedding vector_cosine_ops)
	WITH (lists = 100);`

	// Insert with explicit vector casting
	INSERT_CODE_CHUNK = `
	INSERT INTO code_chunks (file_path, chunk_text, embedding)
	VALUES ($1, $2::jsonb, $3::vector);`

	// Search using cosine similarity
	SEARCH_SIMILAR_CHUNKS = `
	SELECT file_path, chunk_text, 
		   1 - (embedding <=> $1::vector) as similarity
	FROM code_chunks
	WHERE 1 - (embedding <=> $1::vector) > 0.7  -- Similarity threshold
	ORDER BY embedding <=> $1::vector
	LIMIT $2;`
)

const (
	InsertFileMetadata = `
        INSERT INTO file_metadata (path, type, size, last_modified)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `

	InsertFunctionDetails = `
        INSERT INTO functions (file_id, name, params, returns, comment)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
    `

	GetFunctionByName = `
        SELECT id, file_id, name, params, returns, comment
        FROM functions
        WHERE name = $1;
    `

	SearchDocs = `
        SELECT id, title, content
        FROM documentation
        WHERE to_tsvector('english', content) @@ plainto_tsquery($1);
    `
)
