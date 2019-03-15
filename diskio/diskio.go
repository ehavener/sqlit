/* UNR CS 457 | SPRING 2019 | emerson@nevada.unr.edu */

// Package diskio is a library of ideally-atomic operations
// used on persisted databases and tables.
package diskio

// TODO: rename to utility

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
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

// InsertRecord ...
func InsertRecord(name string, values []string) error {
	f, err := os.OpenFile(path+database+"/"+name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer f.Close()

	var writeBuffer string

	writeBuffer += "\n"

	for i := 0; i < len(values); i++ {
		if i > 0 {
			writeBuffer += " | "
		}
		writeBuffer += values[i] + " "
	}

	f2, err2 := f.Write([]byte(writeBuffer))
	if f2 == 0 {
		return err
	}
	check(err2)

	// n, err := f.WriteString(f, "\n")
	// check(err)

	return err
}

// UpdateRecord updates a record
// get line 0, count how many cols it takes to find whereCol, store as colOffset
// iterate thru lines 1 - n until row[colOffset] == whereValue
// store that line (record) -- like it's selected
// look in selected record for toCol, replace record[toCol] with toValue
func UpdateRecord(table string, whereCol string, whereValue string, toCol string, toValue string) (string) {
	// fmt.Println("table " + table)
	// fmt.Println("whereCol " + whereCol)
	// fmt.Println("whereValue " + whereValue)
	// fmt.Println("toCol " + toCol)
	// fmt.Println("toValue " + toValue)
	// fmt.Println("")

	input, err := ioutil.ReadFile(path + database + "/" + table)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	recordAmount := getAmountOfRecordsInTable(table)
	// fmt.Println("recordAmount ", recordAmount)

	f, err := os.Open(path + database + "/" + table)
	check(err)
	defer f.Close()

	reader := bufio.NewReader(f)
	tableMetaLine, _ := reader.ReadString('\n')
	colDefs := strings.Split(tableMetaLine, "|")

	// break the constraint off the colName
	colNames := make([]string, 0, len(colDefs))
	for i := range colDefs {
		colNames = append(colNames, strings.Fields(colDefs[i])[0])
	}

	// var whereColOffset int
	var toColOffset int
	// var selectedRowOffset int

	// find col offset of whereCol
	// for i := range colNames {
	// 	fmt.Println(colNames[i], whereCol)
	// 	if strings.EqualFold(colNames[i], whereCol) {
	// 		whereColOffset = i
	// 	}
	// }

	// find col offset of toCol
	for i := range colNames {
		if colNames[i] == whereCol {
			toColOffset = i
		}
	}

	// find row offset of whereValue
	for i := 0; i < recordAmount; i++ {
		record, _ := reader.ReadString('\n')
		values := strings.Split(record, "|")
		for j := range values {
			if strings.EqualFold(strings.TrimSpace(values[j]), strings.TrimSpace(whereValue)) {
				// selectedRowOffset = i

				// find toColOffset within table[selectedRowOffset] and assign it toValue
				newValues := values
				newValues[toColOffset] = toValue

				newRecord := strings.Join(newValues, "|")

				output := bytes.Replace(input, []byte(record), []byte(newRecord), -1)

				// fmt.Println(output)
				if err = ioutil.WriteFile(path+database+"/"+table, output, 0666); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
	}

	return "1"
}

// DeleteRecord deletes a record
func DeleteRecord(table string, whereCol string, whereValue string) {

}

func getAmountOfRecordsInTable(table string) int {
	file, _ := os.Open(path + database + "/" + table)
	fileScanner := bufio.NewScanner(file)
	recordCount := -1
	for fileScanner.Scan() {
		recordCount++
	}
	return recordCount
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
