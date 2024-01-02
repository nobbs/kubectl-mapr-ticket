package util

import (
	"math"
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
)

// HumanDurationComparedToNow returns a human readable string representing the
// duration between the given time and now. Always returns a positive duration.
func HumanDurationComparedToNow(t time.Time) string {
	d := time.Since(t)

	// edge case: if the difference is the maximum absolute value of an int64,
	// then we hit the max possible duration in any direction. In this case,
	// we'll just return "inf" as the duration.
	if d.Abs() == math.MaxInt64 {
		return "inf"
	}

	// if the difference is negative, then we'll just return the human readable string
	// for the positive difference, but with a negative sign in front of it.
	if d < 0 {
		return HumanDuration(-d)
	}

	return HumanDuration(d)
}

// ShortHumanDurationComparedToNow returns a short human readable string
// representing the duration between the given time and now. Always returns a
// positive duration.
func ShortHumanDurationComparedToNow(t time.Time) string {
	d := time.Since(t)

	// edge case: if the difference is the maximum absolute value of an int64,
	// then we hit the max possible duration in any direction. In this case,
	// we'll just return "inf" as the duration.
	if d.Abs() == math.MaxInt64 {
		return "inf"
	}

	// if the difference is negative, then we'll just return the human readable string
	// for the positive difference, but with a negative sign in front of it.
	if d < 0 {
		return ShortHumanDuration(-d)
	}

	return ShortHumanDuration(d)
}

// HumanDurationUntilNow returns a human readable string representing the
// duration between the given time and now.
func HumanDurationUntilNow(t time.Time) string {
	d := time.Since(t)

	// edge case: if the difference is the maximum absolute value of an int64,
	// then we hit the max possible duration in any direction. In this case,
	// we'll just return "inf" as the duration.
	if d.Abs() == math.MaxInt64 {
		return "inf"
	}

	return HumanDuration(d)
}

// ShortHumanDurationUntilNow returns a short human readable string
// representing the duration between the given time and now.
func ShortHumanDurationUntilNow(t time.Time) string {
	d := time.Since(t)

	// edge case: if the difference is the maximum absolute value of an int64,
	// then we hit the max possible duration in any direction. In this case,
	// we'll just return "inf" as the duration.
	if d.Abs() == math.MaxInt64 {
		return "inf"
	}

	return ShortHumanDuration(d)
}

// HumanDuration returns a human readable string representing the given
// duration.
func HumanDuration(d time.Duration) string {
	return duration.HumanDuration(d)
}

// ShortHumanDuration returns a short human readable string representing the
// given duration.
func ShortHumanDuration(d time.Duration) string {
	return duration.ShortHumanDuration(d)
}
