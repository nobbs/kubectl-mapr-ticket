package claim_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/pkg/claim"
)

func TestValidateSortOptions(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name        string
		sortOptions []string
		expectedErr error
	}{
		{
			name:        "empty sort options",
			sortOptions: []string{},
			expectedErr: nil,
		},
		{
			name: "one valid sort option",
			sortOptions: []string{
				"name",
			},
			expectedErr: nil,
		},
		{
			name:        "all valid sort options",
			sortOptions: []string{"namespace", "name", "secretNamespace", "secretName", "volumeName", "volumePath", "volumeHandle", "expiryTime", "age"},
			expectedErr: nil,
		},
		{
			name:        "invalid sort option",
			sortOptions: []string{"invalidOption"},
			expectedErr: fmt.Errorf("invalid sort option: invalidOption. Must be one of: namespace|name|secretNamespace|secretName|volumeName|volumePath|volumeHandle|expiryTime|age"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateSortOptions(test.sortOptions)

			assert.Equal(test.expectedErr, err)
		})
	}
}
