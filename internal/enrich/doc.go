// Package enrich provides post-capture enrichment for SSH sessions.
//
// Enrichment adds derived metadata that is not available at capture time,
// such as reverse-DNS hostnames for remote IP addresses and command
// classification labels that make filtering and alerting easier.
//
// Usage:
//
//	e := enrich.New(
//		enrich.WithHostnameResolution(),
//		enrich.WithCommandClassification(),
//	)
//	enrichedSession, err := e.Apply(s)
//
// Command categories assigned by WithCommandClassification:
//
//	"privilege-escalation"  sudo / su
//	"remote-access"         ssh / scp / sftp
//	"file-read"             cat / less / more / tail / head
//	"file-modify"           rm / mv / cp
//	"network-fetch"         curl / wget
//	"general"               everything else
package enrich
