package rotate

import "time"

// PresetPolicy returns a named rotation policy.
// Known names: "strict", "relaxed", "unlimited".
// Unknown names fall back to DefaultPolicy.
func PresetPolicy(name string) Policy {
	switch name {
	case "strict":
		return Policy{
			MaxAge:   7 * 24 * time.Hour, // 7 days
			MaxCount: 1000,
		}
	case "relaxed":
		return Policy{
			MaxAge:   90 * 24 * time.Hour, // 90 days
			MaxCount: 50000,
		}
	case "unlimited":
		return Policy{
			MaxAge:   0,
			MaxCount: 0,
		}
	default:
		return DefaultPolicy()
	}
}
