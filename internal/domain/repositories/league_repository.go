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
	GetLeaguesByTeamID(teamID uint) ([]*models.League, error)
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
	err := r.db.Preload("Teams").Find(&leagues).Error
	return leagues, err
}

func (r *LeagueRepositoryImpl) GetLeaguesByTeamID(teamID uint) ([]*models.League, error) {
	var leagues []*models.League

	// Joins the league_teams table to the leagues table and filters by the team ID
	err := r.db.Joins("JOIN league_teams ON league_teams.league_id = leagues.id").
		Where("league_teams.team_id = ?", teamID).
		Find(&leagues).Error
	return leagues, err
}
