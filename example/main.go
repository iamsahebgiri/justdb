package main

import (
	"fmt"

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

	db.Close()
}
