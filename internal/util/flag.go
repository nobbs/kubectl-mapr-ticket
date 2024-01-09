package util

import (
	"time"

	"github.com/xhit/go-str2duration/v2"
)

// DurationValue is a wrapper around time.Duration that implements the
// pflag.Value interface.
//
// Reason: the default time.Duration implementation of pflag.Value only supports
// units up to hours. That's not really enough for us as we most likely want to
// check for tickets that expire in a few days or even weeks. This wrapper uses
// the go-str2duration library to additionally support days and weeks.
type DurationValue time.Duration

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

// Duration returns the underlying time.Duration value
func (d *DurationValue) Duration() time.Duration {
	return time.Duration(*d)
}
