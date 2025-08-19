package api

import (
	"encoding/json"
	"net/http"

	"intelligent-doc-assistant/internal/llm"
	"intelligent-doc-assistant/internal/parser"
	"intelligent-doc-assistant/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	Router  *mux.Router
	Parser  *parser.Parser
	Storage *storage.Store
	LLM     *llm.Client
}

func NewServer() *Server {
	s := &Server{
		Router:  mux.NewRouter(),
		Parser:  parser.NewParser(),
		Storage: storage.NewStore(),
		LLM:     llm.NewClient(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("/ingest", s.handleIngest).Methods("POST")
	s.Router.HandleFunc("/ask", s.handleAsk).Methods("POST")
}

type IngestRequest struct {
	RepoPath string `json:"repoPath"`
}

type AskRequest struct {
	Question string `json:"question"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) {
	var req IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Parse the codebase
	chunks, err := s.Parser.Parse(req.RepoPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse codebase")
		return
	}

	// Store the chunks and their embeddings
	if err := s.Storage.StoreChunks(r.Context(), chunks); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to store code chunks")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]interface{}{"message": "Successfully ingested codebase"},
	})
}

func (s *Server) handleAsk(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Search for relevant chunks
	chunks, err := s.Storage.SearchChunks(r.Context(), req.Question)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search code chunks")
		return
	}

	// Generate answer using LLM
	answer, err := s.LLM.GenerateAnswer(r.Context(), req.Question, chunks)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate answer")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]interface{}{"answer": answer},
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{
		Success: false,
		Error:   message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
