// Package schedule provides time-based scheduling for session processing tasks.
package schedule

import (
	"errors"
	"time"
)

// Job represents a recurring task to be executed on a schedule.
type Job struct {
	Name     string
	Interval time.Duration
	Task     func() error
	stop     chan struct{}
	done     chan struct{}
}

// Scheduler manages a collection of periodic jobs.
type Scheduler struct {
	jobs []*Job
}

// New creates a new Scheduler instance.
func New() *Scheduler {
	return &Scheduler{}
}

// Add registers a new job with the scheduler.
// Returns an error if name is empty, interval is non-positive, or task is nil.
func (s *Scheduler) Add(name string, interval time.Duration, task func() error) error {
	if name == "" {
		return errors.New("schedule: job name must not be empty")
	}
	if interval <= 0 {
		return errors.New("schedule: interval must be positive")
	}
	if task == nil {
		return errors.New("schedule: task must not be nil")
	}
	s.jobs = append(s.jobs, &Job{
		Name:     name,
		Interval: interval,
		Task:     task,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	})
	return nil
}

// Start begins executing all registered jobs in separate goroutines.
func (s *Scheduler) Start() {
	for _, j := range s.jobs {
		go runJob(j)
	}
}

// Stop signals all jobs to cease execution and waits for them to finish.
func (s *Scheduler) Stop() {
	for _, j := range s.jobs {
		close(j.stop)
		<-j.done
	}
}

// Jobs returns the list of registered jobs.
func (s *Scheduler) Jobs() []*Job {
	return s.jobs
}

func runJob(j *Job) {
	defer close(j.done)
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = j.Task()
		case <-j.stop:
			return
		}
	}
}
