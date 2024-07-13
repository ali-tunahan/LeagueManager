package repositories

import (
	"LeagueManager/internal/domain/models"
	"gorm.io/gorm"
)

type StandingRepository interface {
	CreateStanding(standing *models.Standing) error
	GetStandingByID(id uint) (*models.Standing, error)
	UpdateStanding(standing *models.Standing) error
	DeleteStanding(id uint) error
	GetAllStandings() ([]*models.Standing, error)
	GetStandingByTeam(leagueID uint, teamID uint) (*models.Standing, error)
}

type StandingRepositoryImpl struct {
	db *gorm.DB
}

func NewStandingRepository(db *gorm.DB) StandingRepository {
	return &StandingRepositoryImpl{db: db}
}

func (r *StandingRepositoryImpl) CreateStanding(standing *models.Standing) error {
	return r.db.Create(&standing).Error
}

func (r *StandingRepositoryImpl) GetStandingByID(id uint) (*models.Standing, error) {
	var standing *models.Standing
	err := r.db.First(&standing, id).Error
	return standing, err
}

func (r *StandingRepositoryImpl) UpdateStanding(standing *models.Standing) error {
	return r.db.Save(&standing).Error
}

func (r *StandingRepositoryImpl) DeleteStanding(id uint) error {
	return r.db.Delete(&models.Standing{}, id).Error
}

func (r *StandingRepositoryImpl) GetAllStandings() ([]*models.Standing, error) {
	var standings []*models.Standing
	err := r.db.Find(&standings).Error
	return standings, err
}

func (r *StandingRepositoryImpl) GetStandingByTeam(leagueID uint, teamID uint) (*models.Standing, error) {
	var standing *models.Standing

	// query standings with leagueID and teamID matching the requested one
	err := r.db.Where("league_id = ? AND team_id = ?", leagueID, teamID).
		First(&standing).Error

	return standing, err
}
