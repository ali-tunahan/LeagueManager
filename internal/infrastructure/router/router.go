package router

import (
	"LeagueManager/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

func Init(init *config.Initialization) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		team := api.Group("/teams")
		team.GET("", init.TeamCtrl.GetAllTeams)
		team.POST("", init.TeamCtrl.AddTeam)
		team.GET("/:teamID", init.TeamCtrl.GetTeamByID)
		team.PUT("/:teamID", init.TeamCtrl.UpdateTeam)
		team.DELETE("/:teamID", init.TeamCtrl.DeleteTeam)

		// Add the league routes
		league := api.Group("/leagues")
		league.POST("/debug", init.LeagueCtrl.CreateAndAdvanceLeague)
		league.GET("/:leagueID", init.LeagueCtrl.GetLeagueByID)
	}

	return router
}
