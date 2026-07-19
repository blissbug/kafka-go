// Package storage section to append data and read data
package storage

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/go-kafka/internal/utils"
)

type Log struct {
	maxSize int64
	file    *os.File
}

func Open(f string, maxSize int64) (*Log, error) {
	err := os.MkdirAll("data", 0777)
	if err != nil {
		return nil, fmt.Errorf("could not create directory: %w", err)
	}
	fileName := f + ".log"
	file, err := os.OpenFile("data/"+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	l := &Log{
		maxSize: maxSize,
		file:    file,
	}

	return l, nil
}

func (l *Log) Append(data []byte) (int64, error) {
	//we receive a byte array, we need to append it to the file
	//if active segment is full, we need to create a new segment
	//normally appending for now
	//return the ending byte position
	currentPosition, err := l.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	data = SerializeRecord(data)
	_, err = l.file.Write(data)

	if err != nil {
		return 0, err
	}

	return currentPosition, nil
}

func (l *Log) Read(position int64) ([]byte, error) {
	_, err := l.file.Seek(position, io.SeekStart)
	if err != nil {
		return nil, err
	}

	byteLength := make([]byte, HEADER_SIZE)
	_, err = io.ReadFull(l.file, byteLength)

	if err != nil {
		return nil, err
	}

	length := utils.FromByteArr32(byteLength, HEADER_SIZE)
	data := make([]byte, length)
	_, err = io.ReadFull(l.file, data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *Log) Size() int64 {
	fileInfo, err := l.file.Stat()
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

func (l *Log) Close() error {
	err := l.file.Close()
	if err != nil {
		return err
	}
	return nil
}
