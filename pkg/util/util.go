package util

import "strings"

// StringSliceToCommaSeparatedString returns a normalized string representation of a string slice
// separated by commas and spaces, suitable for use in CLI flag usage strings, e.g.
// "one, two, three".
func StringSliceToCommaSeparatedString(slice []string) string {
	return strings.Join(slice, ", ")
}
