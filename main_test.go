package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"kompit-recruitment/jr-be-eng-assignment/csv"
	
)

func TestLeaderboard(t *testing.T) {
	// Inisialisasi mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Menyiapkan ekspektasi query mock
	rows := sqlmock.NewRows([]string{"competition_id", "date", "team_1", "team_1_score", "team_2", "team_2_score"}).
		AddRow(1, "2023-11-11", "Team A", 1, "Team B", 0).
		AddRow(2, "2023-11-11", "Team B", 1, "Team C", 0).
		AddRow(3, "2023-11-11", "Team A", 1, "Team D", 0).
		AddRow(4, "2023-11-12", "Team B", 2, "Team A", 2).
		AddRow(5, "2023-11-12", "Team C", 2, "Team D", 1).
		AddRow(6, "2023-11-12", "Team A", 2, "Team B", 1)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	// Menentukan argumen untuk fungsi leaderboard
	args := []string{"1", "test_leaderboard.csv"}

	// Memanggil fungsi yang akan diuji
	csv.Leaderboard(args)

	// Memastikan bahwa semua ekspektasi query mock terpenuhi
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Unfulfilled expectations: %s", err)
	}
}
