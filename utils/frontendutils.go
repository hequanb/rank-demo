package utils

import (
	"time"
)

func BackendToFrontendTime(t time.Time) time.Time {
	return t.Add(8 * time.Hour)
}


