package config

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/repositories"
	"LeagueManager/internal/presentation/controllers"
)

type Initialization struct {
	TeamRepo repositories.TeamRepository
	TeamSvc  services.TeamService
	TeamCtrl *controllers.TeamController
}

func NewInitialization(
	teamRepo repositories.TeamRepository,
	teamSvc services.TeamService,
	teamCtrl *controllers.TeamController,
) *Initialization {
	return &Initialization{
		TeamRepo: teamRepo,
		TeamSvc:  teamSvc,
		TeamCtrl: teamCtrl,
	}
}
