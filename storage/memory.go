package storage

import (
	"errors"
	"sync"

	"example.com/jobqueue/domain"
)

var ErrJobNotFound = errors.New("Job not found")

type MemoryStorage struct {
	mu   sync.RWMutex
	jobs map[string]*domain.JobResult
}

func NewInMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		jobs: make(map[string]*domain.JobResult),
	}
}

func (s *MemoryStorage) Save(r *domain.JobResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[r.JobID] = r
	return nil
}

func (s *MemoryStorage) Get(jobID string) (domain.JobResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.jobs[jobID]
	if !ok {
		return domain.JobResult{}, ErrJobNotFound
	}

	return *r, nil
}

func (s *MemoryStorage) List() ([]domain.JobResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]domain.JobResult, 0, len(s.jobs))
	for _, r := range s.jobs {
		results = append(results, *r)
	}

	return results, nil
}

func (s *MemoryStorage) Update(jobID string, fn func(*domain.JobResult) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.jobs[jobID]
	if !ok {
		return ErrJobNotFound
	}

	if err := fn(r); err != nil {
		return err
	}

	s.jobs[jobID] = r
	return nil
}
