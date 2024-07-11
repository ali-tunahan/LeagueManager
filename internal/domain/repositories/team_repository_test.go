package repositories

import (
	"LeagueManager/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestTeamRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Team{})
	assert.NoError(t, err)

	repo := NewTeamRepository(db)

	// Create
	team := models.Team{Name: "Team A", AttackStrength: 80, DefenseStrength: 70}
	err = repo.CreateTeam(&team)
	assert.NoError(t, err)
	assert.NotZero(t, team.ID)

	// Read
	readTeam, err := repo.GetTeamByID(team.ID)
	assert.NoError(t, err)
	assert.Equal(t, team.Name, readTeam.Name)
	assert.Equal(t, team.AttackStrength, readTeam.AttackStrength)
	assert.Equal(t, team.DefenseStrength, readTeam.DefenseStrength)

	// Update
	readTeam.Name = "Team A Updated"
	err = repo.UpdateTeam(readTeam)
	assert.NoError(t, err)

	updatedTeam, err := repo.GetTeamByID(readTeam.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Team A Updated", updatedTeam.Name)

	// Delete
	err = repo.DeleteTeam(team.ID)
	assert.NoError(t, err)

	_, err = repo.GetTeamByID(team.ID)
	assert.Error(t, err)
}
