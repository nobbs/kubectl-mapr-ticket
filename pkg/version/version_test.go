// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/version"
)

func TestVersionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version Version
		want    string
	}{
		{
			name: "Clean Version",
			version: Version{
				Version: "1.2.3",
				Commit:  "abcdefg",
				Dirty:   false,
			},
			want: "v1.2.3-abcdefg",
		},
		{
			name: "Dirty Version",
			version: Version{
				Version: "1.2.3",
				Commit:  "abcdefg",
				Dirty:   true,
			},
			want: "v1.2.3-abcdefg-dirty",
		},
		{
			name: "No Commit",
			version: Version{
				Version: "1.2.3",
				Commit:  "",
				Dirty:   false,
			},
			want: "v1.2.3",
		},
		{
			name: "Full Hash",
			version: Version{
				Version: "1.2.3",
				Commit:  "abcdefg1234567890",
				Dirty:   false,
			},
			want: "v1.2.3-abcdefg",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.version.String()

			assert.Equal(t, test.want, got)
		})
	}
}
