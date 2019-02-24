// Package generator ...
package generator

import (
	"errors"
	"sqlit/io"
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
	}

	return operation
}

func generateCreateDatabase(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	assert := func() error {
		if io.CheckIfDatabaseExists(name) == true {
			return errors.New("!Failed to create database " + name + " because it already exists.")
		}
		return nil
	}

	invoke := func() (string, error) {
		err := io.CreateDatabase(name)
		return "Database " + name + " created.", err
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateDropDatabase(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	assert := func() error {
		if io.CheckIfDatabaseExists(name) == false {
			return errors.New("!Failed to delete " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		err := io.DeleteDatabase(name)
		return "Database " + name + " deleted.", err
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateUseDatabase(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "DATABASE_NAME")

	assert := func() error {
		if io.CheckIfDatabaseExists(name) == false {
			return errors.New("!Failed to use database " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		io.UseDatabase(name)
		return "Using database " + name, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateCreateTable(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")
	columns := getAllTokensOfName(statement, "COL_NAME")
	constraints := getAllTokensOfName(statement, "COL_TYPE")

	// TODO: move these two calls to tokenizer?
	columns = sanitizeColumns(columns)
	constraints = sanitizeConstraints(constraints)

	assert := func() error {
		if io.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to create table " + name + " because no database is in use.")
		}

		if io.CheckIfTableExists(name) == true {
			return errors.New("!Failed to create table " + name + " because it already exists.")
		}
		return nil
	}

	invoke := func() (string, error) {
		io.CreateTable(name, columns, constraints)
		return "Table " + name + " created.", nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateDropTable(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")

	assert := func() error {
		if io.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to delete table " + name + " because no database is in use.")
		}

		if io.CheckIfTableExists(name) == false {
			return errors.New("!Failed to delete table " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		io.DropTable(name)
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
		if io.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to alter table " + name + " because no database is in use.")
		}

		if io.CheckIfTableExists(name) == false {
			return errors.New("!Failed to alter table " + name + " because it does not exist.")
		}

		return nil
	}

	invoke := func() (string, error) {
		io.AlterTable(name, method, column, constraint)
		return "Table " + name + " modified.", nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateSelect(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")

	assert := func() error {
		if io.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + name + " because no database is in use.")
		}

		if io.CheckIfTableExists(name) == false {
			return errors.New("!Failed to query table " + name + " because it does not exist.")
		}
		return nil
	}

	invoke := func() (string, error) {
		result := io.SelectAll(name)
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

func generateInsert(statement tokenizer.Statement) Operation {
	name := getFirstTokenOfName(statement, "TABLE_NAME")

	assert := func() error {
		if io.CheckIfAnyDatabaseIsInUse() == false {
			return errors.New("!Failed to query table " + name + " because no database is in use.")
		}

		if io.CheckIfTableExists(name) == false {
			return errors.New("!Failed to query table " + name + " because it does not exist.")
		}

		// if io.CheckIfTypesMatch() == false {
		// 	return errors.New("!Failed to query table " + name + " because of type mismatch.")
		// }

		return nil
	}

	invoke := func() (string, error) {
		result := io.SelectAll(name)
		return result, nil
	}

	return Operation{Assert: assert, Invoke: invoke}
}

/*
 *	Helper Functions
 */
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
	return "err"
}

func sanitizeColumns(columns []string) []string {
	for i := 0; i < len(columns); i++ {
		columns[i] = strings.Replace(columns[i], "(", "", 1)
	}
	return columns
}

func sanitizeConstraints(constraints []string) []string {
	for i := 0; i < len(constraints); i++ {
		constraints[i] = strings.Replace(constraints[i], ")", "", 1)
		constraints[i] = strings.Replace(constraints[i], ",", "", -1)
	}
	return constraints
}
