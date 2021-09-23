package session

import (
	"os"
	"time"
)

type File struct {
	*os.File
	Filename    string
	Transferred int64
	Size        int64
	StartTime   time.Time
	EndTime     time.Time
}

func (f *File) Close() error {
	f.EndTime = time.Now()

	err := f.File.Close()
	if err != nil {
		return err
	}

	return nil
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

func (f *File) Read(p []byte) (int, error) {
	n, err := f.File.Read(p)
	f.Transferred += int64(n)

	return n, err
}

func (f *File) Write(p []byte) (int, error) {
	n, err := f.File.Write(p)
	f.Transferred += int64(n)

	return n, err
}

func CreateFile(filename string, filepath string) (*File, error) {
	opened, err := os.Open(filepath)
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
