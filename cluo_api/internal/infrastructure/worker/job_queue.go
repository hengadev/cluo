package worker

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// JobQueue is an in-memory queue for job IDs using Go channels.
type JobQueue struct {
	jobs  chan uuid.UUID
	logger *slog.Logger
}

// NewJobQueue creates a new job queue with the specified size.
func NewJobQueue(size int, logger *slog.Logger) *JobQueue {
	return &JobQueue{
		jobs:  make(chan uuid.UUID, size),
		logger: logger,
	}
}

// Enqueue adds a job ID to the queue.
// Returns false if the queue is full.
func (q *JobQueue) Enqueue(jobID uuid.UUID) bool {
	select {
	case q.jobs <- jobID:
		q.logger.Debug("job enqueued", "jobID", jobID)
		return true
	default:
		// Queue is full
		q.logger.Warn("job queue full, job not enqueued", "jobID", jobID)
		return false
	}
}

// Dequeue removes and returns a job ID from the queue.
// Blocks until a job is available or the context is cancelled.
func (q *JobQueue) Dequeue(ctx context.Context) (uuid.UUID, bool) {
	select {
	case jobID := <-q.jobs:
		q.logger.Debug("job dequeued", "jobID", jobID)
		return jobID, true
	case <-ctx.Done():
		return uuid.Nil, false
	}
}

// Size returns the current number of jobs in the queue.
func (q *JobQueue) Size() int {
	return len(q.jobs)
}

// Close closes the queue.
func (q *JobQueue) Close() {
	close(q.jobs)
}
