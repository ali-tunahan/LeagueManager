package repositories

import (
	"LeagueManager/internal/domain/models"
	"gorm.io/gorm"
)

type MatchRepository interface {
	CreateMatch(match *models.Match) error
	GetMatchByID(id uint) (*models.Match, error)
	UpdateMatch(match *models.Match) error
	DeleteMatch(id uint) error
	GetAllMatches() ([]*models.Match, error)
}

type MatchRepositoryImpl struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &MatchRepositoryImpl{db: db}
}

func (r *MatchRepositoryImpl) CreateMatch(match *models.Match) error {
	return r.db.Create(&match).Error
}

func (r *MatchRepositoryImpl) GetMatchByID(id uint) (*models.Match, error) {
	var match *models.Match
	err := r.db.First(&match, id).Error
	return match, err
}

func (r *MatchRepositoryImpl) UpdateMatch(match *models.Match) error {
	return r.db.Save(&match).Error
}

func (r *MatchRepositoryImpl) DeleteMatch(id uint) error {
	return r.db.Delete(&models.Match{}, id).Error
}

func (r *MatchRepositoryImpl) GetAllMatches() ([]*models.Match, error) {
	var matches []*models.Match
	err := r.db.Find(&matches).Error
	return matches, err
}
