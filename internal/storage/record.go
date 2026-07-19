package storage

import (
	"errors"

	"github.com/yourusername/go-kafka/internal/utils"
)

const HEADER_SIZE = 4

// SerializeRecord whatever msg we receive, we need to serialize and append the length of it to the beginning
func SerializeRecord(msg []byte) []byte {
	length := len(msg)
	v := utils.ToByteArr32(length, HEADER_SIZE)
	return append(v, msg...)
}

func DeserializeRecord(data []byte) ([]byte, error) {
	if len(data) < HEADER_SIZE {
		return nil, errors.New("invalid record")
	}
	length := data[0:HEADER_SIZE]
	v := utils.FromByteArr(length, HEADER_SIZE)
	if len(data) < v+HEADER_SIZE {
		return nil, errors.New("corrupted record")
	}
	return data[HEADER_SIZE : v+HEADER_SIZE], nil
}
