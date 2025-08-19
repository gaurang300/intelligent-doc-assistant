# Intelligent Documentation Assistant

A smart documentation assistant that uses Google's Gemini AI to process, understand, and retrieve information from your codebase. It creates semantic embeddings of your code and documentation, stores them in a vector database, and provides intelligent natural language querying capabilities.

## Features

- 📚 Codebase Ingestion: Processes and analyzes your entire codebase
- 🧠 Semantic Understanding: Uses Gemini AI for advanced code comprehension
- 🔍 Natural Language Queries: Ask questions about your codebase in plain English
- 🗄️ Vector Storage: Efficient storage and retrieval using pgvector
- 🚀 High Performance: Built with Go for optimal speed and resource usage

## Prerequisites

- Go 1.21 or higher
- PostgreSQL with pgvector extension
- Google Cloud Gemini API key

## Project Structure

```
.
├── api/          # API definitions and handlers
├── bin/          # Compiled binaries
├── cmd/          # Application entry points
│   ├── ingest/   # Codebase ingestion tool
│   └── server/   # API server
├── config/       # Configuration management
├── docs/         # Documentation generation
├── internal/     # Internal packages
│   ├── embeddings/   # Embedding generation using Gemini
│   ├── llm/         # LLM integration with Gemini
│   ├── parser/      # Code parsing and chunking
│   └── storage/     # Vector database operations
└── utils/       # Utility functions
```

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd intelligent-doc-assistant
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up PostgreSQL with pgvector:
   ```sql
   CREATE EXTENSION vector;
   CREATE DATABASE docassistant;
   ```

4. Configure environment variables:
   ```bash
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=docassistant
   GEMINI_API_KEY=your_api_key_here
   SERVER_PORT=8080
   ```

## Usage

### Ingesting a Codebase

1. Build the ingest tool:
   ```bash
   go build -o bin/ingest cmd/ingest/main.go
   ```

2. Run the ingestion:
   ```bash
   ./bin/ingest /path/to/your/codebase
   ```

### Running the Server

1. Build the server:
   ```bash
   go build -o bin/server cmd/server/main.go
   ```

2. Start the server:
   ```bash
   ./bin/server
   ```

The server will start on the configured port (default: 8080).

### Development

For development, you can use the provided VS Code launch configurations:

- **Ingest Codebase**: Debug the ingestion process
- **Run Server**: Start the API server in debug mode
- **Run Current File**: Debug the currently open Go file

## Architecture

1. **Code Parsing**
   - Chunks code and documentation into semantic units
   - Preserves context and relationships between code elements

2. **Embedding Generation**
   - Uses Gemini AI to generate semantic embeddings
   - Captures meaning and relationships in code

3. **Vector Storage**
   - Stores embeddings in PostgreSQL with pgvector
   - Enables efficient similarity search

4. **Query Processing**
   - Converts natural language queries to embeddings
   - Performs semantic search to find relevant code
   - Uses LLM to generate contextual responses

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## License

[MIT License](LICENSE)

## Acknowledgments

- Google Gemini AI for embeddings and language processing
- pgvector for vector similarity search
- The Go community for excellent tools and libraries