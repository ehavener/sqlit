package io

// A separate B-tree is used for each table and index in the database.
// All B-trees are stored in the same disk file.

import (
	// "bufio"
	// "fmt"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	path = "tmp/"
	ext  = ".db"
	ext2 = ".lit"
)

var database string

// CheckIfDatabaseExists ...
func CheckIfDatabaseExists(name string) bool {
	_, err := os.Stat(path + name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDatabase ...
func CreateDatabase(name string) {
	os.Mkdir(path+name, os.ModePerm)
}

// UseDatabase ...
func UseDatabase(name string) {
	database = name
	return
}

// CheckIfAnyDatabaseIsInUse ...
func CheckIfAnyDatabaseIsInUse() bool {
	if database == "" {
		return false
	}
	return true
}

// CheckIfTableExists ...
func CheckIfTableExists(name string) bool {
	_, err := os.Stat(path + database + "/" + name + ext2)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateTable ...
func CreateTable(name string, columns []string, constraints []string) {
	f, err := os.Create(path + database + "/" + name + ext2)
	check(err)
	defer f.Close()

	for i := 0; i < len(columns); i++ {
		if i > 0 {
			f.WriteString(" | ")
		}
		f.WriteString(columns[i] + " " + constraints[i])
	}
}

// DropTable ...
func DropTable(name string) {
	err := os.Remove(path + database + "/" + name + ext2)
	check(err)
}

// AlterTable ... TODO: fic ASAP
func AlterTable(name string, method string, column string, constraint string) string {
	fileContents, err := ioutil.ReadFile(path + database + "/" + name + ext2)
	check(err)

	f, err := os.Create(path + database + "/" + name + ext2)
	check(err)
	defer f.Close()

	if method == "ADD" {
		f.WriteString(string(fileContents) + " | " + column + " " + constraint)
	}

	return SelectAll(name)
}

// DeleteDatabase ...
func DeleteDatabase(name string) {
	err := os.Remove(path + database + "/" + name)
	check(err)
}

// SelectAll ...
func SelectAll(name string) string {
	fileContents, err := ioutil.ReadFile(path + database + "/" + name + ext2)
	check(err)
	return string(fileContents)
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
		panic(e)
	}
}
