package notify

import "fmt"

// MultiNotifier fans a single Message out to multiple Notifiers.
// All notifiers are called; errors are collected and returned together.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMulti returns a Notifier that delivers to all provided notifiers.
func NewMulti(notifiers ...Notifier) Notifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Send delivers msg to every registered Notifier.
// If one or more fail, a combined error is returned.
func (m *MultiNotifier) Send(msg Message) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Send(msg); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("notify: %d error(s): %v", len(errs), errs)
}
