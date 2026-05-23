package export

import "fmt"

// AvailableFormats returns a slice of all supported export format strings.
func AvailableFormats() []Format {
	return []Format{FormatJSON, FormatText}
}

// ParseFormat converts a string to a Format, returning an error if unrecognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatJSON, FormatText:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unknown format %q: supported formats are json, text", s)
	}
}

// String implements the Stringer interface for Format.
func (f Format) String() string {
	return string(f)
}
