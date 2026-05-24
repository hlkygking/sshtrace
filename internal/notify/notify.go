// Package notify provides pluggable notification backends for sshtrace alerts.
// Supported channels: log (stdout), webhook (HTTP POST).
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Channel represents a notification delivery channel.
type Channel string

const (
	ChannelLog     Channel = "log"
	ChannelWebhook Channel = "webhook"
)

// Message holds the data sent to a notification backend.
type Message struct {
	SessionID string    `json:"session_id"`
	User      string    `json:"user"`
	Event     string    `json:"event"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}

// Notifier sends a Message to a destination.
type Notifier interface {
	Send(msg Message) error
}

// LogNotifier writes notifications to an io.Writer.
type LogNotifier struct {
	w io.Writer
}

// NewLog returns a Notifier that writes to w.
func NewLog(w io.Writer) Notifier {
	return &LogNotifier{w: w}
}

func (l *LogNotifier) Send(msg Message) error {
	_, err := fmt.Fprintf(l.w, "[notify] %s user=%s event=%s detail=%s\n",
		msg.Timestamp.Format(time.RFC3339), msg.User, msg.Event, msg.Detail)
	return err
}

// WebhookNotifier posts notifications as JSON to an HTTP endpoint.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhook returns a Notifier that POSTs JSON to url.
func NewWebhook(url string) Notifier {
	return &WebhookNotifier{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (wh *WebhookNotifier) Send(msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("notify: marshal: %w", err)
	}
	resp, err := wh.client.Post(wh.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
