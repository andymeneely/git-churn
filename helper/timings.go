package helper

import (
	. "github.com/andymeneely/git-churn/print"
	"time"
)

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	PrintInBlue("%v: %v\n", msg, time.Since(start))
}
