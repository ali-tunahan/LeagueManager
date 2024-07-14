//go:build wireinject
// +build wireinject

package internal

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/repositories"
	"LeagueManager/internal/infrastructure/config"
	"LeagueManager/internal/presentation/controllers"
	"github.com/google/wire"
)

func Init() (*config.Initialization, error) {
	wire.Build(
		config.ConnectToDB,
		repositories.NewTeamRepository,
		repositories.NewLeagueRepository,
		repositories.NewStandingRepository,
		repositories.NewMatchRepository,
		services.NewTeamService,
		controllers.NewTeamController,
		services.NewLeagueService,
		controllers.NewLeagueController,
		config.NewInitialization,
	)
	return &config.Initialization{}, nil
}
