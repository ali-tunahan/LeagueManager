// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/repositories"
	"LeagueManager/internal/infrastructure/config"
	"LeagueManager/internal/presentation/controllers"
)

// Injectors from injector.go:

func Init() (*config.Initialization, error) {
	db, err := config.ConnectToDB()
	if err != nil {
		return nil, err
	}
	teamRepository := repositories.NewTeamRepository(db)
	teamService := services.NewTeamService(teamRepository)
	teamController := controllers.NewTeamController(teamService)
	initialization := config.NewInitialization(teamRepository, teamService, teamController)
	return initialization, nil
}
