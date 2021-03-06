<!-- markdownlint-disable -->

# CS 457 Programming Assignment 2: Basic Data Manipulation

## Running

The quickest way to test the program is to just run a precompiled build

```sh
> cd sqlit/
> ./sqlit-linux-386
```

## Building

To build the source, first make sure go is installed:

```sh
> sudo apt install golang-go
> go version
```

Export your gopath

```sh
> export GOPATH=$HOME/go
```

Then move the project directory to the default go workspace directory:

```sh
> mkdir ~/go && mkdir ~/go/src
> mv sqlit $HOME/go/src
> cd $HOME/go/src
```

And build & run!

```sh
> go install sqlit
> cd ~/go/bin # default build output dir
> ./sqlit
```

The clean flag  deletes all databases (stored in the sqlit/tmp directory) before running.

```sh
> go run sqlit --clean
```

A test script can be piped in as so:

```sh
> go run sqlit --clean < test/PA2_test.sql
```

If these steps don't work, there might be a problem with your $GOPATH. Check the docs:
[https://golang.org/doc/code.html](https://golang.org/doc/code.html)

## Implementation

The project is designed after sqlite's architecture. The main loop reads in a string which is piped through tokenizer.go, parser.go, and generator.go, into a disk operation. Much more detailed documentation is available in the comments.

## Organizing multiple databases (PA1)

Databases are represented as directories, just as mentioned in the project spec. Currently they're nested within the tmp/ directory. Inside each is a .meta file with creation details. The program makes checks to prevent duplicate databases or other errors from occuring. The name of the database that is being `USE`'d by the system is stored in memory only.

## Managing multiple tables (PA1)

Each table is a file nested within its database. The constraint metadata is stored in the first line of the table file, again, just as mentioned in the project spec. One goal is to implement all tables in a single file, to allow for pagination and large tables.

## Tuple insertion, deletion, modification, and query (PA2)

The primary functons that handle tuple CRUD are InsertRecord(), SelectWhere(), UpdateRecord(), and DeleteRecord(). They all behave similarly, and live in the diskio library. They are fairly abstract, but some compromises were made to reduce complexity. The only operators that are implemented are "*", ">", "=", and "=", and only the update/select clause formats used in the test files have been tested.

DeleteRecord() will
- open the table file in Read mode
- read the table's metadata (column names)
- determine the offset of the clause's column
- iterate through all records in the table, reading their values at the clause column
- compare their values to the clause
- records are physically deleted by opening a new file reader instance in write mode and erasing the selected record by it's bytestring

UpdateRecord() will behave like DeleteRecord(), only replacing the record with a newly constructed bytestring.

InsertRecord() will
- open the table file in append mode
- construct a record from the given tuple, in a similar format to table metadata (pipe delimited)
- append the new record to the last line of the table file

SelectWhere() will
- construct a temporary table of the columns specified
- read through a table, accumulating the values at the specified columns which pass the clause
- output the temporary table

## Table Joins (PA3)

### A new data structure is added, Sets.

Sets contain a 2D matrix of column values as well as a 1D matrix of column definitions. Sets make it easier to implement in-memory operations like joins. Along with them are helper methods to translate between structured and stringified data.


### An inner join is preformed using nested loops.

First, the column definitions are read from the table. They're then concatenated to create a new column definition for our join set. The new horizontal indices of the clause columns (e.g. id and employee id) are stored, as well as their original table indices. We then completely iterate through both tables with nested for loops, row by row. If two  values pass the given expression (in this case just ==), then we allocate space for and construct a new record in our join set. This record consists of all other values in the matching records, their offsets adjusted to match our concatenated column defs. The set is then returned to be serialized and output. This is approximately O(n^2) where n is every record in both tables which is not great and could certainly be improved upon.


### A left (outer) join simply extends an inner join.

An outer join is preformed by first inner-joining the two tables. The inner-join set's records are then iterated though again, and unmatched records from the leftmost table are appended to the end of the set. A  real world implementation would probably use a specialized algorithm, one that populates the join's columns individually. It may for example, allow the leftmost set's values to bypass the condition and copy them to the set intermittently, sorting them to the end afterwards. Another, better alternative would be to store the indicies of unpaired left columns during the inner join to be appended in O(n) afterwards.

## Locking & Transactions (PA4)

This project introduces locks as a way of ensuring that a sequence of operations are completed atomically, in their defined order, and reversibly. When a begin transaction token is read, the main loop enters "transaction mode". In transaction mode, operations are still asserted as interpreted but are not executed. Instead, they are pushed into a queue. This queue continues to grow until a commit token is read, or an assertion fails. If a commit token is reached without any errors, the operation queue will then be asserted for a second time and then executed as a single transaction. This assertion will lock all resources touched by the transaction. They will then be unlocked upon any error or upon total completion. Any error will also terminate the transaction block.

Locks are persisted as a {{TABLE_NAME}}.lock file at the root database directory. Inside of each .lock is the id of the process that is currently accessing the resource. A shortcoming of my implementation is that if a process is exited while still owning a resource, the lock will not be deleted. This could be fixed rather easily but is unnecessary for this proof of concept.

## Resources

SQLite Architecture
[https://www.sqlite.org/draft/arch.html](https://www.sqlite.org/draft/arch.html)

SQLite SQL Syntax Diagrams
[https://www.sqlite.org/draft/lang.html](https://www.sqlite.org/draft/lang.html)

Multi platform build script
[https://gist.github.com/DimaKoz](https://gist.github.com/DimaKoz/06b7475317b12e7ffa724ef0e115a4ec)

Auto build go project on changes [https://github.com/canthefason/go-watcher](https://github.com/canthefason/go-watcher)
