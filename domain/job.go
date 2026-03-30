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
	Error error
}

type JobResultResponse struct {
	JobID        string     `json:"job_id"`
	Status       JobStatus  `json:"status"`
	Attempts     int        `json:"attempts"`

	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`

	Output       string     `json:"output,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
}

func ToJobResponse(r JobResult) JobResultResponse {
	var errMsg string
	if r.Error != nil {
		errMsg = r.Error.Error()
	}

	return JobResultResponse{
		JobID:        r.JobID,
		Status:       r.Status,
		Attempts:     r.Attempts,
		CreatedAt:    r.CreatedAt,
		StartedAt:    r.StartedAt,
		FinishedAt:   r.FinishedAt,
		Output:       r.Output,
		ErrorMessage: errMsg,
	}
}

func ToJobResponseList(jobs []JobResult) []JobResultResponse {
	responses := make([]JobResultResponse, 0, len(jobs))
	for _, job := range jobs {
		responses = append(responses, ToJobResponse(job))
	}

	return responses
}

type Job interface {
	Execute(ctx context.Context) (string, error)
	Type() string
}
