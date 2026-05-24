package mask

// PresetMasker returns a pre-configured Masker for common presets.
//
// Supported preset names:
//
//	"default"    – masks "password" and "token" fields
//	"strict"     – masks "password", "token", "key", "secret", and "credential"
//	"minimal"    – masks "password" only
//
func PresetMasker(name string) (*Masker, error) {
	presets := map[string][]string{
		"default": {"password", "token"},
		"strict":  {"password", "token", "key", "secret", "credential"},
		"minimal": {"password"},
	}
	fields, ok := presets[name]
	if !ok {
		fields = presets["default"]
	}
	return New(fields, "")
}
