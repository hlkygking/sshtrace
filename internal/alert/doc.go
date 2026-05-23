// Package alert implements rule-based alerting over SSH session events.
//
// An Evaluator holds a list of Rules. Each Rule specifies a keyword and a
// severity Level (info, warn, crit). When Evaluate is called with a Session,
// every event whose Data field contains the keyword (case-insensitive) produces
// an Alert.
//
// Example:
//
//	rules := []alert.Rule{
//		{Name: "sudo detected", Level: alert.LevelWarn, Keyword: "sudo"},
//		{Name: "rm -rf",        Level: alert.LevelCrit, Keyword: "rm -rf"},
//	}
//	ev := alert.New(rules)
//	alerts := ev.Evaluate(sess)
package alert
