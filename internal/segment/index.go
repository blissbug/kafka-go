package segment

import (
	"fmt"

	"github.com/yourusername/go-kafka/internal/index"
	"github.com/yourusername/go-kafka/internal/storage"
)

type Segment struct {
	BaseOffset int64
	NextOffset int64

	LogFile   *storage.Log
	IndexFile *index.Index

	maxSize int64
}

func (s *Segment) Open(baseOffset int64, maxSize int64) error {
	s.BaseOffset = baseOffset
	s.NextOffset = baseOffset
	s.maxSize = maxSize

	logFile, err := storage.Open(fmt.Sprintf("%d", s.BaseOffset), s.maxSize)

	if err != nil {
		return err
	}

	s.LogFile = logFile

	indexFile, err := index.Open(fmt.Sprintf("%d", s.BaseOffset))

	if err != nil {
		return err
	}

	s.IndexFile = indexFile

	return nil
}

func (s *Segment) Append(data []byte) (index int, err error) {
	logicalOffset := s.NextOffset
	byteOffset, err := s.LogFile.Append(data)
	if err != nil {
		return 0, err
	}
	err = s.IndexFile.Append(int32(s.NextOffset-s.BaseOffset), byteOffset)
	if err != nil {
		return 0, err
	}
	s.NextOffset++
	return int(logicalOffset), err
}

func (s *Segment) Read(offset int64) ([]byte, error) {
	relativeOffset := offset - s.BaseOffset
	index, err := s.IndexFile.Read(relativeOffset)

	if err != nil {
		return nil, err
	}

	data, err := s.LogFile.Read(int64(index))

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Segment) IsFull() bool {
	return s.LogFile.Size() >= s.maxSize
}

func (s *Segment) Close() error {
	err := s.LogFile.Close()
	if err != nil {
		return err
	}
	err = s.IndexFile.Close()
	if err != nil {
		return err
	}
	return nil
}
