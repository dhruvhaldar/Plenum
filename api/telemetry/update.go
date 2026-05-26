package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"plenum/pkg/telemetry"
)

var store = telemetry.NewMemoryStore()

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req telemetry.Update
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json payload", http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.TimestampUTC = time.Now().UTC().Format(time.RFC3339)
		store.Set(req)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	case http.MethodGet:
		jobID := r.URL.Query().Get("jobId")
		if jobID == "" {
			http.Error(w, "jobId is required", http.StatusBadRequest)
			return
		}
		state, ok := store.Get(jobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(state)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
