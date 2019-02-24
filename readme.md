<!-- markdownlint-disablee -->
# CS 457 Programming Assignment 1: Metadata Management

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

Then move the project directory to the default go workspace directory:

```sh
> mv sqlite $home/go/src
> cd $home/go/src
```

And run! 

```sh
> go run sqlit
# or 
> go install sqlit
``` 

The clean flag  deletes all databases (stored in the sqlit/tmp directory) before running.

```sh
> go run sqlit --clean
```

If these steps don't work, there might be a problem with your $GOPATH. Check the docs: 
[https://golang.org/doc/code.html](https://golang.org/doc/code.html) 

## Implementation

The project is designed after sqlite's architecture. The main loop reads in a string which is piped through tokenizer.go, parser.go, and generator.go, into a disk operation. Much more detailed documentation is available in the comments.

## Organizing multiple databases

Databases are represented as directories, just as mentioned in the project spec. Currently they're nested within the tmp/ directory. Inside each is a .meta file with creation details. The program makes checks to prevent duplicate databases or other errors from occuring. The name of the database that is being `USE`'d by the system is stored in memory only.

## Managing multiple tables

Each table is a file nested within its database. The constraint metadata is stored in the first line of the table file, again, just as mentioned in the project spec. One goal is to implement all tables in a single file, to allow for pagination and large tables. 


## Resources

SQLite Architecture
[https://www.sqlite.org/draft/arch.html](https://www.sqlite.org/draft/arch.html)

SQLite SQL Syntax Diagrams
[https://www.sqlite.org/draft/lang.html](https://www.sqlite.org/draft/lang.html)

Multi platform build script
[https://gist.github.com/DimaKoz](https://gist.github.com/DimaKoz/06b7475317b12e7ffa724ef0e115a4ec)
