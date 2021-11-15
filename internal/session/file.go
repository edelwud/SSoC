package session

import (
	"os"
	"time"
)

// File responsible for filesystem io
type File struct {
	*os.File
	Filename    string
	Transferred int64
	Size        int64
	StartTime   time.Time
	EndTime     time.Time
}

// Close closes file handle, sets EndTime as current
func (f *File) Close() error {
	f.EndTime = time.Now()

	err := f.File.Close()
	if err != nil {
		return err
	}

	return nil
}

// Duration returns delta between EndTime and StartTime in nanoseconds
func (f File) Duration() int {
	return f.EndTime.Nanosecond() - f.StartTime.Nanosecond()
}

// Completed returns true if Transferred equals Size
func (f File) Completed() bool {
	return f.Transferred == f.Size
}

// Bitrate returns bitrate as Transferred / Duration (in MB/s)
func (f File) Bitrate() float64 {
	if f.Duration() == 0 {
		return ^0 >> 1
	}
	return float64(f.Transferred) / float64(f.Duration()) * 1000000000 / 1024 / 1024
}

// Read implements io.Reader interface, calculates Transferred
func (f *File) Read(p []byte) (int, error) {
	n, err := f.File.Read(p)
	if err == nil {
		f.Transferred += int64(n)
	}

	return n, err
}

// Write implements io.Writer interface, calculates Transferred
func (f *File) Write(p []byte) (int, error) {
	n, err := f.File.Write(p)
	if err == nil {
		f.Transferred += int64(n)
	}

	return n, err
}

// CreateFile creates file (flags: os.O_CREATE|os.O_RDWR, permissions: 0777), receives stats
func CreateFile(filename string, filepath string) (*File, error) {
	opened, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	stat, err := opened.Stat()
	if err != nil {
		return nil, err
	}

	return &File{
		File:        opened,
		Filename:    filename,
		Transferred: 0,
		Size:        stat.Size(),
		StartTime:   time.Now(),
		EndTime:     time.Time{},
	}, nil
}
