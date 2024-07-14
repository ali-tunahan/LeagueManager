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
	RemoveTeamFromLeague(leagueID, teamID uint) error
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

	// Include all related entities when a single league is retrieved by ID
	err := r.db.Preload("Teams").Preload("Matches").Preload("Standings").First(&league, id).Error
	return league, err
}

func (r *LeagueRepositoryImpl) UpdateLeague(league *models.League) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Save(league).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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

func (r *LeagueRepositoryImpl) RemoveTeamFromLeague(leagueID, teamID uint) error {
	league := models.League{Model: gorm.Model{ID: leagueID}}
	team := models.Team{Model: gorm.Model{ID: teamID}}
	return r.db.Model(&league).Association("Teams").Delete(&team)
}
