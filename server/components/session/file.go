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
	return f.EndTime.Second() - f.StartTime.Second()
}

func (f File) Completed() bool {
	return f.Transferred == f.Size
}

func (f File) Bitrate() int {
	if f.Duration() == 0 {
		return ^0 >> 1
	}
	return f.Transferred / f.Duration()
}
