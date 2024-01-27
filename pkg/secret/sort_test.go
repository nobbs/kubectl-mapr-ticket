// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package secret_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/secret"
)

func TestValidateSortOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		sortOptions []string
		wantErr     bool
	}{
		{
			name:        "empty sort options",
			sortOptions: []string{},
			wantErr:     false,
		},
		{
			name:        "one valid sort option",
			sortOptions: []string{"name"},
			wantErr:     false,
		},
		{
			name:        "all valid sort options",
			sortOptions: []string{"name", "namespace", "mapr.cluster", "mapr.user", "age", "expiration", "npvcs"},
			wantErr:     false,
		},
		{
			name:        "invalid sort option",
			sortOptions: []string{"invalidOption"},
			wantErr:     true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSortOptions(test.sortOptions)

			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
