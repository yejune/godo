// Package memory provides memory management utilities for the Do Worker Service.
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/do-focus/worker/internal/db"
	"github.com/do-focus/worker/pkg/models"
)

// Store manages in-memory caching and batch operations.
type Store struct {
	db    db.Adapter
	cache *sync.Map

	// Batch queue for observations
	obsBatch    []models.Observation
	obsMu       sync.Mutex
	batchTicker *time.Ticker
	done        chan struct{}
}

// NewStore creates a new memory store.
func NewStore(adapter db.Adapter) *Store {
	s := &Store{
		db:       adapter,
		cache:    &sync.Map{},
		obsBatch: make([]models.Observation, 0, 100),
		done:     make(chan struct{}),
	}

	// Start batch flusher
	s.batchTicker = time.NewTicker(5 * time.Second)
	go s.batchFlusher()

	return s
}

// Close stops the store and flushes pending data.
func (s *Store) Close() error {
	close(s.done)
	s.batchTicker.Stop()
	return s.flushObservations()
}

// batchFlusher periodically flushes batched observations.
func (s *Store) batchFlusher() {
	for {
		select {
		case <-s.batchTicker.C:
			_ = s.flushObservations()
		case <-s.done:
			return
		}
	}
}

// QueueObservation adds an observation to the batch queue.
func (s *Store) QueueObservation(obs models.Observation) {
	s.obsMu.Lock()
	defer s.obsMu.Unlock()

	s.obsBatch = append(s.obsBatch, obs)

	// Flush immediately if batch is large enough
	if len(s.obsBatch) >= 50 {
		go func() {
			_ = s.flushObservations()
		}()
	}
}

// flushObservations writes all queued observations to the database.
func (s *Store) flushObservations() error {
	s.obsMu.Lock()
	if len(s.obsBatch) == 0 {
		s.obsMu.Unlock()
		return nil
	}

	batch := s.obsBatch
	s.obsBatch = make([]models.Observation, 0, 100)
	s.obsMu.Unlock()

	ctx := context.Background()
	for i := range batch {
		if err := s.db.CreateObservation(ctx, &batch[i]); err != nil {
			// Log error but continue with other observations
			continue
		}
	}

	return nil
}

// CacheKey represents a cache key type.
type CacheKey string

const (
	CacheKeySession CacheKey = "session"
	CacheKeyPlan    CacheKey = "plan"
)

// SetCache stores a value in the cache.
func (s *Store) SetCache(key CacheKey, id string, value interface{}) {
	s.cache.Store(string(key)+":"+id, value)
}

// GetCache retrieves a value from the cache.
func (s *Store) GetCache(key CacheKey, id string) (interface{}, bool) {
	return s.cache.Load(string(key) + ":" + id)
}

// DeleteCache removes a value from the cache.
func (s *Store) DeleteCache(key CacheKey, id string) {
	s.cache.Delete(string(key) + ":" + id)
}

// ClearCache removes all cached values.
func (s *Store) ClearCache() {
	s.cache.Range(func(key, value interface{}) bool {
		s.cache.Delete(key)
		return true
	})
}
