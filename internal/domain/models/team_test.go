package models

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestTeamModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&Team{})
	assert.NoError(t, err)

	// Create
	team := &Team{Name: "Team A", AttackStrength: 5, DefenseStrength: 3}
	err = db.Create(team).Error
	assert.NoError(t, err)
	assert.NotZero(t, team.ID)

	// Read
	var readTeam Team
	err = db.First(&readTeam, team.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, team.Name, readTeam.Name)
	assert.Equal(t, team.AttackStrength, readTeam.AttackStrength)
	assert.Equal(t, team.DefenseStrength, readTeam.DefenseStrength)

	// Update
	readTeam.AttackStrength = 6
	err = db.Save(&readTeam).Error
	assert.NoError(t, err)

	var updatedTeam Team
	err = db.First(&updatedTeam, team.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, 6, updatedTeam.AttackStrength)

	// Delete
	err = db.Delete(&Team{}, team.ID).Error
	assert.NoError(t, err)

	var deletedTeam Team
	err = db.First(&deletedTeam, team.ID).Error
	assert.Error(t, err)
}
