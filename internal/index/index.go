package index

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/go-kafka/internal/utils"
)

type Index struct {
	file *os.File
}

const ENTRY_SIZE = 12

func Open(f string) (*Index, error) {
	err := os.MkdirAll("data", 0777)
	if err != nil {
		return nil, fmt.Errorf("could not create directory: %w", err)
	}
	fileName := f + ".index"
	file, err := os.OpenFile("data/"+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	l := &Index{
		file: file,
	}

	return l, nil
}

// Append add a new entry to the index - offset is a logical entry and bytePosition is the position in the file
func (index *Index) Append(offset int32, bytePosition int64) error {
	//serialize and append a new entry to the index
	_, err := index.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	serializedOffset := utils.ToByteArr32(offset, 4)
	serializedBytePosition := utils.ToByteArr64(bytePosition, 8)
	data := append(serializedOffset, serializedBytePosition...)
	_, err = index.file.Write(data)

	if err != nil {
		return err
	}

	return nil
}

// Read reads an entry from the index, it gets the logical offset and returns the position in the file
func (index *Index) Read(relativeOffset int64) (uint64, error) {
	entryByteOffset := ENTRY_SIZE * relativeOffset
	_, err := index.file.Seek(entryByteOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	offsetLength := make([]byte, 4)
	_, err = io.ReadFull(index.file, offsetLength)

	if err != nil {
		return 0, err
	}

	bytePosition := make([]byte, 8)
	_, err = io.ReadFull(index.file, bytePosition)

	if err != nil {
		return 0, err
	}

	bytePositionInt := utils.FromByteArr64(bytePosition, 8)
	return bytePositionInt, nil
}

func (index *Index) EntryCount() (int64, error) {
	fileSize, err := index.file.Stat()
	if err != nil {
		return 0, err
	}
	return fileSize.Size() / ENTRY_SIZE, nil
}

func (index *Index) Close() error {
	err := index.file.Close()
	if err != nil {
		return err
	}
	return nil
}
