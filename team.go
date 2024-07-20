package main

import "log"


type Team struct {
    ID           int
    Name         string
    OffensiveStrength     int
	DefensiveStrength     int
    Points       int
	Wins		 int
	Draws		 int
	Loses		 int
    GoalsFor     int
    GoalsAgainst int
}

func getTeamByID(id int) Team {
	rows, err := db.Query("SELECT id, name, offensive_strength, defensive_strength, points, wins, draws, loses, goals_for, goals_against FROM teams WHERE id = ?", id)
	if err != nil {
		log.Fatalf("Error getting team: %s", err)
	}
	defer rows.Close()

	team := Team{}
	for rows.Next() {
		err := rows.Scan(&team.ID, &team.Name, &team.OffensiveStrength, &team.DefensiveStrength, &team.Points, &team.Wins, &team.Draws, &team.Loses, &team.GoalsFor, &team.GoalsAgainst)
		if err != nil {
			log.Fatalf("Error scanning team: %s", err)
		}
	}
	return team

}


func insertTeams() {

	teams := []Team{
		{Name: "Manchester City", OffensiveStrength: 275, DefensiveStrength: 281, Points: 0, Wins: 0, Draws: 0, Loses: 0, GoalsFor: 0, GoalsAgainst: 0},
		{Name: "Arsenal", OffensiveStrength: 261, DefensiveStrength: 330, Points: 0, Wins: 0, Draws: 0, Loses: 0, GoalsFor: 0, GoalsAgainst: 0},
		{Name: "Liverpool", OffensiveStrength: 246, DefensiveStrength: 233, Points: 0, Wins: 0, Draws: 0, Loses: 0, GoalsFor: 0, GoalsAgainst: 0},
		{Name: "Aston Villa", OffensiveStrength: 218, DefensiveStrength: 157, Points: 0, Wins: 0, Draws: 0, Loses: 0, GoalsFor: 0, GoalsAgainst: 0},
	}
	for _, team := range teams {
		_, err := db.Exec("INSERT INTO teams (name, offensive_strength, defensive_strength, points, wins, draws, loses, goals_for, goals_against) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", team.Name, team.OffensiveStrength, team.DefensiveStrength, team.Points, team.Wins, team.Draws, team.Loses, team.GoalsFor, team.GoalsAgainst)
		if err != nil {
			log.Fatalf("Error inserting team: %s", err)
		}
	}
}

func getTeams() []Team {
	rows, err := db.Query("SELECT id, name, offensive_strength, defensive_strength, points, wins, draws, loses, goals_for, goals_against FROM teams")
	if err != nil {
		log.Fatalf("Error getting teams: %s", err)
	}
	defer rows.Close()

	teams := make([]Team, 0)
	for rows.Next() {
		team := Team{}
		err := rows.Scan(&team.ID, &team.Name, &team.OffensiveStrength, &team.DefensiveStrength, &team.Points, &team.Wins, &team.Draws, &team.Loses, &team.GoalsFor, &team.GoalsAgainst)
		if err != nil {
			log.Fatalf("Error scanning team: %s", err)
		}
		teams = append(teams, team)
	}
	
	return teams
}