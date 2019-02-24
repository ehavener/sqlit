/* UNR CS 457 | SPRING 2019 | emerson@nevada.unr.edu */

// Package io is a library of ideally-atomic operations
// used on persisted databases and tables.
package io

// TODO: rename to utility

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"time"
)

const (
	path = "tmp/"
)

var database string

// CheckIfDatabaseExists checks if a database directory exists
func CheckIfDatabaseExists(name string) bool {
	_, err := os.Stat(path + name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDatabase creates a database directory
func CreateDatabase(name string) error {
	err := os.Mkdir(path+name, os.ModePerm)
	return err
}

// CreateDatabaseMeta places a dotfile inside the database with brief details
func CreateDatabaseMeta(name string) {
	f, err := os.Create(path + name + "/" + ".meta")
	check(err)
	defer f.Close()

	createdAt := time.Now().Format(time.RFC850)
	owner, err := user.Current()
	check(err)

	f.WriteString("owner" + " | " + owner.Name + "\n")
	f.WriteString("createdAt" + " | " + createdAt)
}

// UseDatabase stores the name of the database in memory
func UseDatabase(name string) {
	database = name
	return
}

// DeleteDatabase removes a database directory
func DeleteDatabase(name string) error {
	err := os.Remove(path + database + "/" + name)
	return err
}

// CheckIfAnyDatabaseIsInUse checks if a database name is stored in mem
func CheckIfAnyDatabaseIsInUse() bool {
	if database == "" {
		return false
	}
	return true
}

// CheckIfTableExists ...
func CheckIfTableExists(name string) bool {
	_, err := os.Stat(path + database + "/" + name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateTable ...
func CreateTable(name string, columns []string, constraints []string) {
	f, err := os.Create(path + database + "/" + name)
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
	err := os.Remove(path + database + "/" + name)
	check(err)
}

// AlterTable ... TODO: fix
func AlterTable(name string, method string, column string, constraint string) string {
	fileContents, err := ioutil.ReadFile(path + database + "/" + name)
	check(err)

	f, err := os.Create(path + database + "/" + name)
	check(err)
	defer f.Close()

	if method == "ADD" {
		f.WriteString(string(fileContents) + " | " + column + " " + constraint)
	}

	return SelectAll(name)
}

// SelectAll ...
func SelectAll(name string) string {
	fileContents, err := ioutil.ReadFile(path + database + "/" + name)
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

//
//			Helper functions
//

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
