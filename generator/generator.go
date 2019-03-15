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
	"strings"
)

// Operation ...
type Operation struct {
	Assert func() (err error)
	Invoke func() (success string, err error)
}

// Generate ...
func Generate(statement tokenizer.Statement) Operation {
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
	case parser.Types["INSERT"]:
		operation = generateInsert(statement)
	case parser.Types["UPDATE"]:
		operation = generateUpdate(statement)
	}

	return operation
}

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
		result := diskio.SelectAll(name)
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateInsert(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")
	values := getAllTokensOfName(statement, "VALUE")
	values = removeCommas(values)
	values = removeOuterParenthesisFromValues(values)

	assert := func() error {
		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + name + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(name) == false {
			return errors.New("!Failed to query table " + name + " because it does not exist.")
		}

		// TODO: assert insertions respect constraints, right amount
		// if diskio.CheckIfTypesMatch() == false {
		// 	return errors.New("!Failed to query table " + name + " because of type mismatch.")
		// }

		return nil
	}

	invoke := func() (string, error) {
		diskio.InsertRecord(name, values)
		result := "1 new record inserted."
		// result := diskio.SelectAll(name)
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateUpdate(statement tokenizer.Statement) Operation {
	// fmt.Println("HIT")

	table := getFirstTokenOfName(statement, "TABLE_NAME")
	whereCol := getSecondTokenOfName(statement, "COL_NAME")
	whereValue := getSecondTokenOfName(statement, "COL_VALUE")
	toCol := getFirstTokenOfName(statement, "COL_NAME")
	toValue := getFirstTokenOfName(statement, "COL_VALUE")

	assert := func() error {
		if diskio.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + table + " because no database is in use.")
		}

		if diskio.CheckIfTableExists(table) == false {
			return errors.New("!Failed to query table " + table + " because it does not exist.")
		}

		// TODO: better assertions
		return nil
	}

	invoke := func() (string, error) {
		result := diskio.UpdateRecord(table, whereCol, whereValue, toCol, toValue)
		result = result + " record modified." +  diskio.SelectAll(table)
  	//	result = strings.
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
		// fmt.Println("Token.Name: " + token.Name)
		if token.Name == name {
			return token.Special
		}
	}
	// TODO: handle err
	panic("Token " + name + " doesn't exist")
	// fmt.Println("Token " + name + " doesn't exist")
	// return ""
}

func getSecondTokenOfName(statement tokenizer.Statement, name string) string {
	// fmt.Println("hit")
	foundFirst := false
	for _, token := range statement.Tokens {
		if token.Name == name {
			if foundFirst == true {
				return token.Special
			}

			foundFirst = true
		}
	}
	// TODO: handle err
	panic("Token " + name + " doesn't exist")
	// fmt.Println("Token " + name + " doesn't exist")
	// return ""
}

//	@in	    (pid int, | name varchar(20), | price float)
//  @out		pid int, | name varchar(20), | price float
func removeOuterParenthesis(columns []string, constraints []string) ([]string, []string) {
	columns[0] = strings.Replace(columns[0], "(", "", 1)
	constraints[len(constraints)-1] = strings.Replace(constraints[len(constraints)-1], ")", "", 1)

	return columns, constraints
}

// TOFIX: already   gone
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
