package replay

import (
	"fmt"
	"io"
	"time"

	"sshtrace/internal/session"
)

// Options configures replay behaviour.
type Options struct {
	// Speed multiplier: 1.0 = real-time, 2.0 = double speed, 0 = instant.
	Speed float64
	// Writer is the destination for replayed output. Defaults to os.Stdout.
	Writer io.Writer
}

// Replayer replays a recorded SSH session.
type Replayer struct {
	opts Options
}

// New creates a Replayer with the given options.
func New(opts Options) *Replayer {
	if opts.Speed == 0 {
		opts.Speed = 1.0
	}
	return &Replayer{opts: opts}
}

// Replay writes session events to the configured writer, honouring timing.
func (r *Replayer) Replay(s *session.Session) error {
	if len(s.Events) == 0 {
		return nil
	}

	var prev time.Time
	for i, ev := range s.Events {
		if i > 0 && r.opts.Speed > 0 {
			gap := ev.Timestamp.Sub(prev)
			delay := time.Duration(float64(gap) / r.opts.Speed)
			if delay > 0 {
				time.Sleep(delay)
			}
		}
		prev = ev.Timestamp

		if _, err := fmt.Fprint(r.opts.Writer, ev.Data); err != nil {
			return fmt.Errorf("replay write: %w", err)
		}
	}
	return nil
}
