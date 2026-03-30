package processor

import (
	"context"
	"errors"
	"time"

	"example.com/jobqueue/domain"
	"example.com/jobqueue/storage"
	"github.com/google/uuid"
)

var ErrQueueFull = errors.New("Job queue is full")

type QueuedJob struct {
	job domain.Job
	ctx context.Context
	id  string
}

type JobProcessor struct {
	jobs    chan QueuedJob
	Storage storage.Storage
}

func (p JobProcessor) AddJob(job domain.Job) (string, error) {
	job_uuid := uuid.NewString()
	entry := domain.NewJobResult(job_uuid)

	p.Storage.Save(entry)
	select {
		case p.jobs <- QueuedJob{id: job_uuid, job: job, ctx: context.Background()}:
			return job_uuid, nil
		default:
			return "", ErrQueueFull
	}
}

func NewJobProcessor(buffer, workers int, storage storage.Storage) *JobProcessor {
	p := &JobProcessor{
		jobs: make(chan QueuedJob, buffer),
		Storage: storage,
	}

	for range workers {
		go p.worker()
	}

	return p
}

func (p *JobProcessor) worker() {
	for envelope := range p.jobs {
		func() {
			ctx, cancel := context.WithTimeout(envelope.ctx, 120*time.Second)
			now := time.Now().UTC()
			defer cancel()

			p.Storage.Update(envelope.id, func(r *domain.JobResult) error {
				r.Status = domain.JobStatusRunning
				r.StartedAt = &now
				r.Attempts++

				return nil
			})

			out, err := envelope.job.Execute(ctx)
			now = time.Now().UTC()
			if err != nil {
				p.Storage.Update(envelope.id, func(r *domain.JobResult) error {
					r.Status = domain.JobStatusFailed
					r.FinishedAt = &now
					r.ErrorMessage = err.Error()
					return nil
				})
			} else {
				p.Storage.Update(envelope.id, func(r *domain.JobResult) error {
					r.Status = domain.JobStatusCompleted
					r.Output = out
					r.FinishedAt = &now
					return nil
				})
			}
		}()
	}
}
