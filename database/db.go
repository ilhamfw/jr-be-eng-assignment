package database

import (
	"database/sql"
	"fmt"
	"kompit-recruitment/jr-be-eng-assignment/structs"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "user"
	password = "password"
	dbname   = "database"
)

var db *sql.DB

// InitDB initializes PostgreSQL connection
func InitDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")

	return db, nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
		fmt.Println("Database connection closed.")
	}
}

// Exec is a wrapper around sql.DB.Exec
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

// GetMatchResultsByCompetitionID retrieves match results for a specific competition_id
func GetMatchResultsByCompetitionID(competitionID string) ([]structs.MatchResult, error) {
	// Check if the database connection is initialized
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Query to retrieve match results for a specific competition_id
	query := `
		SELECT competition_id, date, team_1, team_1_score, team_2, team_2_score
		FROM match_results
		WHERE competition_id = $1
	`

	// Execute the query
	rows, err := db.Query(query, competitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []structs.MatchResult

	// Iterate through the rows and populate the results slice
	for rows.Next() {
		var result structs.MatchResult
		err := rows.Scan(&result.CompetitionID, &result.Date, &result.Team1, &result.Team1Score, &result.Team2, &result.Team2Score)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
