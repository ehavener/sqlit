package io

// A separate B-tree is used for each table and index in the database.
// All B-trees are stored in the same disk file.

import (
	// "bufio"
	// "fmt"
	// "io/ioutil"
	"fmt"
	"os"
)

const (
	path = "tmp/"
	ext  = ".db"
)

// Execute ...
func Execute() {
	// CreateDatabase("db_1")
	// deleteDatabase("hello_io")
}

// CheckIfDatabaseExists ...
func CheckIfDatabaseExists(name string) bool {
	_, err := os.Stat(path + name + ext)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDatabase ...
func CreateDatabase(name string) {
	f, err := os.Create(path + name + ext)
	check(err)
	defer f.Close()
}

// DeleteDatabase ...
func DeleteDatabase(name string) {
	err := os.Remove(path + name + ext)
	check(err)
}

func open() {}

func commit() {}

func close() {}

func createTable() {}

func lockTable() {}

func unlockTable() {}

func dropTable() {}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		// panic(e)
	}
}
