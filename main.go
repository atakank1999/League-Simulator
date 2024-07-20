package main

import (
	"log"
	"math/rand"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

var currentWeek = 0

func main() {
	r := gin.Default()

	initDatabase()
	resetDatabase()

	r.GET("/simulate-week", simulateWeekHandler)
	r.GET("/league-table", leagueTableHandler)
	r.GET("/play-all-matches", playAllMatchesHandler)
	r.GET("/insert-teams", insertTeamsHandler)

	r.Run(":8080")
}

func insertTeamsHandler(c *gin.Context) {
	teams := getTeams()
	league := leagueTable()

	if len(teams) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"currentWeek": currentWeek,
			"league":      league,
		})
		return
	}
	insertTeams()
	league = leagueTable()

	c.JSON(http.StatusOK, gin.H{
		"currentWeek": currentWeek,
		"league":      league,
	})
}

func simulateWeekHandler(c *gin.Context) {
	simulateWeek()
	league := leagueTable()
	c.JSON(http.StatusOK, gin.H{
		"currentWeek": currentWeek,
		"league":      league,
	})
}

func leagueTableHandler(c *gin.Context) {
	league := leagueTable()
	c.JSON(http.StatusOK, gin.H{
		"currentWeek": currentWeek,
		"league":      league,
	})
}

func playAllMatchesHandler(c *gin.Context) {

	if currentWeek == 6 {
		simulateWeek()
	}
	
	for currentWeek != 6 {
		simulateWeek()
	}

	league := leagueTable()
	c.JSON(http.StatusOK, gin.H{
		"currentWeek": currentWeek,
		"league":      league,
	})
}

func leagueTable() []Team {
	teams := getTeams()

	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Points == teams[j].Points {
			goalDifferenceI := teams[i].GoalsFor - teams[i].GoalsAgainst
			goalDifferenceJ := teams[j].GoalsFor - teams[j].GoalsAgainst
			return goalDifferenceI > goalDifferenceJ
		}
		return teams[i].Points > teams[j].Points
	})

	return teams
}

func simulateMatch(homeOStrength, homeDStrength, awayOStrength, awayDstrength int) (int, int) {
	homeGoals := max(0, homeOStrength/awayDstrength+rand.Intn(2))
	awayGoals := max(0, awayOStrength/homeDStrength-rand.Intn(2))
	return homeGoals, awayGoals
}

func simulateWeek() {
	if currentWeek == 0 {
		insertMatches()
	}

	if currentWeek == 6 {
		resetDatabase()
		insertMatches()
		currentWeek = 0
		return
	}

	matches := getMatches()
	slice := matches[currentWeek*2 : currentWeek*2+2]
	for _, match := range slice {
		homeTeam := getTeamByID(match.HomeTeamID)
		awayTeam := getTeamByID(match.AwayTeamID)
		homeGoals, awayGoals := simulateMatch(homeTeam.OffensiveStrength, homeTeam.DefensiveStrength, awayTeam.OffensiveStrength, awayTeam.DefensiveStrength)
		match.HomeGoals = homeGoals
		match.AwayGoals = awayGoals
		match.IsPlayed = true
		if homeGoals > awayGoals {
			homeTeam.Points += 3
			homeTeam.Wins++
			awayTeam.Loses++
		}
		if homeGoals == awayGoals {
			homeTeam.Points++
			awayTeam.Points++
			homeTeam.Draws++
			awayTeam.Draws++
		}
		if homeGoals < awayGoals {
			awayTeam.Points += 3
			awayTeam.Wins++
			homeTeam.Loses++
		}

		homeTeam.GoalsFor += homeGoals
		homeTeam.GoalsAgainst += awayGoals
		awayTeam.GoalsFor += awayGoals
		awayTeam.GoalsAgainst += homeGoals

		_, err := db.Exec("UPDATE teams SET points = ?, goals_for = ?, goals_against = ?, wins = ?, draws = ?, loses = ? WHERE id = ?", homeTeam.Points, homeTeam.GoalsFor, homeTeam.GoalsAgainst, homeTeam.Wins, homeTeam.Draws, homeTeam.Loses, homeTeam.ID)
		if err != nil {
			log.Fatalf("Error updating team: %s", err)
		}

		_, err = db.Exec("UPDATE teams SET points = ?, goals_for = ?, goals_against = ?, wins = ?, draws = ?, loses = ? WHERE id = ?", awayTeam.Points, awayTeam.GoalsFor, awayTeam.GoalsAgainst, awayTeam.Wins, awayTeam.Draws, awayTeam.Loses, awayTeam.ID)
		if err != nil {
			log.Fatalf("Error updating team: %s", err)
		}

		_, err = db.Exec("UPDATE matches SET home_goals = ?, away_goals = ?, is_played = ? WHERE id = ?", match.HomeGoals, match.AwayGoals, match.IsPlayed, match.ID)
		if err != nil {
			log.Fatalf("Error updating match: %s", err)
		}
	}
	currentWeek++
}
