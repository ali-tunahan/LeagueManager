package repositories_test

import (
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStandingRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Standing{})
	assert.NoError(t, err)

	repo := repositories.NewStandingRepository(db)

	// Create
	standing := &models.Standing{LeagueID: 1, TeamID: 1, Points: 3, Played: 1, Wins: 1, Draws: 0, Losses: 0, GoalDifference: 2}
	err = repo.CreateStanding(standing)
	assert.NoError(t, err)
	assert.NotZero(t, standing.ID)

	// Read
	readStanding, err := repo.GetStandingByID(standing.ID)
	assert.NoError(t, err)
	assert.Equal(t, standing.LeagueID, readStanding.LeagueID)
	assert.Equal(t, standing.TeamID, readStanding.TeamID)
	assert.Equal(t, standing.Points, readStanding.Points)
	assert.Equal(t, standing.Played, readStanding.Played)
	assert.Equal(t, standing.Wins, readStanding.Wins)
	assert.Equal(t, standing.Draws, readStanding.Draws)
	assert.Equal(t, standing.Losses, readStanding.Losses)
	assert.Equal(t, standing.GoalDifference, readStanding.GoalDifference)

	// Update
	readStanding.Points = 4
	err = repo.UpdateStanding(readStanding)
	assert.NoError(t, err)

	updatedStanding, err := repo.GetStandingByID(readStanding.ID)
	assert.NoError(t, err)
	assert.Equal(t, 4, updatedStanding.Points)

	// Delete
	err = repo.DeleteStanding(standing.ID)
	assert.NoError(t, err)

	_, err = repo.GetStandingByID(standing.ID)
	assert.Error(t, err)

	// Create multiple standings for GetStandingByTeam test
	standing1 := &models.Standing{LeagueID: 2, TeamID: 1, Points: 5, Played: 2, Wins: 1, Draws: 2, Losses: 0, GoalDifference: 3}
	standing2 := &models.Standing{LeagueID: 2, TeamID: 2, Points: 4, Played: 2, Wins: 1, Draws: 1, Losses: 0, GoalDifference: 2}
	err = repo.CreateStanding(standing1)
	assert.NoError(t, err)
	err = repo.CreateStanding(standing2)
	assert.NoError(t, err)

	// Test GetStandingByTeam
	standingByTeam, err := repo.GetStandingByTeam(2, 1)
	assert.NoError(t, err)
	assert.Equal(t, standing1.LeagueID, standingByTeam.LeagueID)
	assert.Equal(t, standing1.TeamID, standingByTeam.TeamID)
	assert.Equal(t, standing1.Points, standingByTeam.Points)
	assert.Equal(t, standing1.Played, standingByTeam.Played)
	assert.Equal(t, standing1.Wins, standingByTeam.Wins)
	assert.Equal(t, standing1.Draws, standingByTeam.Draws)
	assert.Equal(t, standing1.Losses, standingByTeam.Losses)
	assert.Equal(t, standing1.GoalDifference, standingByTeam.GoalDifference)
}
