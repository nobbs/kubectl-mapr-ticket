// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/cmd/common"
)

func TestStringSliceToFlagOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		slice    []string
		expected string
	}{
		{
			name:     "EmptySlice",
			slice:    []string{},
			expected: "",
		},
		{
			name:     "SingleElementSlice",
			slice:    []string{"value"},
			expected: "value",
		},
		{
			name:     "MultipleElementSlice",
			slice:    []string{"value1", "value2", "value3"},
			expected: "value1, value2, value3",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := StringSliceToFlagOptions(test.slice)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCliLongDesc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		desc     string
		args     []any
		expected string
	}{
		{
			name:     "EmptyDesc",
			desc:     "",
			args:     []any{},
			expected: "",
		},
		{
			name:     "SingleLineDesc",
			desc:     "This is a single line description",
			args:     []any{},
			expected: "This is a single line description",
		},
		{
			name:     "MultiLineDesc",
			desc:     "This is a multi line description\nwith multiple lines",
			args:     []any{},
			expected: "This is a multi line description\nwith multiple lines",
		},
		{
			name:     "MultiLineDescWithArgs",
			desc:     "This is a multi line description\nwith multiple lines and %s",
			args:     []any{"args"},
			expected: "This is a multi line description\nwith multiple lines and args",
		},
		{
			name: "MultiLineDescWithBackticks",
			desc: `
			This is a multi line description

			with multiple lines and %s
			`,
			args:     []any{"args"},
			expected: "This is a multi line description\n\nwith multiple lines and args",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := CliLongDesc(test.desc, test.args...)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCliExample(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		example  string
		args     []any
		expected string
	}{
		{
			name:     "EmptyExample",
			example:  "",
			args:     []any{},
			expected: "",
		},
		{
			name:     "SingleLineExample",
			example:  "This is a single line example",
			args:     []any{},
			expected: "  This is a single line example",
		},
		{
			name:     "MultiLineExample",
			example:  "This is a multi line example\nwith multiple lines",
			args:     []any{},
			expected: "  This is a multi line example\n  with multiple lines",
		},
		{
			name:     "MultiLineExampleWithArgs",
			example:  "This is a multi line example\nwith multiple lines and %s",
			args:     []any{"args"},
			expected: "  This is a multi line example\n  with multiple lines and args",
		},
		{
			name: "MultiLineExampleWithBackticks",
			example: `
			This is a multi line example

			with multiple lines and %s
			`,
			args:     []any{"args"},
			expected: "  This is a multi line example\n\n  with multiple lines and args",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := CliExample(test.example, test.args...)
			assert.Equal(t, test.expected, result)
		})
	}
}
