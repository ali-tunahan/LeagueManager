package models

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestStandingsModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Standings{}, &League{}, &Team{})
	assert.NoError(t, err)

	// Create Team
	team := &Team{Name: "Team A", AttackStrength: 5, DefenseStrength: 3}
	db.Create(team)

	// Create League
	league := &League{Name: "Premier League", CurrentWeek: 1}
	db.Create(league)

	// Create Standings
	standings := &Standings{
		LeagueID:       league.ID,
		TeamID:         team.ID,
		Points:         10,
		Played:         4,
		Wins:           3,
		Draws:          1,
		Losses:         0,
		GoalDifference: 5,
	}
	err = db.Create(standings).Error
	assert.NoError(t, err)
	assert.NotZero(t, standings.ID)

	// Read
	var readStandings Standings
	err = db.First(&readStandings, standings.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, standings.Points, readStandings.Points)
	assert.Equal(t, standings.Played, readStandings.Played)

	// Update
	readStandings.Points = 12
	err = db.Save(&readStandings).Error
	assert.NoError(t, err)

	var updatedStandings Standings
	err = db.First(&updatedStandings, standings.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 12, updatedStandings.Points)

	// Delete
	err = db.Delete(&Standings{}, standings.ID).Error
	assert.NoError(t, err)

	var deletedStandings Standings
	err = db.First(&deletedStandings, standings.ID).Error
	assert.Error(t, err)
}
