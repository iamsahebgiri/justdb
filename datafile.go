package justdb

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DataFile struct {
	sync.RWMutex

	id     uint32
	offset uint32
	file   *os.File
}

func NewDataFile(dirPath, fileName string) (*DataFile, error) {
	var id uint32
	name := fmt.Sprintf("%s/%s.dat", dirPath, fileName)

	if file, err := os.Stat(dirPath); os.IsExist(err) {
		timestamp, _ := strconv.ParseUint(strings.TrimRight(file.Name(), filepath.Ext(file.Name())), 10, 32)
		id = uint32(timestamp)
	}

	// check if the directory exists, if not create it also open in append mode rw mode
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if id == 0 {
		id = uint32(time.Now().Unix())
	}

	return &DataFile{file: file, offset: uint32(stat.Size()), id: id}, nil
}

func (df *DataFile) Write(entry *Entry) (uint32, error) {
	df.Lock()
	defer df.Unlock()

	entry.SetChecksum()

	offset, err := df.file.Seek(0, io.SeekEnd)

	if err != nil {
		return uint32(offset), err
	}

	err = gob.NewEncoder(df.file).Encode(entry)
	if err != nil {
		return uint32(offset), err
	}
	offset, err = df.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return uint32(offset), err
	}

	df.offset = uint32(offset)

	return uint32(offset), nil
}

func (df *DataFile) Read(offset uint32) (Entry, error) {
	df.RLock()
	defer df.RUnlock()

	_, err := df.file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return Entry{}, err
	}

	var entry Entry
	err = gob.NewDecoder(df.file).Decode(&entry)
	if err != nil {
		return Entry{}, err
	}

	if !entry.VerifyChecksum() {
		return Entry{}, ErrChecksumMismatch
	}

	return entry, nil
}
