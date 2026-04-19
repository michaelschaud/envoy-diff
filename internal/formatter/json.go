package formatter

import (
	"encoding/json"
	"io"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// jsonPayload is the serialisable representation of a diff result.
type jsonPayload struct {
	OnlyInStaging    map[string]string `json:"only_in_staging"`
	OnlyInProduction map[string]string `json:"only_in_production"`
	Modified         map[string]valuePair `json:"modified"`
	Identical        map[string]string `json:"identical"`
}

type valuePair struct {
	Staging    string `json:"staging"`
	Production string `json:"production"`
}

// JSONWriter encodes the diff result as indented JSON and writes it to w.
func JSONWriter(w io.Writer, result diff.Result) error {
	modified := make(map[string]valuePair, len(result.Modified))
	for k, pair := range result.Modified {
		modified[k] = valuePair{
			Staging:    pair[0],
			Production: pair[1],
		}
	}

	payload := jsonPayload{
		OnlyInStaging:    result.OnlyInStaging,
		OnlyInProduction: result.OnlyInProduction,
		Modified:         modified,
		Identical:        result.Identical,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
