package domain

import (
	"encoding/json"
	"time"
)

type JobBuilder func(payload json.RawMessage) (Job, error)


func NewJobResult(id string) *JobResult {
	return &JobResult{
		JobID: id,
		Status: JobStatusPending,
		CreatedAt: time.Now(),
	}
}

