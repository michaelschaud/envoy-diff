package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseThresholdConfig parses a threshold spec of the form:
//
//	"min:10"         — only keep values >= 10
//	"max:100"        — only keep values <= 100
//	"min:10,max:100" — keep values in [10, 100]
//
// Non-numeric env values are always preserved.
func ParseThresholdConfig(spec string) (ThresholdConfig, error) {
	var cfg ThresholdConfig
	if strings.TrimSpace(spec) == "" {
		return cfg, nil
	}
	parts := strings.Split(spec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return cfg, fmt.Errorf("threshold: invalid segment %q, expected key:value", part)
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return cfg, fmt.Errorf("threshold: invalid number %q in segment %q", val, part)
		}
		switch strings.ToLower(key) {
		case "min":
			cfg.Min = &f
		case "max":
			cfg.Max = &f
		default:
			return cfg, fmt.Errorf("threshold: unknown key %q, expected min or max", key)
		}
	}
	return cfg, nil
}
