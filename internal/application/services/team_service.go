package services

import (
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"errors"
)

type TeamService interface {
	CreateTeam(team *models.Team) error
	GetTeamByID(id uint) (*models.Team, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(id uint) error
	GetAllTeams() ([]*models.Team, error)
}

type TeamServiceImpl struct {
	teamRepo   repositories.TeamRepository
	leagueRepo repositories.LeagueRepository
}

func NewTeamService(teamRepo repositories.TeamRepository, leagueRepo repositories.LeagueRepository) TeamService {
	return &TeamServiceImpl{teamRepo: teamRepo, leagueRepo: leagueRepo}
}

func (s *TeamServiceImpl) CreateTeam(team *models.Team) error {
	return s.teamRepo.CreateTeam(team)
}

func (s *TeamServiceImpl) GetTeamByID(id uint) (*models.Team, error) {
	return s.teamRepo.GetTeamByID(id)
}

func (s *TeamServiceImpl) UpdateTeam(team *models.Team) error {
	return s.teamRepo.UpdateTeam(team)
}

func (s *TeamServiceImpl) DeleteTeam(id uint) error {
	leagues, err := s.leagueRepo.GetLeaguesByTeamID(id)
	if err != nil {
		return err
	}

	for _, league := range leagues {
		if league.IsActive() {
			return errors.New("cannot delete team that is part of an active league")
		}
	}

	return s.teamRepo.DeleteTeam(id)
}

func (s *TeamServiceImpl) GetAllTeams() ([]*models.Team, error) {
	return s.teamRepo.GetAllTeams()
}

func (s *TeamServiceImpl) GetTeamsByLeague(leagueID uint) ([]models.Team, error) {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return nil, err
	}

	return league.Teams, nil
}
