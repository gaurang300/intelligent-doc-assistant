package storage

const (
	CREATE_TABLE_CODE_CHUNKS = `
	CREATE TABLE IF NOT EXISTS code_chunks (
		id SERIAL PRIMARY KEY,
		file_path TEXT,
		chunk_text TEXT,
		embedding vector(768)  -- Gemini embeddings size
	);`

	INSERT_CODE_CHUNK = `
	INSERT INTO code_chunks (file_path, chunk_text, embedding)
	VALUES ($1, $2, $3);`

	SEARCH_SIMILAR_CHUNKS = `
	SELECT file_path, chunk_text
	FROM code_chunks
	ORDER BY embedding <-> $1
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
