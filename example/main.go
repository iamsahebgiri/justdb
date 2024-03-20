package main

import (
	"fmt"
	"os"

	"github.com/iamsahebgiri/justdb"
)

func main() {
	fmt.Printf("Hello, world\n")
	db, err := justdb.New(&justdb.Options{
		DirPath: "data",
	})

	if err != nil {
		panic(err)
	}

	db.Put([]byte("name"), []byte("John"))
	db.Put([]byte("age"), []byte("30"))

	if val, err := db.Get([]byte("age")); err == nil {
		fmt.Println(string(val))
	}

	fmt.Println(db.Keys())

	db.Close()

	// ticker := time.NewTicker(1 * time.Second).C

	// for range ticker {
	// 	os.OpenFile(fmt.Sprintf("%d.dat", time.Now().Unix()), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	// }

	file, err := os.Stat("171096422300.dat")
	fmt.Println(err)
	fmt.Println(file)

	// files, _ := filepath.Glob("./*.dat")
	// var allFiles []uint32

	// for _, file := range files {
	// 	timestamp, _ := strconv.ParseUint(strings.TrimRight(file, filepath.Ext(file)), 10, 32)
	// 	allFiles = append(allFiles, uint32(timestamp))
	// }

	// fmt.Println(allFiles)

	// fmt.Println(sort.SliceIsSorted(allFiles, func(i, j int) bool {
	// 	return allFiles[i] < allFiles[j]
	// }))
}
