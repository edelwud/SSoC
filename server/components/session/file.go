package session

import (
	"time"
)

type File struct {
	Filename    string
	Transferred int
	Size        int
	StartTime   time.Time
	EndTime     time.Time
}

func (f File) Duration() int {
	return f.EndTime.Nanosecond() - f.StartTime.Nanosecond()
}

func (f File) Completed() bool {
	return f.Transferred == f.Size
}

func (f File) Bitrate() float64 {
	if f.Duration() == 0 {
		return ^0 >> 1
	}
	return float64(f.Transferred) / float64(f.Duration()) * 1000000000 / 1024 / 1024
}
