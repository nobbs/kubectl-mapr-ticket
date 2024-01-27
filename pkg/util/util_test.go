// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/util"
)

func TestValidateSortOptions(t *testing.T) {
	t.Parallel()

	type args struct {
		validSortOptions []string
		sortOptions      []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty sort options",
			args: args{
				validSortOptions: []string{"name", "secret.namespace", "secret.name", "claim.namespace", "claim.name", "volume.path", "volume.handle", "expiration", "age"},
				sortOptions:      []string{},
			},
			wantErr: false,
		},
		{
			name: "one valid sort option",
			args: args{
				validSortOptions: []string{"name", "secret.namespace", "secret.name", "claim.namespace", "claim.name", "volume.path", "volume.handle", "expiration", "age"},
				sortOptions:      []string{"name"},
			},
			wantErr: false,
		},
		{
			name: "all valid sort options",
			args: args{
				validSortOptions: []string{"name", "secret.namespace", "secret.name", "claim.namespace", "claim.name", "volume.path", "volume.handle", "expiration", "age"},
				sortOptions:      []string{"name", "secret.namespace", "secret.name", "claim.namespace", "claim.name", "volume.path", "volume.handle", "expiration", "age"},
			},
			wantErr: false,
		},
		{
			name: "invalid sort option",
			args: args{
				validSortOptions: []string{"name", "secret.namespace", "secret.name", "claim.namespace", "claim.name", "volume.path", "volume.handle", "expiration", "age"},
				sortOptions:      []string{"invalid"},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSortOptions(test.args.validSortOptions, test.args.sortOptions)

			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
