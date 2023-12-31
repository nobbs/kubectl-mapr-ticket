package util

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
)

func HumanDurationComparedToNow(t time.Time) string {
	difference := time.Since(t)

	if difference < 0 {
		return HumanDuration(-difference)
	}

	return HumanDuration(difference)
}

func ShortHumanDurationComparedToNow(t time.Time) string {
	difference := time.Since(t)

	if difference < 0 {
		return ShortHumanDuration(-difference)
	}

	return ShortHumanDuration(difference)
}

func HumanDurationUntilNow(t time.Time) string {
	difference := time.Since(t)

	return HumanDuration(difference)
}

func ShortHumanDurationUntilNow(t time.Time) string {
	difference := time.Since(t)

	return ShortHumanDuration(difference)
}

func HumanDuration(time time.Duration) string {
	return duration.HumanDuration(time)
}

func ShortHumanDuration(time time.Duration) string {
	return duration.ShortHumanDuration(time)
}
