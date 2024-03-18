package justdb

import (
	"encoding/gob"
	"io"
	"os"
	"sync"
	"time"
)

type JustDB struct {
	keyDir         map[byte]KeyDirEntry
	activeDataFile *os.File

	mu sync.RWMutex
}

type KeyDirEntry struct {
	FileId        uint32
	ValueSize     uint32
	ValuePosition uint32
	Timestamp     uint32
}

func New(options *Options) (*JustDB, error) {
	// check if the directory exists, if not create it
	file, err := os.OpenFile(options.DirPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &JustDB{activeDataFile: file, keyDir: make(map[byte]KeyDirEntry)}, nil
}

func (db *JustDB) Put(key, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, err := db.activeDataFile.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	entry := &Entry{
		Timestamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Key:       key,
		Value:     value,
	}
	entry.SetChecksum()

	offset, err := db.activeDataFile.Seek(0, io.SeekEnd)

	if err != nil {
		return err
	}

	db.keyDir[key[0]] = KeyDirEntry{
		FileId:        0,
		ValueSize:     entry.ValueSize,
		ValuePosition: uint32(offset),
	}

	return gob.NewEncoder(db.activeDataFile).Encode(entry)
}

func (db *JustDB) Get(key []byte) ([]byte, error) {

	if entry, ok := db.keyDir[key[0]]; ok {
		if _, err := db.activeDataFile.Seek(int64(entry.ValuePosition), io.SeekStart); err != nil {
			return nil, err
		}

		var e Entry
		if err := gob.NewDecoder(db.activeDataFile).Decode(&e); err != nil {
			return nil, err
		}

		if !e.VerifyChecksum() {
			return nil, ErrChecksum
		}

		return e.Value, nil

	}

	return nil, ErrNoKey
}

func (db *JustDB) Delete(key []byte) error {
	return nil
}

func (db *JustDB) Keys(key, value []byte) error {

	return nil
}

func (db *JustDB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.activeDataFile.Close()

	// TODO: write keyDir to hint file

	return nil
}
