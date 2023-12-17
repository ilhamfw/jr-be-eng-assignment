package structs

import "time"

// MatchResult struct to represent a match result
type MatchResult struct {
	CompetitionID string
	Date          time.Time
	Team1         string
	Team1Score    int
	Team2         string
	Team2Score    int
}