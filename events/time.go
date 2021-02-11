package events

import "time"

const (
	defaultTimeFormat = time.RFC3339
)

func TimeNow() string {
	return time.Now().Format(defaultTimeFormat)
}
