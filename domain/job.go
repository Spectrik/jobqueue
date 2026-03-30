package domain

import (
	"context"
	"time"
)

type JobResult struct {
	CreatedAt time.Time
	StartedAt *time.Time
	FinishedAt *time.Time
	Attempts int
	JobID string
	Status JobStatus
	Output string
	ErrorMessage string
}

type Job interface {
	Execute(ctx context.Context) (string, error)
	Type() string
}
