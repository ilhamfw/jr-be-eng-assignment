package main

import (
	
	"fmt"
	"os"

	"kompit-recruitment/jr-be-eng-assignment/csv"
	"kompit-recruitment/jr-be-eng-assignment/database"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Failed to initialize the database:", err)
		return
	}
	defer db.Close()

	command := os.Args[1]
	switch command {
	case "input":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go input <filename>")
			return
		}
		filePath := os.Args[2]
		err := csv.ProcessCSVFile(filePath, db)
		if err != nil {
			fmt.Println("Error processing CSV file:", err)
			return
		}
		fmt.Println("CSV file processed successfully and data stored in the database.")
	case "leaderboard":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go leaderboard <competition_id> <filename>")
			return
		}
		csv.Leaderboard(os.Args[2:])
	
	default:
		fmt.Fprintf(os.Stderr, "Invalid command: %s\n", command)
		os.Exit(1)
	}
}




