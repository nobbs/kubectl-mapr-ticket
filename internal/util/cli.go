package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

const (
	indend = `  `
)

// DurationValue is a wrapper around time.Duration that implements the
// pflag.Value interface.
//
// Reason: the default time.Duration implementation of pflag.Value only supports
// units up to hours. That's not really enough for us as we most likely want to
// check for tickets that expire in a few days or even weeks. This wrapper uses
// the go-str2duration library to additionally support days and weeks.
type DurationValue time.Duration

// NewDurationValue returns a new DurationValue with the specified time.Duration
func NewDurationValue(val time.Duration) *DurationValue {
	return (*DurationValue)(&val)
}

// Set implements the pflag.Value interface for DurationValue
func (d *DurationValue) Set(s string) error {
	v, err := str2duration.ParseDuration(s)
	*d = DurationValue(v)
	return err
}

// Type implements the pflag.Value interface for DurationValue
func (d *DurationValue) Type() string {
	return "duration"
}

// String implements the pflag.Value interface for DurationValue
func (d *DurationValue) String() string {
	return time.Duration(*d).String()
}

// Cast returns the underlying time.Duration value
func (d *DurationValue) Cast() time.Duration {
	return time.Duration(*d)
}

type stringNormalizer struct {
	string
}

func CliLongDesc(desc string, args ...any) string {
	if desc == "" || len(desc) == 0 {
		return ""
	}

	desc = fmt.Sprintf(desc, args...)

	return stringNormalizer{desc}.trim().string
}

func CliExample(example string, args ...any) string {
	if example == "" || len(example) == 0 {
		return ""
	}

	example = fmt.Sprintf(example, args...)

	return stringNormalizer{example}.trim().indent().string
}

func (s stringNormalizer) trim() stringNormalizer {
	s.string = strings.TrimSpace(s.string)

	lines := strings.Split(s.string, "\n")

	for i, line := range strings.Split(s.string, "\n") {
		lines[i] = strings.TrimSpace(line)
	}

	s.string = strings.Join(lines, "\n")
	return s
}

func (s stringNormalizer) indent() stringNormalizer {
	lines := strings.Split(s.string, "\n")

	for i, line := range lines {
		lines[i] = indend + line
	}

	s.string = strings.Join(lines, "\n")
	return s
}

func StringSliceToFlagOptions(slice []string) string {
	return strings.Join(slice, ", ")
}
