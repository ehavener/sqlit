/* UNR CS 457 | SPRING 2019 | emerson@nevada.unr.edu */

// Package generator determines a block of assertions and a block of operations
// needed to perform a statement. This is our high level analog of SQLite's
// "bytecode generator".
package generator

import (
	"errors"
	// "fmt"
	"sqlit/diskio"
	"sqlit/parser"
	"sqlit/tokenizer"
	"strconv"
	"strings"
)

// Operation ...
type Operation struct {
	Assert func() (err error)
	Invoke func() (success string, err error)
}

// Generate ...
func Generate(statement tokenizer.Statement, inTransactionMode bool) Operation {

	operation := Operation{}

	switch statement.Type {
	case parser.Types["CREATE_DATABASE"]:
		operation = generateCreateDatabase(statement)
	case parser.Types["DROP_DATABASE"]:
		operation = generateDropDatabase(statement)
	case parser.Types["USE_DATABASE"]:
		operation = generateUseDatabase(statement)
	case parser.Types["CREATE_TABLE"]:
		operation = generateCreateTable(statement)
	case parser.Types["ALTER_TABLE"]:
		operation = generateAlterTable(statement)
	case parser.Types["DROP_TABLE"]:
		operation = generateDropTable(statement)
	case parser.Types["SELECT"]:
		operation = generateSelect(statement)
	case parser.Types["SELECT_INNER"]:
		operation = generateSelectInnerJoin(statement)
	case parser.Types["SELECT_LEFT"]:
		operation = generateSelectLeftJoin(statement)
	case parser.Types["INSERT"]:
		operation = generateInsert(statement)
	case parser.Types["UPDATE"]:
		operation = generateUpdate(statement, inTransactionMode)
	case parser.Types["DELETE"]:
		operation = generateDelete(statement)
		//	case parser.Types["BEGIN"]:
		// operation = generateBegin(statement)
		// 	case parser.Types["COMMIT"]:
		// operation = generateCommit(statement)
	}

	return operation
}

// func generateBegin(statement tokenizer.Statement) Operation {

// 	assert := func() error {
// 		return nil
// 	}

// 	invoke := func() (string, error) {
// 		return "Transaction starts.", nil
// 	}

// 	return Operation{Assert: assert, Invoke: invoke}
// }

// func generateCommit(statement tokenizer.Statement) Operation {

// 	assert := func() error {
// 		return nil
// 	}

// 	invoke := func() (string, error) {
// 		return "Attempting to commit transaction.", nil
// 	}

// 	return Operation{Assert: assert, Invoke: invoke}
// }

