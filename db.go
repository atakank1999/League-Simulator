package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./football_league.db")
	if err != nil {
		log.Fatal(err)
	}
	createTables()
}

func createTables() {

	createTeamTable := `CREATE TABLE IF NOT EXISTS teams (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
		offensive_strength INTEGER,
		defensive_strength INTEGER,
		wins INTEGER DEFAULT 0,
		draws INTEGER DEFAULT 0,
		loses INTEGER DEFAULT 0,
        points INTEGER DEFAULT 0,
        goals_for INTEGER DEFAULT 0,
        goals_against INTEGER DEFAULT 0
    );`

	createMatchTable := `CREATE TABLE IF NOT EXISTS matches (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        home_team_id INTEGER,
        away_team_id INTEGER,
        home_goals INTEGER DEFAULT 0,
        away_goals INTEGER DEFAULT 0,
		is_played BOOLEAN DEFAULT FALSE,
        FOREIGN KEY(home_team_id) REFERENCES teams(id),
        FOREIGN KEY(away_team_id) REFERENCES teams(id)
    );`

	_, err := db.Exec(createTeamTable)
	if err != nil {
		log.Fatalf("Error creating teams table: %s", err)
	}

	_, err = db.Exec(createMatchTable)
	if err != nil {
		log.Fatalf("Error creating matches table: %s", err)
	}
}
func resetDatabase() {
    _, err := db.Exec("DELETE FROM matches")
    if err != nil {
        log.Fatalf("Error clearing matches table: %s", err)
    }

    _, err = db.Exec("UPDATE teams SET points = 0, goals_for = 0, goals_against = 0, wins = 0, draws = 0, loses = 0")
    if err != nil {
        log.Fatalf("Error resetting teams: %s", err)
    }
}
