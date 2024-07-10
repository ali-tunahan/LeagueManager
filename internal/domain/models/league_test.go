package models

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestLeagueModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&League{}, &Team{})
	assert.NoError(t, err)

	// Create
	league := &League{Name: "Premier League", CurrentWeek: 1}
	err = db.Create(league).Error
	assert.NoError(t, err)
	assert.NotZero(t, league.ID)

	// Read
	var readLeague League
	err = db.First(&readLeague, league.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, league.Name, readLeague.Name)
	assert.Equal(t, league.CurrentWeek, readLeague.CurrentWeek)

	// Update
	readLeague.CurrentWeek = 2
	err = db.Save(&readLeague).Error
	assert.NoError(t, err)

	var updatedLeague League
	err = db.First(&updatedLeague, league.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, updatedLeague.CurrentWeek)

	// Delete
	err = db.Delete(&League{}, league.ID).Error
	assert.NoError(t, err)

	var deletedLeague League
	err = db.First(&deletedLeague, league.ID).Error
	assert.Error(t, err)
}
