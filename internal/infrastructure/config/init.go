package config

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/repositories"
	"LeagueManager/internal/presentation/controllers"
)

type Initialization struct {
	TeamRepo     repositories.TeamRepository
	LeagueRepo   repositories.LeagueRepository
	StandingRepo repositories.StandingRepository
	MatchRepo    repositories.MatchRepository

	TeamSvc  services.TeamService
	TeamCtrl *controllers.TeamController
}

func NewInitialization(
	teamRepo repositories.TeamRepository,
	leagueRepo repositories.LeagueRepository,
	standingRepo repositories.StandingRepository,
	matchRepo repositories.MatchRepository,
	teamSvc services.TeamService,
	teamCtrl *controllers.TeamController,
) *Initialization {
	return &Initialization{
		TeamRepo:     teamRepo,
		LeagueRepo:   leagueRepo,
		StandingRepo: standingRepo,
		MatchRepo:    matchRepo,
		TeamSvc:      teamSvc,
		TeamCtrl:     teamCtrl,
	}
}
