package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/rs/cors"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Source    string                 `json:"source"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type SearchRequest struct {
	Query     string `json:"query"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Size      int    `json:"size"`
	From      int    `json:"from"`
}

type App struct {
	osClient *opensearch.Client
}

func main() {
	osURL := os.Getenv("OPENSEARCH_URL")
	if osURL == "" {
		osURL = "http://opensearch:9200"
	}

	cfg := opensearch.Config{
		Addresses: []string{osURL},
		Username:  os.Getenv("OPENSEARCH_USER"),
		Password:  os.Getenv("OPENSEARCH_PASSWORD"),
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating OpenSearch client: %s", err)
	}

	app := &App{osClient: client}

	router := mux.NewRouter()
	router.HandleFunc("/api/logs", app.ingestLog).Methods("POST")
	router.HandleFunc("/api/logs/search", app.searchLogs).Methods("POST")
	router.HandleFunc("/api/health", app.healthCheck).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func (a *App) ingestLog(w http.ResponseWriter, r *http.Request) {
	var logEntry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if logEntry.Timestamp == "" {
		logEntry.Timestamp = time.Now().Format(time.RFC3339)
	}

	data, _ := json.Marshal(logEntry)
	
	indexName := "logs-" + time.Now().Format("2006.01.02")
	
	req := opensearchapi.IndexRequest{
		Index: indexName,
		Body:  bytes.NewReader(data),
	}

	res, err := req.Do(context.Background(), a.osClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (a *App) searchLogs(w http.ResponseWriter, r *http.Request) {
	var searchReq SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if searchReq.Size == 0 {
		searchReq.Size = 100
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"query_string": map[string]interface{}{
							"query": searchReq.Query,
						},
					},
				},
			},
		},
		"size": searchReq.Size,
		"from": searchReq.From,
		"sort": []interface{}{
			map[string]interface{}{
				"timestamp": map[string]string{
					"order": "desc",
				},
			},
		},
	}

	if searchReq.StartTime != "" && searchReq.EndTime != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = []interface{}{
			map[string]interface{}{
				"range": map[string]interface{}{
					"timestamp": map[string]interface{}{
						"gte": searchReq.StartTime,
						"lte": searchReq.EndTime,
					},
				},
			},
		}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := opensearchapi.SearchRequest{
		Index: []string{"logs-*"},
		Body:  &buf,
	}

	res, err := req.Do(context.Background(), a.osClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	req := opensearchapi.PingRequest{}
	res, err := req.Do(context.Background(), a.osClient)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": res.String()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
