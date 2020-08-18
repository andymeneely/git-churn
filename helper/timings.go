package helper

import (
	. "github.com/andymeneely/git-churn/print"
	"time"
)

// Track returns the msg and the current time to start the tracking.
func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

// Duration prints the time since start with the msg.
func Duration(msg string, start time.Time) {
	PrintInBlue("%v: %v\n", msg, time.Since(start))
}
