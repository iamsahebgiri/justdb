package justdb

import (
	"os"
	"sync"
)

type DataFile struct {
	sync.RWMutex

	id     uint32
	offset uint32
	writer *os.File
	reader *os.File
}

func NewDataFile(dirPath string) (*DataFile, error) {
	// check if the directory exists, if not create it
	writer, err := os.OpenFile(dirPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	reader, err := os.OpenFile(dirPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := writer.Stat()
	if err != nil {
		return nil, err
	}

	return &DataFile{writer: writer, reader: reader, offset: uint32(stat.Size())}, nil
}

func (df *DataFile) Write(key, value []byte) error {
	return nil
}
