package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"kagent-operator/mcp"

	"github.com/gorilla/mux"
)

const (
	defaultPort     = "8080"
	defaultMCPServerURL = "http://localhost:3000/mcp"
)

type AnalyzeRequest struct {
	ClusterID  string                 `json:"clusterId"`
	ToolNames  []string               `json:"toolNames"`
	Parameters map[string]interface{} `json:"parameters"`
}

type AnalyzeResponse struct {
	ClusterID   string                 `json:"clusterId"`
	ToolResults map[string]interface{} `json:"toolResults"`
	Error       string                 `json:"error,omitempty"`
}

type Server struct {
	mcpClient *mcp.MCPClient
	port      string
}

func NewServer(mcpServerURL string, port string) *Server {
	return &Server{
		mcpClient: mcp.NewMCPClient(mcpServerURL),
		port:      port,
	}
}

func (s *Server) analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if len(req.ToolNames) == 0 {
		respondWithError(w, http.StatusBadRequest, "toolNames cannot be empty")
		return
	}

	// Вызываем MCP Server для выполнения tools
	results, err := s.mcpClient.CallTools(r.Context(), req.ToolNames, req.Parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to call MCP Server: "+err.Error())
		return
	}

	response := AnalyzeResponse{
		ClusterID:   req.ClusterID,
		ToolResults: results,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "kagent-operator",
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mcpServerURL := os.Getenv("MCP_SERVER_URL")
	if mcpServerURL == "" {
		mcpServerURL = defaultMCPServerURL
	}

	server := NewServer(mcpServerURL, port)

	r := mux.NewRouter()
	r.HandleFunc("/analyze", server.analyzeHandler).Methods("POST")
	r.HandleFunc("/health", server.healthHandler).Methods("GET")

	log.Printf("Starting server on port %s", port)
	log.Printf("MCP Server URL: %s", mcpServerURL)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
