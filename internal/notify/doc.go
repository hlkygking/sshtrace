// Package notify provides notification backends for sshtrace.
//
// # Usage
//
// Create a Notifier and call Send with a Message:
//
//	logN := notify.NewLog(os.Stdout)
//	err := logN.Send(notify.Message{
//		SessionID: "abc123",
//		User:      "alice",
//		Event:     "alert",
//		Detail:    "keyword 'sudo' detected",
//		Timestamp: time.Now(),
//	})
//
// # Channels
//
//   - ChannelLog     — writes human-readable lines to any io.Writer
//   - ChannelWebhook — HTTP POST JSON payload to a configured URL
package notify
