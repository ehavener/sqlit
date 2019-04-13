/* UNR CS 457 | SPRING 2019 | emerson@nevada.unr.edu */

// Package diskio is a library of ideally-atomic operations
// used on persisted databases and tables.
package diskio

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const (
	path = "tmp/"
)

var database string

// ColumnDef ...
type ColumnDef struct {
	ColumnName string
	TypeName   string
}

// Set ...
type Set struct {
	Name       string
	ColumnDefs []ColumnDef
	Records     [][]string
}

// SerializeSet ...
func SerializeSet(set Set) string {
	return SerializeColumnDefs(set.ColumnDefs) + "\n" + SerializeRecords(set.Records)
}

// SerializeColumnDefs ...
func SerializeColumnDefs(columnDefs []ColumnDef) string {
	var serializedColumnDef string

	for index, columnDef := range columnDefs {
		serializedColumnDef += columnDef.ColumnName + " " + columnDef.TypeName
		if (index + 1) < len(columnDefs) {
			serializedColumnDef += "|"
		}
	}

	return serializedColumnDef
}

// ConstructColumnDefs ...
func ConstructColumnDefs(columnDefsLine string) []ColumnDef {
	columnDefsPairs := strings.Split(columnDefsLine, "|")

	var columnDefs []ColumnDef

	for _, columnDefsPair := range columnDefsPairs {
		columnDefsPair := strings.Fields(columnDefsPair)
		columnDefs = append(columnDefs, ColumnDef{ColumnName: columnDefsPair[0], TypeName: columnDefsPair[1]})
	}

	return columnDefs
}

// SerializeRecords ...
func SerializeRecords(records [][]string) string {
	var recordsSerialized string

	for _, row := range records {
		for colIndex, col := range row {
			recordsSerialized += strings.Replace(col, "\n", "", -1)
			if (colIndex +  1) < len(row) {
				recordsSerialized += "|"
			} else {
				recordsSerialized += "\n"
			}
		}
	}

	return recordsSerialized
}

// ConstructRecords ...
func ConstructRecords(reader *bufio.Reader, recordAmount int) [][]string {

	records := make([][]string, recordAmount)

	for i := range records {
		recordsLine, _ := reader.ReadString('\n')
		recordsPair := strings.Split(recordsLine, "|")
		records[i] = make([]string, 2)
		records[i][0] = recordsPair[0]
		records[i][1] = recordsPair[1]
	}

	return records
}

// SelectSet ...
func SelectSet(tableName string) Set {
	// open the table file contents
	f, err := os.Open(path + database + "/" + tableName)
	check(err)
	defer f.Close()

	// read in the metadata line, parse the table's col names
	reader := bufio.NewReader(f)
	columnDefsLine, _ := reader.ReadString('\n')

	recordAmount := getAmountOfRecordsInTable(tableName)

	columnDefs := ConstructColumnDefs(columnDefsLine)
	// columnDefsSerialized := SerializeColumnDefs(columnDefs)
	// fmt.Println("columnDefsSerialized: " + columnDefsSerialized)

	records := ConstructRecords(reader, recordAmount)
	// recordsSerialized := SerializeRecords(records)
	// fmt.Println("recordsSerialized: " + recordsSerialized)

	set := Set{Name: tableName, ColumnDefs: columnDefs, Records: records}
	// setSerialized := SerializeSet(set)
	// fmt.Println("setSerialized: " + setSerialized)

	return set
}

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

	f.WriteString("owner" + "|" + owner.Name + "\n")
	f.WriteString("createdAt" + "|" + createdAt)
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

