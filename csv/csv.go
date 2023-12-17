package csv

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"kompit-recruitment/jr-be-eng-assignment/database"
	"kompit-recruitment/jr-be-eng-assignment/structs"
)

// TeamStats struct to represent team statistics
type TeamStats struct {
	Name   string
	Play   int
	Win    int
	Draw   int
	Lose   int
	Points int
}

// ProcessCSVFile reads and validates match results from a CSV file
func ProcessCSVFile(filePath string, db *sql.DB) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var results []structs.MatchResult

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Check the number of columns in the CSV record
		if len(record) < 6 {
			return fmt.Errorf("invalid number of columns in the CSV record")
		}

		// Parse score values from the CSV record
		team1Score, err := strconv.Atoi(record[3])
		if err != nil {
			return fmt.Errorf("error parsing Team1Score: %v", err)
		}

		// Use strings.TrimSpace to remove extra spaces before parsing
		team2Score, err := strconv.Atoi(strings.TrimSpace(record[5]))
		if err != nil {
			return fmt.Errorf("error parsing Team2Score: %v", err)
		}

		// Create a MatchResult instance with parsed scores
		result := structs.MatchResult{
			CompetitionID: record[0],
			Team1:         record[2],
			Team1Score:    team1Score,
			Team2:         record[4],
			Team2Score:    team2Score,
		}

		// Parse date with the expected format
		result.Date, err = time.Parse("2006-01-02", record[1])
		if err != nil {
			return err
		}

		// Append the MatchResult to the results slice
		results = append(results, result)
	}

	// Insert match results into the database
	if err := InsertMatchResults(db, results); err != nil {
		return err
	}

	return nil
}

// InsertMatchResults inserts match results into the database
func InsertMatchResults(db *sql.DB, results []structs.MatchResult) error {
	for _, result := range results {
		_, err := database.Exec(`
			INSERT INTO match_results (competition_id, date, team_1, team_1_score, team_2, team_2_score)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			result.CompetitionID, result.Date, result.Team1, result.Team1Score, result.Team2, result.Team2Score)

		if err != nil {
			return fmt.Errorf("error inserting data into the database: %v", err)
		}
	}

	return nil
}

func Leaderboard(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: leaderboard <competition_id> <filename>")
		return
	}

	competitionID := args[0]
	fileName := args[1]

	// Connect to the database
	_, err := database.InitDB()
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}
	defer database.CloseDB()

	// Fetch match results for the specified competition_id
	results, err := database.GetMatchResultsByCompetitionID(competitionID)
	if err != nil {
		fmt.Printf("Error fetching match results: %v\n", err)
		return
	}

	// Process match results to generate leaderboard
	leaderboard := processLeaderboard(results)

	// Write leaderboard to CSV file
	err = writeLeaderboardToCSV(fileName, leaderboard)
	if err != nil {
		fmt.Printf("Error writing leaderboard to CSV: %v\n", err)
		return
	}

	fmt.Printf("Leaderboard successfully written to %s\n", fileName)
}

func processLeaderboard(results []structs.MatchResult) map[string]TeamStats {
	teamStats := make(map[string]TeamStats)

	// Calculate team stats based on match results
	for _, result := range results {
		updateTeamStats(teamStats, result.Team1, result.Team1Score, result.Team2Score)
		updateTeamStats(teamStats, result.Team2, result.Team2Score, result.Team1Score)
	}

	// Sort teams based on points, number of plays, and alphabetical name
	sortedTeams := sortTeams(teamStats)

	// Create a map for leaderboard output
	leaderboard := make(map[string]TeamStats)

	// Generate leaderboard with sorted teams
	for _, team := range sortedTeams {
		leaderboard[team] = teamStats[team]
	}

	return leaderboard
}

func updateTeamStats(teamStats map[string]TeamStats, team string, teamScore, opponentScore int) {
	stats, ok := teamStats[team]
	if !ok {
		stats = TeamStats{Name: team}
	}

	stats.Play++
	if teamScore > opponentScore {
		stats.Win++
		stats.Points += 3
	} else if teamScore == opponentScore {
		stats.Draw++
		stats.Points++
	} else {
		stats.Lose++
	}

	teamStats[team] = stats
}

func sortTeams(teamStats map[string]TeamStats) []string {
	// Create a slice to hold team names
	teams := make([]string, 0, len(teamStats))

	// Populate the slice with team names
	for team := range teamStats {
		teams = append(teams, team)
	}

	// Sort the teams based on points, number of plays, and alphabetical name
	sort.Slice(teams, func(i, j int) bool {
		// Sort by points in descending order
		if teamStats[teams[i]].Points != teamStats[teams[j]].Points {
			return teamStats[teams[i]].Points > teamStats[teams[j]].Points
		}

		// Sort by number of plays in descending order
		if teamStats[teams[i]].Play != teamStats[teams[j]].Play {
			return teamStats[teams[i]].Play > teamStats[teams[j]].Play
		}

		// Sort by alphabetical order of team names
		return teams[i] < teams[j]
	})

	return teams
}

func writeLeaderboardToCSV(fileName string, leaderboard map[string]TeamStats) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{"team", "play", "win", "draw", "lose", "points"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, team := range leaderboard {
		row := []string{
			team.Name,
			strconv.Itoa(team.Play),
			strconv.Itoa(team.Win),
			strconv.Itoa(team.Draw),
			strconv.Itoa(team.Lose),
			strconv.Itoa(team.Points),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
