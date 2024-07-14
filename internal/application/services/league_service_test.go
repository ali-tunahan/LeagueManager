package services_test

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupLeagueServiceTest() (*gorm.DB, services.LeagueService, services.TeamService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	err = db.AutoMigrate(&models.Team{}, &models.League{}, &models.Match{}, &models.Standing{})
	if err != nil {
		panic("failed to connect to migrate database")
	}

	teamRepo := repositories.NewTeamRepository(db)
	leagueRepo := repositories.NewLeagueRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	standingRepo := repositories.NewStandingRepository(db)

	leagueService := services.NewLeagueService(leagueRepo, teamRepo, matchRepo, standingRepo)
	teamService := services.NewTeamService(teamRepo, leagueRepo)

	return db, leagueService, teamService
}

func createTestTeamsForLeague(teamService services.TeamService) {
	teams := []models.Team{
		{Name: "Team A", AttackStrength: 80, DefenseStrength: 75},
		{Name: "Team B", AttackStrength: 70, DefenseStrength: 80},
		{Name: "Team C", AttackStrength: 65, DefenseStrength: 70},
		{Name: "Team D", AttackStrength: 60, DefenseStrength: 65},
	}

	for _, team := range teams {
		err := teamService.CreateTeam(&team)
		if err != nil {
			panic("failed to create test teams")
		}
	}
}

func createTestLeagueForService(leagueService services.LeagueService, teamService services.TeamService) *models.League {
	createTestTeamsForLeague(teamService) // Ensure teams are created first

	teams, err := teamService.GetAllTeams()
	if err != nil {
		panic(fmt.Sprint("failed to retrieve test teams", " error is ", err, " retrieved length is ", len(teams)))
	}

	// if not 4 teams panic
	if len(teams) != 4 {
		panic("more than 4 teams created")
	}

	var teamStructs []models.Team
	// Cast each team to a Team struct
	for _, team := range teams {
		teamStructs = append(teamStructs, *team)
	}

	league := &models.League{
		Name:  "Test League",
		Teams: teamStructs,
	}

	err = leagueService.CreateLeague(league)
	if err != nil {
		panic("failed to create test league")
	}
	return league
}

func TestCreateLeague(t *testing.T) {
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	createTestTeamsForLeague(teamService)
	teams, err := teamService.GetAllTeams()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(teams))

	var teamStructs []models.Team
	// Cast each team to a Team struct
	for _, team := range teams {
		teamStructs = append(teamStructs, *team)
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
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	league := createTestLeagueForService(leagueService, teamService)

	err := leagueService.StartLeague(league.ID)
	assert.NoError(t, err)

	updatedLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedLeague.CurrentWeek)

	err = leagueService.AdvanceWeek(league.ID)
	assert.NoError(t, err)

	updatedLeague, err = leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, updatedLeague.CurrentWeek)

	err = leagueService.AdvanceWeek(999) // Non-existent league
	assert.Error(t, err)
}

func TestPlayAllMatches(t *testing.T) {
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	league := createTestLeagueForService(leagueService, teamService)

	err := leagueService.StartLeague(league.ID)
	assert.NoError(t, err)

	err = leagueService.PlayAllMatches(league.ID)
	assert.NoError(t, err)

	updatedLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 38, updatedLeague.CurrentWeek)

	matches, err := leagueService.ViewMatchResults(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(matches)) // Assuming 2 matches per week
}

func TestEditMatchResults(t *testing.T) {
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	league := createTestLeagueForService(leagueService, teamService)

	err := leagueService.StartLeague(league.ID)
	assert.NoError(t, err)

	err = leagueService.AdvanceWeek(league.ID)
	assert.NoError(t, err)

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
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	league := createTestLeagueForService(leagueService, teamService)

	err := leagueService.StartLeague(league.ID)
	assert.NoError(t, err)

	for i := 0; i < 4; i++ {
		err := leagueService.AdvanceWeek(league.ID)
		assert.NoError(t, err)
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
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	league := createTestLeagueForService(leagueService, teamService)

	err := leagueService.StartLeague(league.ID)
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		err := leagueService.AdvanceWeek(league.ID)
		assert.NoError(t, err)
	}

	newLeague, err := leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(newLeague.Standings))
}

func TestAddTeamToLeague(t *testing.T) {
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

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
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	league := createTestLeagueForService(leagueService, teamService)
	assert.Equal(t, 4, len(league.Teams), "Initial team count should be 4")

	teamToRemove := league.Teams[0]
	err := leagueService.RemoveTeamFromLeague(league.ID, teamToRemove.ID)
	assert.NoError(t, err)

	league, err = leagueService.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	t.Logf("Remaining teams: %v", league.Teams)
	assert.Equal(t, 3, len(league.Teams))
}

func TestDeleteTeam(t *testing.T) {
	db, leagueService, teamService := setupLeagueServiceTest()

	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("failed to close database connection")
		}
	}(sqlDB)

	league := createTestLeagueForService(leagueService, teamService)

	teams := league.Teams

	assert.Equal(t, 4, len(teams))

	err := leagueService.StartLeague(league.ID)
	assert.NoError(t, err)

	err = teamService.DeleteTeam(teams[0].ID)
	assert.Error(t, err) // Should fail as team is part of active league

	for i := 1; i < len(teams); i++ {
		err = teamService.DeleteTeam(teams[i].ID)
		assert.Error(t, err)
	}

	err = teamService.CreateTeam(&models.Team{Name: "Team X", AttackStrength: 55, DefenseStrength: 60})
	assert.NoError(t, err)

	allTeams, _ := teamService.GetAllTeams()

	for i := 0; i < len(allTeams); i++ {
		err = teamService.DeleteTeam(allTeams[i].ID)
		if allTeams[i].ID == uint(5) {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}

	remainingTeams, err := teamService.GetAllTeams()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(remainingTeams))
	if len(remainingTeams) == 0 {
		t.Error("No teams left")
		return
	}
}
