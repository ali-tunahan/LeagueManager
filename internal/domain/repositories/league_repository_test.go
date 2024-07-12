package repositories

import (
	"LeagueManager/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestLeagueRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.League{}, &models.Team{})
	assert.NoError(t, err)

	repo := NewLeagueRepository(db)

	// Create league
	league := &models.League{Name: "Premier League", CurrentWeek: 0}
	err = repo.CreateLeague(league)
	assert.NoError(t, err)
	assert.NotZero(t, league.ID)

	// Add teams to league
	teamA := &models.Team{Name: "Team A", AttackStrength: 80, DefenseStrength: 70}
	teamB := &models.Team{Name: "Team B", AttackStrength: 75, DefenseStrength: 65}
	err = db.Create(&teamA).Error
	assert.NoError(t, err)
	err = db.Create(&teamB).Error
	assert.NoError(t, err)

	league.Teams = []models.Team{*teamA, *teamB}
	err = repo.UpdateLeague(league)
	assert.NoError(t, err)

	// Get leagues by team ID
	leagues, err := repo.GetLeaguesByTeamID(teamA.ID)
	assert.NoError(t, err)
	assert.Len(t, leagues, 1)
	assert.Equal(t, "Premier League", leagues[0].Name)

	// Get league by ID
	readLeague, err := repo.GetLeagueByID(league.ID)
	assert.NoError(t, err)
	assert.Equal(t, league.Name, readLeague.Name)

	// Update league
	readLeague.CurrentWeek = 1
	err = repo.UpdateLeague(readLeague)
	assert.NoError(t, err)

	// Delete league
	err = repo.DeleteLeague(league.ID)
	assert.NoError(t, err)
}
