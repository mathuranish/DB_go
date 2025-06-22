package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mathuranish/DB_go/internals"
)

func main() {
	db, err := internals.NewDatabase()
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}
	defer db.Close()

	fmt.Println("B+ Tree Database CLI")
	fmt.Println("Commands:")
	fmt.Println("  put <key> <value> - Insert a key-value pair")
	fmt.Println("  get <key>         - Retrieve a value by key")
	fmt.Println("  exit              - Exit the program")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			fmt.Println("Error: Empty command")
			continue
		}

		command := strings.ToLower(args[0])
		switch command {
		case "put":
			if len(args) < 3 {
				fmt.Println("Error: Usage: put <key> <value>")
				continue
			}
			key, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Error: Key must be an integer")
				continue
			}
			value := strings.Join(args[2:], " ") // Join remaining args as value
			db.Put(key, []byte(value))
			fmt.Printf("Inserted key %d with value %q\n", key, value)

		case "get":
			if len(args) != 2 {
				fmt.Println("Error: Usage: get <key>")
				continue
			}
			key, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Error: Key must be an integer")
				continue
			}
			if value, ok := db.Get(key); ok {
				fmt.Printf("Key %d: %s\n", key, string(value))
			} else {
				fmt.Printf("Key %d not found\n", key)
			}

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Error: Unknown command. Use put, get, or exit")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}