# Database in GO

A simple database implementation in Go using B+ trees to explore how databases manage data internally. This project serves as a learning tool to understand the mechanics of databases internals, file I/O for persistence, and Go's code structure.

### Features

- B+ Tree Implementation: A basic B+ tree structure for efficient data storage and retrieval, with self-balancing properties.
- Leaf and Non-Leaf Nodes: Separation of leaf nodes (storing key-value pairs) and non-leaf nodes (managing tree structure).
- File I/O: Persistent storage using a database.db file, with fixed-size pages (4096 bytes) to align with OS block sizes.
- Thread Safety: Mutexes ensure safe concurrent access to the database file.
- CLI Interface: Interactive terminal interface for put <key> <value>, get <key>, and exit commands.
- Go Code Structure: Organized package structure (internals, cmd) to learn Go's best practices.

## Notes
![alt text](https://github.com/mathuranish/DB_go/blob/master/notes_1.png?raw=true)
![alt text](https://github.com/mathuranish/DB_go/blob/master/notes_2.png?raw=true)

## Learnings
Through this project, I explored and implemented:

- How B+ trees manage data and maintain self-balancing properties.
- The distinction between leaf and non-leaf nodes in B+ trees and their roles.
- Go's code organization, including packages and modular design.
- File I/O operations for storing data in a .db file, simulating database persistence.
- The use of page sizes to optimize data storage and align with OS block sizes.
- Thread safety using mutexes to handle concurrent file access.

## Project Structure
```
database/
├── cmd/
│   └── main.go           # CLI entry point for user interaction
├── internals/
│   ├── api.go            # Database API (Put, Get, Close)
│   ├── disk.go           # Disk I/O for managing .db file
│   ├── node.go           # B+ tree node structure and methods
│   ├── operations.go     # B+ tree operations (search, insert)
├── go.mod                # Go module definition
├── README.md             # Project documentation
├── database.db           # Database file (generated on first run)
```
## Getting Started
### Prerequisites

Go 1.24.3 or later

### Installation

Clone the repository:git clone https://github.com/mathuranish/DB_go.git
```
cd DB_go
```
Initialize the Go module (if not already done):go mod init database

### Running the Database

Build and run the CLI:
```
go run cmd/main.go
```

Use the interactive CLI:
```
Commands:
  put <key> <value> - Insert a key-value pair
  get <key>         - Retrieve a value by key
  exit              - Exit the program
> put 1 hello world
Inserted key 1 with value "hello world"
> get 1
Key 1: hello world
> exit
Exiting...
```

_Data is stored in database.db and persists between runs._

## Notes

- Value Size Limit: Values are capped at 256 bytes (configurable in disk.go).
- Persistence: The database.db file stores all data. To start fresh, delete it before running.
- Thread Safety: The database is thread-safe for file operations but designed for single-user CLI interaction.

### Future Improvements

- Add delete operation for removing key-value pairs.
- Implement a write ahead log.
- Support configurable page sizes and value limits.

### Acknowledgments
Inspired by the desire to understand database internals and B+ tree mechanics. Built as a learning exercise in Go.
