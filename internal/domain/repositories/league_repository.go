package repositories

import (
	"LeagueManager/internal/domain/models"
	"gorm.io/gorm"
)

type LeagueRepository interface {
	CreateLeague(league *models.League) error
	GetLeagueByID(id uint) (*models.League, error)
	UpdateLeague(league *models.League) error
	DeleteLeague(id uint) error
	GetAllLeagues() ([]*models.League, error)
}

type LeagueRepositoryImpl struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) LeagueRepository {
	return &LeagueRepositoryImpl{db: db}
}

func (r *LeagueRepositoryImpl) CreateLeague(league *models.League) error {
	return r.db.Create(&league).Error
}

func (r *LeagueRepositoryImpl) GetLeagueByID(id uint) (*models.League, error) {
	var league *models.League
	err := r.db.Preload("Teams").First(&league, id).Error
	return league, err
}

func (r *LeagueRepositoryImpl) UpdateLeague(league *models.League) error {
	return r.db.Save(&league).Error
}

func (r *LeagueRepositoryImpl) DeleteLeague(id uint) error {
	return r.db.Delete(&models.League{}, id).Error
}

func (r *LeagueRepositoryImpl) GetAllLeagues() ([]*models.League, error) {
	var leagues []*models.League
	err := r.db.Preload("Teams").Preload("Matches").Preload("Standing").Find(&leagues).Error
	return leagues, err
}
