package services_test

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupLeagueServiceTest() (services.LeagueService, services.TeamService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.Team{}, &models.League{}, &models.Match{}, &models.Standing{})

	teamRepo := repositories.NewTeamRepository(db)
	leagueRepo := repositories.NewLeagueRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	standingRepo := repositories.NewStandingRepository(db)

	leagueService := services.NewLeagueService(leagueRepo, teamRepo, matchRepo, standingRepo)
	teamService := services.NewTeamService(teamRepo, leagueRepo)

	return leagueService, teamService
}

func createTestTeamsForLeague(teamService services.TeamService) {
	teams := []models.Team{
		{Name: "Team A", AttackStrength: 80, DefenseStrength: 75},
		{Name: "Team B", AttackStrength: 70, DefenseStrength: 80},
		{Name: "Team C", AttackStrength: 65, DefenseStrength: 70},
		{Name: "Team D", AttackStrength: 60, DefenseStrength: 65},
	}

	for _, team := range teams {
		teamService.CreateTeam(&team)
	}
}

func createTestLeagueForService(leagueService services.LeagueService, teamService services.TeamService) *models.League {
	teams, _ := teamService.GetAllTeams()

	var teamStructs []models.Team
	// Cast each team to a Team struct
	for i, team := range teams {
		teamStructs[i] = *team
	}

	league := &models.League{
		Name:  "Test League",
		Teams: teamStructs,
	}

	err := leagueService.CreateLeague(league)
	if err != nil {
		return nil
	}
	return league
}

func TestCreateLeague(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	teams, err := teamService.GetAllTeams()
	assert.NoError(t, err)

	var teamStructs []models.Team
	// Cast each team to a Team struct
	for i, team := range teams {
		teamStructs[i] = *team
	}

	league := &models.League{
		Name:  "New League",
		Teams: teamStructs,
	}

	err = leagueService.CreateLeague(league)
	assert.NoError(t, err)
	assert.NotZero(t, league.ID)

	createdLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, league.Name, createdLeague.Name)
	assert.Equal(t, len(teams), len(createdLeague.Teams))
}

func TestAdvanceWeek(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)

	err := leagueService.AdvanceWeek(league.ID)
	assert.NoError(t, err)

	updatedLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedLeague.CurrentWeek)

	err = leagueService.AdvanceWeek(999) // Non-existent league
	assert.Error(t, err)
}

func TestPlayAllMatches(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)

	err := leagueService.PlayAllMatches(league.ID)
	assert.NoError(t, err)

	updatedLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 38, updatedLeague.CurrentWeek)

	matches, err := leagueService.ViewMatchResults(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 38*2, len(matches)) // Assuming 2 matches per week in a 4 team league
}

func TestEditMatchResults(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)

	leagueService.AdvanceWeek(league.ID)
	matches, err := leagueService.ViewMatchResults(league.ID)
	assert.NoError(t, err)
	match := matches[0]

	updatedMatch := &models.Match{
		HomeTeamScore: 2,
		AwayTeamScore: 2,
		LeagueID:      match.LeagueID,
		HomeTeamID:    match.HomeTeamID,
		AwayTeamID:    match.AwayTeamID,
		Week:          match.Week,
	}

	err = leagueService.EditMatchResults(match.ID, updatedMatch)
	assert.NoError(t, err)

	newLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)

	standings := newLeague.Standings
	for _, standing := range standings {
		if standing.TeamID == match.HomeTeamID || standing.TeamID == match.AwayTeamID {
			assert.Equal(t, 1, standing.Played)
			assert.Equal(t, 1, standing.Draws)
			assert.Equal(t, 1, standing.Points)
			assert.Equal(t, 0, standing.Wins)
			assert.Equal(t, 0, standing.Losses)
		}
	}

	err = leagueService.EditMatchResults(0, &models.Match{}) // Non-existent match
	assert.Error(t, err)
}

func TestPredictChampion(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)

	for i := 0; i < 4; i++ {
		leagueService.AdvanceWeek(league.ID)
	}

	predictions, err := leagueService.PredictChampion(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(predictions))

	totalProbability := 0.0
	for _, prediction := range predictions {
		totalProbability += prediction.WinProbability
	}
	assert.InEpsilon(t, 1.0, totalProbability, 0.01)
}

func TestGetLeagueStandings(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)

	for i := 0; i < 5; i++ {
		err := leagueService.AdvanceWeek(league.ID)
		assert.NoError(t, err)
	}

	newLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(newLeague.Standings))
}

func TestAddTeamToLeague(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)

	newTeam := models.Team{Name: "Team E", AttackStrength: 55, DefenseStrength: 60}
	err := teamService.CreateTeam(&newTeam)
	assert.NoError(t, err)

	err = leagueService.AddTeamToLeague(league.ID, newTeam.ID)
	assert.Error(t, err) // Should fail as league already has 4 teams

	league, err = leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(league.Teams))
}

func TestRemoveTeamFromLeague(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)

	teamToRemove := league.Teams[0]
	err := leagueService.RemoveTeamFromLeague(league.ID, teamToRemove.ID)
	assert.NoError(t, err)

	league, err = leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(league.Teams))
}

func TestDeleteTeam(t *testing.T) {
	leagueService, teamService := setupLeagueServiceTest()

	createTestTeamsForLeague(teamService)
	league := createTestLeagueForService(leagueService, teamService)
	teams, _ := teamService.GetAllTeams()

	leagueService.AddTeamToLeague(league.ID, teams[0].ID)

	err := teamService.DeleteTeam(teams[0].ID)
	assert.Error(t, err) // Should fail as team is part of active league

	for i := 1; i < len(teams); i++ {
		teamService.DeleteTeam(teams[i].ID)
	}

	remainingTeams, err := teamService.GetAllTeams()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(remainingTeams))
	assert.Equal(t, teams[0].ID, remainingTeams[0].ID)
}
