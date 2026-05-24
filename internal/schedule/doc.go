// Package schedule provides time-based scheduling for periodic session
// processing tasks within sshtrace.
//
// A Scheduler manages one or more named Jobs, each associated with a
// fixed interval and a task function. Jobs run concurrently in the
// background and can be stopped gracefully.
//
// Typical use cases include:
//   - Triggering log rotation on a daily schedule
//   - Periodically flushing buffered audit records
//   - Running quota resets at configured intervals
//
// Example:
//
//	s := schedule.New()
//	_ = s.Add("rotate", 24*time.Hour, func() error {
//		return rotator.Run()
//	})
//	s.Start()
//	defer s.Stop()
package schedule
