// Package storage section to append data and read data
package storage

type Log struct{}

func (*Log) append(data []byte) error {
	return nil
}

func (*Log) read(position int64) ([]byte, error) {
	return nil, nil
}
