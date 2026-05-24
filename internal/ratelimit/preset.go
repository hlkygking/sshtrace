package ratelimit

import "time"

// PresetConfig returns a named preset configuration.
// Known presets: "strict", "relaxed", "unlimited".
// Unknown names fall back to DefaultConfig.
func PresetConfig(name string) Config {
	switch name {
	case "strict":
		return Config{
			MaxConnections: 3,
			Window:         time.Minute,
		}
	case "relaxed":
		return Config{
			MaxConnections: 30,
			Window:         time.Minute,
		}
	case "unlimited":
		return Config{
			MaxConnections: 1<<31 - 1,
			Window:         time.Hour * 24 * 365,
		}
	default:
		return DefaultConfig()
	}
}
