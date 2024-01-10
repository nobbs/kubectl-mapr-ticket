package version_test

import (
	"testing"

	. "github.com/nobbs/kubectl-mapr-ticket/internal/version"
	"github.com/stretchr/testify/assert"
)

func TestVersionString(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.version.String()

			assert.Equal(t, tt.want, got)
		})
	}
}
