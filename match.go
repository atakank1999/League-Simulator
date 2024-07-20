package main

import (
	"log"
	"math/rand"
)

type Match struct {
	ID         int
	HomeTeamID int
	AwayTeamID int
	HomeGoals  int
	AwayGoals  int
	IsPlayed   bool
}

func insertMatches() {
	teams := getTeams()
	n := len(teams)

	if n%2 != 0 {

		teams = append(teams, Team{ID: 0})
		n++
	}

	rand.Shuffle(n, func(i, j int) {
		teams[i], teams[j] = teams[j], teams[i]
	})

	weeks := n - 1
	halfSize := n / 2

	matches := make([]Match, 0)

	for week := 0; week < weeks; week++ {
		for i := 0; i < halfSize; i++ {
			home := teams[i]
			away := teams[n-1-i]

			if home.ID == 0 || away.ID == 0 {
				continue
			}

			match := Match{
				HomeTeamID: home.ID,
				AwayTeamID: away.ID,
			}
			matches = append(matches, match)
		}

		last := teams[n-1]
		copy(teams[2:], teams[1:n-1])
		teams[1] = last
	}

	for _, match := range matches {
		_, err := db.Exec("INSERT INTO matches (home_team_id, away_team_id) VALUES (?, ?)", match.HomeTeamID, match.AwayTeamID)
		if err != nil {
			log.Fatalf("Error inserting match: %s", err)
		}
	}

	for _, match := range matches {

		_, err := db.Exec("INSERT INTO matches (home_team_id, away_team_id) VALUES (?, ?)", match.AwayTeamID, match.HomeTeamID)
		if err != nil {
			log.Fatalf("Error inserting reverse match: %s", err)
		}
		
	}

}

func getMatches() []Match {
	rows, err := db.Query("SELECT id, home_team_id, away_team_id, home_goals, away_goals, is_played FROM matches")
	if err != nil {
		log.Fatalf("Error getting matches: %s", err)
	}
	defer rows.Close()

	matches := make([]Match, 0)
	for rows.Next() {
		match := Match{}
		err := rows.Scan(&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.HomeGoals, &match.AwayGoals, &match.IsPlayed)
		if err != nil {
			log.Fatalf("Error scanning match: %s", err)
		}
		matches = append(matches, match)
	}
	return matches
}