// CheckIfTableExists does as named
func CheckIfTableExists(name string) bool {
	_, err := os.Stat(path + database + "/" + name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateTable creates a table file and initializes its metadata
func CreateTable(name string, columns []string, constraints []string) {
	f, err := os.Create(path + database + "/" + name)
	check(err)

	defer f.Close()

	for i := 0; i < len(columns); i++ {
		if i > 0 {
			f.WriteString("|")
		}
		f.WriteString(columns[i] + " " + constraints[i])
	}
}

// DropTable deletes a table's file
func DropTable(name string) {
	err := os.Remove(path + database + "/" + name)
	check(err)
}

// AlterTable modifies a table's metadata
func AlterTable(name string, method string, column string, constraint string) string {
	fileContents, err := ioutil.ReadFile(path + database + "/" + name)
	check(err)

	f, err := os.Create(path + database + "/" + name)
	check(err)
	defer f.Close()

	if method == "ADD" {
		f.WriteString(string(fileContents) + "|" + column + " " + constraint)
	}

	return SelectAll(name)
}

// SelectAll selects all records in a table
func SelectAll(name string) string {
	fileContents, err := ioutil.ReadFile(path + database + "/" + name)
	check(err)
	return string(fileContents)
}

// SelectWhere selects a subset of records in a table, provided a clause
func SelectWhere(table string, colNames []string, whereColName string, whereColVal string) string {

	// open the table file contents
	f, err := os.Open(path + database + "/" + table)
	check(err)
	defer f.Close()

	// read in the metadata line, parse the table's col names
	reader := bufio.NewReader(f)
	tableMetaLine, _ := reader.ReadString('\n')
	colDefs := strings.Split(tableMetaLine, "|")

	// find pertinent offsets of clause
	var firstColOffset int
	var secondColOffset int
	var whereColOffset int

	for i := range colDefs {
		if strings.EqualFold(strings.Fields(colDefs[i])[0], whereColName) {
			whereColOffset = i
		}
		if strings.EqualFold(strings.Fields(colDefs[i])[0], colNames[0]) {
			firstColOffset = i
		}
		if strings.EqualFold(strings.Fields(colDefs[i])[0], colNames[1]) {
			secondColOffset = i
		}
	}

	// initialize our selection subset
	var selection string

	// populate selection by iterating through records
	for i := 0; i < getAmountOfRecordsInTable(table); i++ {
		record, _ := reader.ReadString('\n')
		records := strings.Split(record, "|")

		if strings.EqualFold(strings.TrimSpace(records[whereColOffset]), strings.TrimSpace(whereColVal)) == false {
			selection = records[firstColOffset] + "|" + records[secondColOffset]
		}
	}

	// make new table meta based on queried cols
	newMeta := strings.TrimSpace(colDefs[firstColOffset]) + "|" + colDefs[secondColOffset]

	// append selection meta to selection's resultant records
	result := newMeta + selection

	return string(result)
}

// InsertRecord inserts a single record to a table
func InsertRecord(name string, records []string) error {

	// Open the table in append mode
	f, err := os.OpenFile(path+database+"/"+name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer f.Close()

	var writeBuffer string

	writeBuffer += "\n"

	// construct new record from records
	for i := 0; i < len(records); i++ {
		if i > 0 {
			writeBuffer += "|"
		}
		writeBuffer += records[i]
	}

	// write the new record to the end of the table
	f2, err2 := f.Write([]byte(writeBuffer))
	if f2 == 0 {
		return err
	}
	check(err2)

	return err
}

// UpdateRecord updates a record
// get line 0, count how many cols it takes to find whereCol, store as colOffset
// iterate thru lines 1 - n until row[colOffset] == whereValue
// store that line (record) -- like it's selected
// look in selected record for toCol, replace record[toCol] with toValue
func UpdateRecord(table string, whereCol string, whereValue string, toCol string, toValue string) int {

	recordAmount := getAmountOfRecordsInTable(table)
	colNames := getTableColNames(table)
	toColOffset := getIndexOfColName(colNames, toCol)
	recordsModified := 0

	// for each record
	for i := 0; i < recordAmount; i++ {

		// open table file
		f, err := os.Open(path + database + "/" + table)
		check(err)
		defer f.Close()

		// open reader on table file contents
		reader := bufio.NewReader(f)
		for rp := 0; rp <= i; rp++ {
			reader.ReadString('\n')
		}

		// read record and split by colVals
		record, _ := reader.ReadString('\n')
		records := strings.Split(record, "|")

		// for each colVal in record
		for j := range records {

			// if colVal matches whereValue
			if strings.EqualFold(strings.TrimSpace(records[j]), strings.TrimSpace(whereValue)) {

				// replace colVal with toValue
				newRecords := records
				newRecords[toColOffset] = toValue

				// rebuild record
				newRecord := strings.Join(newRecords, "|")

				if strings.Contains(newRecord, "\n") == false {
					newRecord += "\n"
				}

				replaceRecord(table, record, newRecord)

				// save the last updated index
				recordsModified++
			}
		}
	}

	return recordsModified
}

// DeleteRecord deletes a record that matches a clause
func DeleteRecord(table string, whereCol string, whereValue string, comparator string) int {

	recordAmount := getAmountOfRecordsInTable(table)

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

	recordsModified := 0

	// for each record
	for i := 0; i < recordAmount; i++ {

		// open table file
		f, err := os.Open(path + database + "/" + table)
		check(err)
		defer f.Close()

		// open reader on table file contents
		reader := bufio.NewReader(f)
		for rp := 0; rp <= i; rp++ {
			reader.ReadString('\n')
		}

		// read record and split by colVals
		record, _ := reader.ReadString('\n')
		records := strings.Split(record, "|")

		var whereColOffset int
		for i := range colNames {
			if strings.EqualFold(colNames[i], whereCol) {
				whereColOffset = i
			}
		}

		// for each colVal in record
		for j := range records {

			// if colVal matches whereValue
			if strings.EqualFold(strings.TrimSpace(comparator), "GREATER_THAN") {

				if j == whereColOffset {
					curVal := strings.TrimSpace(records[j])
					whereVal := strings.TrimSpace(whereValue)

					// parse existing records into floats to allow comparison
					curValFloat, err := strconv.ParseFloat(curVal, 32)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					whereValFloat, err2 := strconv.ParseFloat(whereVal, 32)
					if err2 != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					if curValFloat > whereValFloat {
						removeRecord(table, record)
						// save the last updated index
						recordsModified++
					}
				}

			} else if strings.EqualFold(strings.TrimSpace(comparator), "EQUALS") {

				if strings.EqualFold(strings.TrimSpace(records[j]), strings.TrimSpace(whereValue)) {
					removeRecord(table, record)
					recordsModified++
				}
			}
		}
	}

	return recordsModified
}

//
//			Helper functions
//

func getIndexOfColName(colNames []string, colName string) int {
	var index int

	for i := range colNames {
		if colNames[i] == colName {
			index = i
		}
	}

	return index
}

func getTableColNames(table string) []string {
	f, err := os.Open(path + database + "/" + table)
	check(err)
	defer f.Close()

	reader := bufio.NewReader(f)
	tableMetaLine, _ := reader.ReadString('\n')
	colDefs := strings.Split(tableMetaLine, "|")

	colNames := make([]string, 0, len(colDefs))
	for i := range colDefs {
		colNames = append(colNames, strings.Fields(colDefs[i])[0])
	}
	return colNames
}

func replaceRecord(table string, record string, newRecord string) {
	lockTable(table)

	// open up another instance of table file
	input, err := ioutil.ReadFile(path + database + "/" + table)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// replace the record with our new record
	output := bytes.Replace(input, []byte(record), []byte(newRecord), -1)

	// save write to the new file
	if err = ioutil.WriteFile(path+database+"/"+table, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	unlockTable(table)
}

func removeRecord(table string, record string) {
	lockTable(table)

	// open up another instance of table file
	input, err := ioutil.ReadFile(path + database + "/" + table)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// replace the record with our new record
	output := bytes.Replace(input, []byte(record), []byte(""), -1)

	// save write to the new file
	if err = ioutil.WriteFile(path+database+"/"+table, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	unlockTable(table)
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

func lockTable(table string) {}

func unlockTable(table string) {}

func dropTable() {}

func open() {}

func commit() {}

func close() {}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
