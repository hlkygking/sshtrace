// Package tag provides lightweight tagging support for SSH sessions captured
// by sshtrace.
//
// Tags are short, lowercase labels (e.g. "prod", "ci", "admin") that can be
// attached to a [session.Session] to enable downstream filtering, reporting,
// and alerting.
//
// # Usage
//
//	tgr := tag.New()
//	if err := tgr.Add(sess, "prod", "admin"); err != nil {
//		log.Fatal(err)
//	}
//	if tag.Has(sess, "prod") {
//		// handle production session
//	}
//	tgr.Remove(sess, "admin")
//
// Tag names must match the pattern [a-z0-9_\-]{1,32}; any other value causes
// [ErrInvalidTag] to be returned.
package tag
