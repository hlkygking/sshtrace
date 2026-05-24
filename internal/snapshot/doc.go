// Package snapshot provides point-in-time capture and persistence of SSH
// session state. A Snapshot records the full session — including all events
// captured so far — at the moment Take is called.
//
// Snapshots are written as JSON files to a configurable directory and can be
// reloaded at any time via Load. This is useful for checkpointing long-running
// sessions or for creating restore points before applying transformations such
// as redaction or anonymisation.
//
// Basic usage:
//
//	mgr, err := snapshot.New("/var/lib/sshtrace/snapshots")
//	if err != nil { ... }
//
//	snap, err := mgr.Take(sess)
//	if err != nil { ... }
//
//	reloaded, err := mgr.Load(sess.ID)
//	if err != nil { ... }
package snapshot
