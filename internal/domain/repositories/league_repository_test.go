package repositories_test

import (
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLeagueRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.League{}, &models.Team{}, &models.Match{}, &models.Standing{})
	assert.NoError(t, err)

	repo := repositories.NewLeagueRepository(db)

	// Create
	league := &models.League{Name: "Premier League", CurrentWeek: 1}
	err = repo.CreateLeague(league)
	assert.NoError(t, err)
	assert.NotZero(t, league.ID)

	// Read
	readLeague, err := repo.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, league.Name, readLeague.Name)
	assert.Equal(t, league.CurrentWeek, readLeague.CurrentWeek)

	// Update
	readLeague.CurrentWeek = 2
	err = repo.UpdateLeague(readLeague)
	assert.NoError(t, err)

	updatedLeague, err := repo.GetLeagueByID(readLeague.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, updatedLeague.CurrentWeek)

	// Delete
	err = repo.DeleteLeague(league.ID)
	assert.NoError(t, err)

	_, err = repo.GetLeagueByID(league.ID)
	assert.Error(t, err)
}
