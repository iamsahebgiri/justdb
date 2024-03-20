package justdb

import (
	"fmt"
	"os"
	"time"
)

func (db *JustDB) rotateActiveDatafilePeriodically(options *Options) {
	ticker := time.NewTicker(5 * time.Minute).C

	for range ticker {
		if err := db.rotateActiveDatafile(options); err != nil {
			fmt.Println(err)
		}
	}
}

func (db *JustDB) rotateActiveDatafile(options *Options) error {

	db.mu.Lock()
	defer db.mu.Unlock()

	file, err := db.activeDataFile.file.Stat()
	if err != nil {
		return err
	}

	if file.Size() < options.ActiveDataFileMaxSize {
		return nil
	}

	// close the current active data file
	if err := db.activeDataFile.file.Close(); err != nil {
		return err
	}

	// move the current active data file to inactive data file
	if err := os.Rename(db.activeDataFile.file.Name(), fmt.Sprintf("%s/%d.dat", options.DirPath, time.Now().Unix())); err != nil {
		return err
	}

	db.inactiveDataFile[db.activeDataFile.file.Name()] = db.activeDataFile

	// create a new active data file
	activeDataFile, err := NewDataFile(options.DirPath, "active")
	if err != nil {
		return err
	}

	db.activeDataFile = activeDataFile

	return nil

}
