package telemetry

import (
	"errors"
	"sync"
)

type Update struct {
	JobID        string             `json:"jobId" validate:"required"`
	Iteration    int                `json:"iteration" validate:"required,gte=0"`
	Courant      float64            `json:"courant"`
	Residuals    map[string]float64 `json:"residuals"`
	TimestampUTC string             `json:"timestampUtc"`
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]Update
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: map[string]Update{}}
}

func (s *MemoryStore) Set(update Update) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[update.JobID] = update
}

func (s *MemoryStore) Get(jobID string) (Update, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.data[jobID]
	return u, ok
}


func (u Update) Validate() error {
	if u.JobID == "" {
		return errors.New("jobId is required")
	}
	if u.Iteration < 0 {
		return errors.New("iteration must be non-negative")
	}
	return nil
}
