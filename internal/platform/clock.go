// Package platform defines replaceable host capabilities for SeekMoon.
package platform

import "time"

// Clock abstracts current time for deterministic source evidence.
type Clock interface {
	Now() time.Time
}

// SystemClock reads wall-clock time.
type SystemClock struct{}

// Now returns the current wall-clock time.
func (SystemClock) Now() time.Time {
	return time.Now()
}

// FixedClock returns a configured time for tests.
type FixedClock struct {
	Time time.Time
}

// Now returns the configured fixed time.
func (c FixedClock) Now() time.Time {
	return c.Time
}
