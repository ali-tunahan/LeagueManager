package models

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestMatchModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Match{}, &League{}, &Team{})
	assert.NoError(t, err)

	// Create Teams
	homeTeam := &Team{Name: "Team A", AttackStrength: 5, DefenseStrength: 3}
	awayTeam := &Team{Name: "Team B", AttackStrength: 4, DefenseStrength: 4}
	db.Create(homeTeam)
	db.Create(awayTeam)

	// Create League
	league := &League{Name: "Premier League", CurrentWeek: 1}
	db.Create(league)

	// Create Match
	match := &Match{
		LeagueID:      league.ID,
		HomeTeamID:    homeTeam.ID,
		AwayTeamID:    awayTeam.ID,
		HomeTeamScore: 2,
		AwayTeamScore: 1,
		Week:          1,
	}
	err = db.Create(match).Error
	assert.NoError(t, err)
	assert.NotZero(t, match.ID)

	// Read
	var readMatch Match
	err = db.First(&readMatch, match.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, match.HomeTeamScore, readMatch.HomeTeamScore)
	assert.Equal(t, match.AwayTeamScore, readMatch.AwayTeamScore)

	// Update
	readMatch.HomeTeamScore = 3
	err = db.Save(&readMatch).Error
	assert.NoError(t, err)

	var updatedMatch Match
	err = db.First(&updatedMatch, match.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 3, updatedMatch.HomeTeamScore)

	// Delete
	err = db.Delete(&Match{}, match.ID).Error
	assert.NoError(t, err)

	var deletedMatch Match
	err = db.First(&deletedMatch, match.ID).Error
	assert.Error(t, err)
}
