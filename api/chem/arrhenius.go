package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"plenum/pkg/chem"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	A, err := strconv.ParseFloat(q.Get("A"), 64)
	if err != nil {
		http.Error(w, "invalid A", http.StatusBadRequest)
		return
	}
	Ea, err := strconv.ParseFloat(q.Get("Ea"), 64)
	if err != nil {
		http.Error(w, "invalid Ea", http.StatusBadRequest)
		return
	}
	Tmin, err := strconv.ParseFloat(q.Get("Tmin"), 64)
	if err != nil {
		http.Error(w, "invalid Tmin", http.StatusBadRequest)
		return
	}
	Tmax, err := strconv.ParseFloat(q.Get("Tmax"), 64)
	if err != nil || Tmax <= Tmin {
		http.Error(w, "invalid Tmax", http.StatusBadRequest)
		return
	}
	points := 25
	if p := q.Get("points"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 1 {
			points = parsed
		}
	}
	resp := chem.SweepArrhenius(A, Ea, Tmin, Tmax, points)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