func generateCreateDatabase(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	assert := func() error {
		if diskio.CheckIfDatabaseExists(name) == true {
			return errors.New("!Failed to create database " + name + " because it already exists.")
		}
		return nil
	}

	invoke := func() (string, error) {
		err := diskio.CreateDatabase(name)
		diskio.CreateDatabaseMeta(name)
		return "Database " + name + " created.", err
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateDropDatabase(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	assert := func() error {
		if diskio.CheckIfDatabaseExists(name) == false {
			return errors.New("!Failed to delete " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		err := diskio.DeleteDatabase(name)
		return "Database " + name + " deleted.", err
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateUseDatabase(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	assert := func() error {
		if diskio.CheckIfDatabaseExists(name) == false {
			return errors.New("!Failed to use database " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		diskio.UseDatabase(name)
		return "Using database " + name, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateCreateTable(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")
	columns := getAllTokensOfName(statement, "COL_NAME")
	constraints := getAllTokensOfName(statement, "COL_TYPE")

	// TODO: move these to tokenizer?
	columns, constraints = removeOuterParenthesis(columns, constraints)

	columns = removeCommas(columns)
	constraints = removeCommas(constraints)

	assert := func() error {
		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to create table " + name + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(name) == true {
			return errors.New("!Failed to create table " + name + " because it already exists.")
		}
		return nil
	}

	invoke := func() (string, error) {
		diskio.CreateTable(name, columns, constraints)
		return "Table " + name + " created.", nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateDropTable(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")

	assert := func() error {
		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to delete table " + name + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(name) == false {
			return errors.New("!Failed to delete table " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		diskio.DropTable(name)
		return "Table " + name + " deleted.", nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateAlterTable(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")
	method := getFirstTokenOfName(statement, "ADD_COL")
	column := getFirstTokenOfName(statement, "COL_NAME")
	constraint := getFirstTokenOfName(statement, "COL_TYPE")

	assert := func() error {
		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to alter table " + name + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(name) == false {
			return errors.New("!Failed to alter table " + name + " because it does not exist.")
		}

		return nil
	}

	invoke := func() (string, error) {
		diskio.AlterTable(name, method, column, constraint)
		return "Table " + name + " modified.", nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateSelect(statement tokenizer.Statement) Operation {

	name := getFirstTokenOfName(statement, "TABLE_NAME")
	clause := statement.Tokens[1].Name

	assert := func() error {
		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + name + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(name) == false {
			return errors.New("!Failed to query table " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		var result string

		if clause == "ALL" {
			result = diskio.SelectAll(name)
		} else {
			name := getFirstTokenOfName(statement, "TABLE_NAME")

			var colNames []string

			colNames = append(colNames, getFirstTokenOfName(statement, "COL_NAME"))
			colNames = append(colNames, getSecondTokenOfName(statement, "COL_NAME"))

			whereColName := getThirdTokenOfName(statement, "COL_NAME")
			whereColVal := getFirstTokenOfName(statement, "COL_VALUE")

			colNames = removeCommas(colNames)

			result = diskio.SelectWhere(name, colNames, whereColName, whereColVal)
		}
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateSelectInnerJoin(statement tokenizer.Statement) Operation {

	leftTableName := getFirstTokenOfName(statement, "TABLE_NAME")
	rightTableName := getSecondTokenOfName(statement, "TABLE_NAME")

	assert := func() error {
		return nil
	}

	invoke := func() (string, error) {
		setOne := diskio.SelectSet(leftTableName)
		setTwo := diskio.SelectSet(rightTableName)
		innerJoinOfSets := InnerJoin(setOne, setTwo, "id", "employeeID")

		result := diskio.SerializeSet(innerJoinOfSets)

		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateSelectLeftJoin(statement tokenizer.Statement) Operation {
	leftTableName := getFirstTokenOfName(statement, "TABLE_NAME")
	rightTableName := getSecondTokenOfName(statement, "TABLE_NAME")

	assert := func() error {
		return nil
	}

	invoke := func() (string, error) {
		setOne := diskio.SelectSet(leftTableName)
		setTwo := diskio.SelectSet(rightTableName)
		leftJoinOfSets := LeftJoin(setOne, setTwo, "id", "employeeID")
		result := diskio.SerializeSet(leftJoinOfSets)
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

// InnerJoin ...
func InnerJoin(setOne diskio.Set, setTwo diskio.Set, setOneColName string, setTwoColName string) diskio.Set {
	var columnDefs []diskio.ColumnDef

	for _, columnDef := range setOne.ColumnDefs {
		columnDefs = append(columnDefs, columnDef)
	}

	for _, columnDef := range setTwo.ColumnDefs {
		columnDefs = append(columnDefs, columnDef)
	}

	var setOneColIndex int
	var setTwoColIndex int

	for index, columnDef := range setOne.ColumnDefs {
		if columnDef.ColumnName == setOneColName {
			setOneColIndex = index
		}
	}

	for index, columnDef := range setTwo.ColumnDefs {
		if columnDef.ColumnName == setTwoColName {
			setTwoColIndex = index
		}
	}

	records := make([][]string, 5)
	recordsRowIndex := 0

	for _, setOneRecords := range setOne.Records {
		for _, setTwoRecords := range setTwo.Records {
			if setOneRecords[setOneColIndex] == setTwoRecords[setTwoColIndex] {
				recordsRowIndex++
				records[recordsRowIndex] = make([]string, 4)
				records[recordsRowIndex][0] = setOneRecords[0]
				records[recordsRowIndex][1] = setOneRecords[1]
				records[recordsRowIndex][2] = setTwoRecords[0]
				records[recordsRowIndex][3] = setTwoRecords[1]
			}
		}
	}

	set := diskio.Set{Name: "inner-join", ColumnDefs: columnDefs, Records: records}

	return set
}

// LeftJoin ...
func LeftJoin(leftSet diskio.Set, rightSet diskio.Set, leftSetColName string, rightSetColName string) diskio.Set {
	innerJoinOfSets := InnerJoin(leftSet, rightSet, leftSetColName, rightSetColName)

	var columnDefsOfSets []diskio.ColumnDef

	for _, columnDef := range leftSet.ColumnDefs {
		columnDefsOfSets = append(columnDefsOfSets, columnDef)
	}

	for _, columnDef := range rightSet.ColumnDefs {
		columnDefsOfSets = append(columnDefsOfSets, columnDef)
	}

	var leftSetColIndex int
	var rightSetColIndex int

	for index, columnDef := range leftSet.ColumnDefs {
		if columnDef.ColumnName == leftSetColName {
			leftSetColIndex = index
		}
	}

	for index, columnDef := range rightSet.ColumnDefs {
		if columnDef.ColumnName == rightSetColName {
			rightSetColIndex = index
		}
	}

	leftExclusiveRecords := make([][]string, 5)
	leftExclusiveRecordsTailRowIndex := -1

	for _, leftSetRowRecord := range leftSet.Records {
		leftExclusive := true

		for _, rightSetRowRecord := range rightSet.Records {
			if leftSetRowRecord[leftSetColIndex] == rightSetRowRecord[rightSetColIndex] {
				leftExclusive = false
			}
		}

		if leftExclusive == true {
			leftExclusiveRecordsTailRowIndex++
			leftExclusiveRecords[leftExclusiveRecordsTailRowIndex] = make([]string, 4)
			leftExclusiveRecords[leftExclusiveRecordsTailRowIndex][0] = leftSetRowRecord[0]
			leftExclusiveRecords[leftExclusiveRecordsTailRowIndex][1] = leftSetRowRecord[1]
			leftExclusiveRecords[leftExclusiveRecordsTailRowIndex][2] = ""
			leftExclusiveRecords[leftExclusiveRecordsTailRowIndex][3] = ""
		}
	}

	leftJoinRecords := append(innerJoinOfSets.Records, leftExclusiveRecords...)

	leftJoinOfSets := diskio.Set{Name: "left-join", ColumnDefs: columnDefsOfSets, Records: leftJoinRecords}

	return leftJoinOfSets
}

func generateInsert(statement tokenizer.Statement) Operation {
	tableName := getFirstTokenOfName(statement, "TABLE_NAME")
	values := getAllTokensOfName(statement, "VALUE")
	values = removeCommas(values)
	values = removeOuterParenthesisFromValues(values)
	values = removeQuotes(values)

	assert := func() error {
		if diskio.CheckIfTableIsLockedByOtherProcess(tableName) == true {
			return errors.New("Error: Table " + tableName + " is locked!")
		} else {
			diskio.LockTable(tableName)
		}

		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + tableName + " because no database is in use.")
		}

		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + tableName + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(tableName) == false {
			return errors.New("!Failed to query table " + tableName + " because it does not exist.")
		}

		return nil
	}

	invoke := func() (string, error) {
		diskio.UnlockTable(tableName)

		diskio.InsertRecord(tableName, values)
		result := "1 new record inserted."
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateUpdate(statement tokenizer.Statement, inTransactionMode bool) Operation {
	tableName := getFirstTokenOfName(statement, "TABLE_NAME")
	whereCol := getSecondTokenOfName(statement, "COL_NAME")
	whereValue := getSecondTokenOfName(statement, "COL_VALUE")
	toCol := getFirstTokenOfName(statement, "COL_NAME")
	toValue := getFirstTokenOfName(statement, "COL_VALUE")

	whereValue = strings.Replace(whereValue, "'", "", -1)
	toValue = strings.Replace(toValue, "'", "", -1)

	assert := func() error {
		if diskio.CheckIfTableIsLockedByOtherProcess(tableName) == true {
			return errors.New("!Error: Table " + tableName + " is locked!")
		} else {
			diskio.LockTable(tableName)
		}

		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + tableName + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(tableName) == false {
			return errors.New("!Failed to query table " + tableName + " because it does not exist.")
		}

		if inTransactionMode {

		}

		return nil
	}

	invoke := func() (string, error) {
		diskio.UnlockTable(tableName)

		recordsModified := diskio.UpdateRecord(tableName, whereCol, whereValue, toCol, toValue)
		result := strconv.Itoa(recordsModified)
		result = result + " record(s) modified."
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateDelete(statement tokenizer.Statement) Operation {

	table := getFirstTokenOfName(statement, "TABLE_NAME")
	whereCol := getFirstTokenOfName(statement, "COL_NAME")
	whereValue := getFirstTokenOfName(statement, "COL_VALUE")
	comparator := statement.Tokens[5].Name

	whereValue = strings.Replace(whereValue, "'", "", -1)

	assert := func() error {
		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + table + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(table) == false {
			return errors.New("!Failed to query table " + table + " because it does not exist.")
		}

		return nil
	}

	invoke := func() (string, error) {
		recordsDeleted := diskio.DeleteRecord(table, whereCol, whereValue, comparator)
		result := strconv.Itoa(recordsDeleted)
		result = result + " record(s) deleted."
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

//
//			Helper functions
//

func getAllTokensOfName(statement tokenizer.Statement, name string) []string {
	var specials []string
	for _, token := range statement.Tokens {
		if token.Name == name {
			specials = append(specials, token.Special)
		}
	}
	return specials
}

func getFirstTokenOfName(statement tokenizer.Statement, name string) string {
	for _, token := range statement.Tokens {
		if token.Name == name {
			return token.Special
		}
	}

	panic("Token " + name + " doesn't exist")
}

func getSecondTokenOfName(statement tokenizer.Statement, name string) string {
	foundFirst := false
	for _, token := range statement.Tokens {
		if token.Name == name {
			if foundFirst == true {
				return token.Special
			}

			foundFirst = true
		}
	}

	panic("Token " + name + " doesn't exist")
}

func getThirdTokenOfName(statement tokenizer.Statement, name string) string {
	foundFirst := false
	foundSecond := false
	for _, token := range statement.Tokens {
		if token.Name == name {
			if foundFirst == true {
				if foundSecond == true {

					return token.Special
				}

				foundSecond = true
			}
			foundFirst = true
		}
	}

	panic("Token " + name + " doesn't exist")
}

//	@in	    (pid int, | name varchar(20), | price float)
//  @out		pid int, | name varchar(20), | price float
func removeOuterParenthesis(columns []string, constraints []string) ([]string, []string) {
	columns[0] = strings.Replace(columns[0], "(", "", 1)
	constraints[len(constraints)-1] = strings.Replace(constraints[len(constraints)-1], ")", "", 1)

	return columns, constraints
}

func removeOuterParenthesisFromValues(values []string) []string {
	values[0] = strings.Replace(values[0], "(", "", 1)
	values[len(values)-1] = strings.Replace(values[len(values)-1], ")", "", 1)

	return values
}

func removeCommas(values []string) []string {
	for i := 0; i < len(values); i++ {
		values[i] = strings.Replace(values[i], ",", "", 1)
	}
	return values
}

func removeQuotes(values []string) []string {
	for i := 0; i < len(values); i++ {
		values[i] = strings.Replace(values[i], "'", "", -1)
	}
	return values
}
