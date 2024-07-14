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

	// Add the LeagueService and LeagueController fields
	LeagueSvc  services.LeagueService
	LeagueCtrl *controllers.LeagueController
}

func NewInitialization(
	teamRepo repositories.TeamRepository,
	leagueRepo repositories.LeagueRepository,
	standingRepo repositories.StandingRepository,
	matchRepo repositories.MatchRepository,
	teamSvc services.TeamService,
	teamCtrl *controllers.TeamController,
	leagueSvc services.LeagueService,
	leagueCtrl *controllers.LeagueController,
) *Initialization {
	return &Initialization{
		TeamRepo:     teamRepo,
		LeagueRepo:   leagueRepo,
		StandingRepo: standingRepo,
		MatchRepo:    matchRepo,
		TeamSvc:      teamSvc,
		TeamCtrl:     teamCtrl,
		LeagueSvc:    leagueSvc,
		LeagueCtrl:   leagueCtrl,
	}
}
