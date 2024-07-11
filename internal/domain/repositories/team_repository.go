package repositories

import (
	"LeagueManager/internal/domain/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(team models.Team) error
	GetTeamByID(id uint) (models.Team, error)
	UpdateTeam(team models.Team) error
	DeleteTeam(id uint) error
	GetAllTeams() ([]models.Team, error)
}

type TeamRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &TeamRepositoryImpl{db: db}
}

func (r *TeamRepositoryImpl) CreateTeam(team models.Team) error {
	return r.db.Create(&team).Error
}

func (r *TeamRepositoryImpl) GetTeamByID(id uint) (models.Team, error) {
	var team models.Team
	err := r.db.First(&team, id).Error
	return team, err
}

func (r *TeamRepositoryImpl) UpdateTeam(team models.Team) error {
	return r.db.Save(&team).Error
}

func (r *TeamRepositoryImpl) DeleteTeam(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}

func (r *TeamRepositoryImpl) GetAllTeams() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Find(&teams).Error
	return teams, err
}
