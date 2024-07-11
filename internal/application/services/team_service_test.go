package services

import (
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestTeamService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Team{})
	assert.NoError(t, err)

	repo := repositories.NewTeamRepository(db)
	service := NewTeamService(repo)

	// Create
	team := &models.Team{Name: "Team A", AttackStrength: 80, DefenseStrength: 70}
	err = service.CreateTeam(team)
	assert.NoError(t, err)
	assert.NotZero(t, team.ID)

	// Read
	readTeam, err := service.GetTeamByID(team.ID)
	assert.NoError(t, err)
	assert.Equal(t, team.Name, readTeam.Name)
	assert.Equal(t, team.AttackStrength, readTeam.AttackStrength)
	assert.Equal(t, team.DefenseStrength, readTeam.DefenseStrength)

	// Update
	readTeam.Name = "Team A Updated"
	err = service.UpdateTeam(readTeam)
	assert.NoError(t, err)

	updatedTeam, err := service.GetTeamByID(readTeam.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Team A Updated", updatedTeam.Name)

	// Delete
	err = service.DeleteTeam(team.ID)
	assert.NoError(t, err)

	_, err = service.GetTeamByID(team.ID)
	assert.Error(t, err)
}
