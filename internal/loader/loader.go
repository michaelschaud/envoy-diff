package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnvFile reads a .env file and returns a map of key-value pairs.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
// Values may optionally be quoted with single or double quotes, which are stripped.
func LoadEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: open %q: %w", path, err)
	}
	defer f.Close()

	envs := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("loader: %q line %d: invalid format %q", path, lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := stripQuotes(strings.TrimSpace(parts[1]))

		if key == "" {
			return nil, fmt.Errorf("loader: %q line %d: empty key", path, lineNum)
		}

		envs[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scan %q: %w", path, err)
	}

	return envs, nil
}

// stripQuotes removes surrounding single or double quotes from a value,
// if the value starts and ends with the same quote character.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
