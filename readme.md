<!-- markdownlint-disable -->
# CS 457 Programming Assignment 1: Metadata Management

## Running

### Running from source

1. From inside the `/src/sqlit` directory, run `go run sqlit` 

2. Or, run `go run sqlit --clean` to delete all databases on launch.

3. 

### Running from a build

1. From inside the `/bin` directory, run `./sqlit` 

## Building the project

- From inside the `/src/sqlit` directory, run `go install sqlit`

- To build executables for all platforms, run `chmod +x ./go-executable-build.sh` and  `./go-executable-build.sh sqlit`

## Organizing multiple databases

Databases are represented as directories. Currently they're nested within the tmp/ directory.

## Managing multiple tables

Each table is a file stored

## Implementation

## Resources

Architecture

https://www.sqlite.org/draft/arch.html

For build info and troubleshooting

[How to Write Go Code | https://golang.org/doc/code.html ](https://golang.org/doc/code.html) 

Building for all platforms

[The Script to Automate Cross-Compilation for Golang(OSX) | https://gist.github.com/DimaKoz]()

## Notes

Data Management Functionality
  1. Persistently store large datasets
  2. Efficiently query & update
    - Must handle complex questions about data
    - Must handle sophisticated updates
    - Performance matters
  3. Change structure (e.g., add attributes)
  4. Concurrency control: enable simultaneous updates
  5. Crash recovery

Data Management Concepts
1. Data independence
  – Physical independence: Can change how data is stored on disk without maintenance to applications
  – Logical independence: can change schema w/o affecting apps
2. Query optimizer and compiler
3. Transactions
  - isolation and atomicity


H1: Metadata management
-- e.g., Create/update/remove tables

H2: Basic data management
  -- e.g., Insert/update/remove tuples

H3: Data aggregates
  -- e.g., Different types of table joins

H4: Advanced topics
  -- e.g., locking, transactions


* Case doesn't matter
* Lock entire DB when writing anything?


Inner join
  SELECT DISTINCT   cname
  FROM    Product, Company
  WHERE   country = ‘USA’ AND category = ‘gadget’ 
          AND manufacturer = cname

Outer join
  SELECT   Product.name, Purchase.store
  FROM    Product LEFT OUTER JOIN Purchase ON
          Product.name = Purchase.prodName

Grouping
SELECT    product, Sum(quantity) AS TotalSales
FROM      Purchase
WHERE     price > 1
GROUP BY product

Aggregations
  select count(*) from Purchase
  select sum(quantity) from Purchase
  select avg(price) from Purchase
  select max(quantity) from Purchase
  select min(quantity) from Purchase




Relational Algebra
SQL-> Relational Algebra -> Physical Plan

Relational Algebra = Logical Plan

Select Operation (σ)
- Returns all tuples which satisfy a condition
- σc(R) 
  σ stands for selection predicate
  R stands for relation
  c is prepositional logic =, <, <=, >, >=, <>

Projection
- Eliminates columns
- ∏A1, A2, An (R)

