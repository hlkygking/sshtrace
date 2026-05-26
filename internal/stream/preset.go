package stream

import (
	"fmt"
	"io"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
)

// WriterHandler returns a Handler that writes a formatted line for each
// event to the provided writer. Useful for quick logging setups.
func WriterHandler(w io.Writer) Handler {
	return func(sess *session.Session, event session.Event) {
		fmt.Fprintf(w, "%s session=%s user=%s kind=%s data=%q\n",
			time.Now().UTC().Format(time.RFC3339),
			sess.ID,
			sess.User,
			event.Kind,
			event.Data,
		)
	}
}

// FilterHandler wraps a Handler and only forwards events whose Kind matches
// one of the provided kinds. An empty kinds list forwards all events.
func FilterHandler(kinds []string, h Handler) Handler {
	if len(kinds) == 0 {
		return h
	}
	allowed := make(map[string]struct{}, len(kinds))
	for _, k := range kinds {
		allowed[k] = struct{}{}
	}
	return func(sess *session.Session, event session.Event) {
		if _, ok := allowed[event.Kind]; ok {
			h(sess, event)
		}
	}
}

// CountingHandler wraps a Handler and increments a counter each time an
// event is received. The counter pointer must not be nil.
func CountingHandler(counter *int64, h Handler) Handler {
	return func(sess *session.Session, event session.Event) {
		*counter++
		if h != nil {
			h(sess, event)
		}
	}
}
