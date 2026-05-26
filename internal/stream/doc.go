// Package stream provides a publish-subscribe broker for real-time
// SSH session event streaming.
//
// Typical usage:
//
//	broker := stream.New()
//
//	err := broker.Subscribe("logger", func(s *session.Session, e session.Event) {
//		fmt.Printf("[%s] %s\n", s.ID, e.Data)
//	})
//
//	// Later, when an event is captured:
//	broker.Publish(sess, event)
//
//	// Remove a subscriber when no longer needed:
//	broker.Unsubscribe("logger")
//
// Publish is safe for concurrent use. Subscribe and Unsubscribe
// are also goroutine-safe.
package stream
