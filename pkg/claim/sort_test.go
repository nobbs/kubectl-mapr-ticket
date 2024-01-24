package claim_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/claim"
)

func TestValidateSortOptions(t *testing.T) {
	assert := assert.New(t)

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
			sortOptions: []string{"namespace", "name", "secret.namespace", "secret.name", "volume.name", "volume.path", "volume.handle", "expiration", "age"},
			wantErr:     false,
		},
		{
			name:        "invalid sort option",
			sortOptions: []string{"invalidOption"},
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateSortOptions(test.sortOptions)

			assert.Equal(test.wantErr, err != nil)
		})
	}
}