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
}
