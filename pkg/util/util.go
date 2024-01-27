// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package util

import (
	"fmt"
	"slices"
	"strings"
)

// StringSliceToCommaSeparatedString returns a normalized string representation of a string slice
// separated by commas and spaces, suitable for use in CLI flag usage strings, e.g.
// "one, two, three".
func StringSliceToCommaSeparatedString(slice []string) string {
	return strings.Join(slice, ", ")
}

// ValidateSortOptions validates the specified sort options to ensure that they are valid.
func ValidateSortOptions(validSortOptions, sortOptions []string) error {
	for _, sortOption := range sortOptions {
		if !slices.Contains(validSortOptions, sortOption) {
			return fmt.Errorf("invalid sort option: %s. Must be one of: (%s)", sortOption, StringSliceToCommaSeparatedString(validSortOptions))
		}
	}

	return nil
}
