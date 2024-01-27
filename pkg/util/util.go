// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package util

import "strings"

// StringSliceToCommaSeparatedString returns a normalized string representation of a string slice
// separated by commas and spaces, suitable for use in CLI flag usage strings, e.g.
// "one, two, three".
func StringSliceToCommaSeparatedString(slice []string) string {
	return strings.Join(slice, ", ")
}
