package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Timestamp returns unix now time
func Timestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Second))
}

// Check check result for errors
func Check(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}
