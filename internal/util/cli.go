package util

import (
	"fmt"
	"strings"
)

const (
	indend = `  `
)

type stringNormalizer struct {
	string
}

// StringSliceToFlagOptions returns a normalized string representation of a string slice
// separated by commas and spaces, suitable for use in CLI flag usage strings, e.g.
// "one, two, three".
func StringSliceToFlagOptions(slice []string) string {
	return strings.Join(slice, ", ")
}

// CliShortDesc returns a normalized long description of a CLI command, the desc string
// is passed through fmt.Sprintf() with the args as arguments.
func CliLongDesc(desc string, args ...any) string {
	if desc == "" || len(desc) == 0 {
		return ""
	}

	desc = fmt.Sprintf(desc, args...)

	return stringNormalizer{desc}.trim().string
}

// CliExample returns a normalized example of a CLI command, the example string
// is passed through fmt.Sprintf() with the args as arguments.
func CliExample(example string, args ...any) string {
	if example == "" || len(example) == 0 {
		return ""
	}

	example = fmt.Sprintf(example, args...)

	return stringNormalizer{example}.trim().indent().string
}

func (s stringNormalizer) trim() stringNormalizer {
	s.string = strings.TrimSpace(s.string)

	lines := strings.Split(s.string, "\n")

	for i, line := range strings.Split(s.string, "\n") {
		lines[i] = strings.TrimSpace(line)
	}

	s.string = strings.Join(lines, "\n")
	return s
}

func (s stringNormalizer) indent() stringNormalizer {
	lines := strings.Split(s.string, "\n")

	for i, line := range lines {
		if line != "" {
			lines[i] = indend + line
		} else {
			lines[i] = ""
		}
	}

	s.string = strings.Join(lines, "\n")
	return s
}
