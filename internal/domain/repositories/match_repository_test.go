package repositories_test

import (
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMatchRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Match{})
	assert.NoError(t, err)

	repo := repositories.NewMatchRepository(db)

	// Create
	match := &models.Match{HomeTeamID: 1, AwayTeamID: 2, HomeTeamScore: 2, AwayTeamScore: 1, Week: 1}
	err = repo.CreateMatch(match)
	assert.NoError(t, err)
	assert.NotZero(t, match.ID)

	// Read
	readMatch, err := repo.GetMatchByID(match.ID)
	assert.NoError(t, err)
	assert.Equal(t, match.HomeTeamID, readMatch.HomeTeamID)
	assert.Equal(t, match.AwayTeamID, readMatch.AwayTeamID)
	assert.Equal(t, match.HomeTeamScore, readMatch.HomeTeamScore)
	assert.Equal(t, match.AwayTeamScore, readMatch.AwayTeamScore)
	assert.Equal(t, match.Week, readMatch.Week)

	// Update
	readMatch.HomeTeamScore = 3
	err = repo.UpdateMatch(readMatch)
	assert.NoError(t, err)

	updatedMatch, err := repo.GetMatchByID(readMatch.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3, updatedMatch.HomeTeamScore)

	// Delete
	err = repo.DeleteMatch(match.ID)
	assert.NoError(t, err)

	_, err = repo.GetMatchByID(match.ID)
	assert.Error(t, err)
}
