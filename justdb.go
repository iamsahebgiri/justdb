package justdb

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type JustDB struct {
	keyDir           map[byte]KeyDirEntry
	activeDataFile   *DataFile
	inactiveDataFile map[string]*DataFile

	mu sync.RWMutex
}

func New(options *Options) (*JustDB, error) {
	// check if the directory exists, if not create it
	if _, err := os.Stat(options.DirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(options.DirPath, 0755); err != nil {
			return nil, err
		}
	}

	// get all the data files (both active and inactive) in the directory
	files, err := filepath.Glob(fmt.Sprintf("%s/*.dat", options.DirPath))
	if err != nil {
		return nil, err
	}

	var activeDataFile *DataFile
	inactiveDataFile := make(map[string]*DataFile)

	if len(files) > 0 {
		lastFile := files[len(files)-1]
		activeDataFile, err = NewDataFile(options.DirPath, lastFile)
		if err != nil {
			return nil, err
		}

		files = files[:len(files)-1]
		for _, file := range files {
			inactiveDataFile[file], err = NewDataFile(options.DirPath, file)
			if err != nil {
				return nil, err
			}
		}
	} else {
		activeDataFile, err = NewDataFile(options.DirPath, fmt.Sprintf("%d", time.Now().Unix()))
		if err != nil {
			return nil, err
		}
	}

	keyDir := make(map[byte]KeyDirEntry)
	hintsPath := filepath.Join(options.DirPath, "hints.dat")

	// if hints file exists, read it and fill the keyDir
	if _, err := os.Stat(hintsPath); err == nil {
		file, err := os.Open(hintsPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		gob.NewDecoder(file).Decode(&keyDir)
	}

	db := &JustDB{
		keyDir:           keyDir,
		activeDataFile:   activeDataFile,
		inactiveDataFile: inactiveDataFile,
	}

	go db.rotateActiveDatafilePeriodically(options)

	return db, nil
}

func (db *JustDB) Put(key, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	entry := &Entry{
		Timestamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Key:       key,
		Value:     value,
	}
	offset, err := db.activeDataFile.Write(entry)

	db.keyDir[key[0]] = KeyDirEntry{
		FileId:        db.activeDataFile.id,
		ValueSize:     entry.ValueSize,
		ValuePosition: offset,
	}

	return err
}

func (db *JustDB) Get(key []byte) ([]byte, error) {

	entry, ok := db.keyDir[key[0]]
	if !ok {
		return nil, ErrNoKey
	}

	if db.activeDataFile.id == entry.FileId {
		e, err := db.activeDataFile.Read(entry.ValuePosition)
		if err != nil {
			return nil, err
		}

		return e.Value, nil
	}

	file, ok := db.inactiveDataFile[fmt.Sprintf("%d.dat", entry.FileId)]
	if !ok {
		return nil, ErrNoKey
	}

	e, err := file.Read(entry.ValuePosition)
	if err != nil {
		return nil, err
	}

	return e.Value, nil
}

func (db *JustDB) Delete(key []byte) error {
	return nil
}

func (db *JustDB) Keys() []byte {
	var keys []byte
	for k := range db.keyDir {
		keys = append(keys, k)
	}
	return keys
}

func (db *JustDB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.activeDataFile.file.Close()

	// TODO: write keyDir to hint file

	return nil
}
