package processor

import (
	"context"
	"errors"
	"sync"
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

	mu 		sync.Mutex
	cancels map[string]context.CancelFunc
}

func (p *JobProcessor) AddJob(job domain.Job) (string, error) {
	job_uuid := uuid.NewString()
	entry := domain.NewJobResult(job_uuid)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	p.cancels[job_uuid] = cancel

	select {
		case p.jobs <- QueuedJob{id: job_uuid, job: job, ctx: ctx}:
			if err := p.Storage.Save(entry); err != nil {
				return "", err
			}

			return job_uuid, nil
		default:
			return "", ErrQueueFull
	}
}

func NewJobProcessor(buffer, workers int, storage storage.Storage) *JobProcessor {
	p := &JobProcessor{
		jobs: make(chan QueuedJob, buffer),
		Storage: storage,
		cancels: make(map[string]context.CancelFunc),
	}

	for range workers {
		go p.worker()
	}

	return p
}

func (p *JobProcessor) worker() {
	for envelope := range p.jobs {
		func() {
			now := time.Now().UTC()
			defer p.cancels[envelope.id]()

			p.Storage.Update(envelope.id, func(r *domain.JobResult) error {
				r.Status = domain.JobStatusRunning
				r.StartedAt = &now
				r.Attempts++

				return nil
			})

			out, err := envelope.job.Execute(envelope.ctx)
			now = time.Now().UTC()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					p.Storage.Update(envelope.id, func(r *domain.JobResult) error {
						r.Status = domain.JobStatusCancelled
						now := time.Now().UTC()
						r.FinishedAt = &now
						return nil
					})
					return
				}

				if updateErr := p.Storage.Update(envelope.id, func(r *domain.JobResult) error {
					r.Status = domain.JobStatusFailed
					r.FinishedAt = &now
					r.Error = err
					return nil
				}); updateErr != nil {
					// TODO: log error
				}
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

func (p *JobProcessor) CancelJob(jobID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cancel, ok := p.cancels[jobID]
	if !ok {
		return errors.New("Job not found or already completed")
	}

	cancel()
	delete(p.cancels, jobID)

	return p.Storage.Update(jobID, func(r *domain.JobResult) error {
		r.Status = domain.JobStatusCancelled
		now := time.Now().UTC()
		r.FinishedAt = &now
		return nil
	})
}
