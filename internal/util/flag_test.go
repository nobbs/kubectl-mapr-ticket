package util_test

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/internal/util"
)

func TestDurationValue_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedVal DurationValue
		expectErr   bool
	}{
		{
			name:        "valid duration string, 1h",
			input:       "1h",
			expectedVal: DurationValue(1 * time.Hour),
			expectErr:   false,
		},
		{
			name:        "valid duration string, 1d",
			input:       "1d",
			expectedVal: DurationValue(24 * time.Hour),
			expectErr:   false,
		},
		{
			name:        "valid duration string, 1w",
			input:       "1w",
			expectedVal: DurationValue(7 * 24 * time.Hour),
			expectErr:   false,
		},
		{
			name:        "valid duration string, 1w2d3h4m5s",
			input:       "1w2d3h4m5s",
			expectedVal: DurationValue(9*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second),
			expectErr:   false,
		},
		{
			name:        "invalid duration string",
			input:       "invalid",
			expectedVal: 0,
			expectErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := new(DurationValue)

			err := d.Set(test.input)

			assert.Equal(t, test.expectedVal, *d)

			switch test.expectErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}

func TestDurationValue_Type(t *testing.T) {
	d := new(DurationValue)

	assert.Equal(t, "duration", d.Type())
}

func TestDurationValueFlag(t *testing.T) {
	var testDuration DurationValue

	testCmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	testCmd.Flags().VarP(&testDuration, "duration", "d", "test duration flag")

	tests := []struct {
		name        string
		args        []string
		expectedVal time.Duration
		expectErr   bool
	}{
		{
			name:        "valid duration string, 1h",
			args:        []string{"--duration", "1h"},
			expectedVal: time.Hour,
			expectErr:   false,
		},
		{
			name:        "valid duration string, 1d",
			args:        []string{"--duration", "1d"},
			expectedVal: time.Hour * 24,
			expectErr:   false,
		},
		{
			name:        "invalid duration string",
			args:        []string{"--duration", "invalid"},
			expectedVal: 0,
			expectErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := testCmd.ParseFlags(test.args)

			assert.Equal(t, test.expectedVal, testDuration.Duration())

			switch test.expectErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}
