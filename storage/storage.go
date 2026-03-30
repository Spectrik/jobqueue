package storage

import (
	"example.com/jobqueue/domain"
)

type Storage interface {
	Update(jobID string, fn func(*domain.JobResult) error) error
	Save(r *domain.JobResult) error
	Get(jobID string) (domain.JobResult, error)
	List() ([]domain.JobResult, error)
}
