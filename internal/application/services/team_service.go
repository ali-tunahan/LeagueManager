package services

import (
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
)

type TeamService interface {
	CreateTeam(team *models.Team) error
	GetTeamByID(id uint) (*models.Team, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(id uint) error
	GetAllTeams() ([]*models.Team, error)
}

type TeamServiceImpl struct {
	repo repositories.TeamRepository
}

func NewTeamService(repo repositories.TeamRepository) TeamService {
	return &TeamServiceImpl{repo: repo}
}

func (s *TeamServiceImpl) CreateTeam(team *models.Team) error {
	return s.repo.CreateTeam(team)
}

func (s *TeamServiceImpl) GetTeamByID(id uint) (*models.Team, error) {
	return s.repo.GetTeamByID(id)
}

func (s *TeamServiceImpl) UpdateTeam(team *models.Team) error {
	return s.repo.UpdateTeam(team)
}

func (s *TeamServiceImpl) DeleteTeam(id uint) error {
	return s.repo.DeleteTeam(id)
}

func (s *TeamServiceImpl) GetAllTeams() ([]*models.Team, error) {
	return s.repo.GetAllTeams()
}
