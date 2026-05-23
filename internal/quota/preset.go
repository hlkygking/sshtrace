package quota

// Preset names for common quota configurations.
const (
	PresetStrict  = "strict"
	PresetRelaxed = "relaxed"
	PresetUnlimited = "unlimited"
)

// PresetLimits returns a Limits struct for a named preset.
// Unknown preset names fall back to DefaultLimits.
func PresetLimits(name string) Limits {
	switch name {
	case PresetStrict:
		return Limits{
			MaxSessions:         10,
			MaxEventsPerSession: 500,
			MaxSessionDuration:  1 * hourDuration,
		}
	case PresetRelaxed:
		return Limits{
			MaxSessions:         500,
			MaxEventsPerSession: 50000,
			MaxSessionDuration:  24 * hourDuration,
		}
	case PresetUnlimited:
		return Limits{}
	default:
		return DefaultLimits()
	}
}

// hourDuration is a typed alias to avoid importing time in the constant block.
const hourDuration = 1_000_000_000 * 3600 // time.Hour in nanoseconds
