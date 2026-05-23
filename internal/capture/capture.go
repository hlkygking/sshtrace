package capture

import (
	"io"
	"time"

	"sshtrace/internal/session"
)

// Capturer intercepts reads/writes on an SSH channel and records them
// as session events.
type Capturer struct {
	session *session.Session
	writer  io.Writer
	reader  io.Reader
}

// New creates a Capturer that wraps the given reader/writer pair and
// appends every chunk of data as an event on sess.
func New(sess *session.Session, r io.Reader, w io.Writer) *Capturer {
	return &Capturer{
		session: sess,
		writer:  w,
		reader:  r,
	}
}

// Write forwards p to the underlying writer and records the data as an
// outbound event (server → client).
func (c *Capturer) Write(p []byte) (int, error) {
	n, err := c.writer.Write(p)
	if n > 0 {
		c.session.AddEvent(session.Event{
			Timestamp: time.Now(),
			Direction: session.DirOutput,
			Data:      append([]byte(nil), p[:n]...),
		})
	}
	return n, err
}

// Read forwards to the underlying reader and records the data as an
// inbound event (client → server).
func (c *Capturer) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	if n > 0 {
		c.session.AddEvent(session.Event{
			Timestamp: time.Now(),
			Direction: session.DirInput,
			Data:      append([]byte(nil), p[:n]...),
		})
	}
	return n, err
}
