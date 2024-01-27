// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	indend = `  `
)

var (
	// CliBinName is the name of the CLI binary as invoked by the user, e.g. "kubectl-mapr-ticket"
	CliBinName = filepath.Base(os.Args[0])
)

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

// stringNormalizer is a helper struct to normalize strings for use in CLI help and usage text
type stringNormalizer struct {
	string
}

// trim returns a stringNormalizer with the string trimmed of leading and trailing whitespace
func (s stringNormalizer) trim() stringNormalizer {
	s.string = strings.TrimSpace(s.string)

	lines := strings.Split(s.string, "\n")

	for i, line := range strings.Split(s.string, "\n") {
		lines[i] = strings.TrimSpace(line)
	}

	s.string = strings.Join(lines, "\n")
	return s
}

// indent returns a stringNormalizer with each line of the string indented by two spaces
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
