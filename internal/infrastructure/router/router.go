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
		league.POST("/create", init.LeagueCtrl.CreateLeague)
		league.POST("/initialize", init.LeagueCtrl.CreateAndInitializeLeague)
		league.POST("/add-team/:leagueID/:teamID", init.LeagueCtrl.AddTeamToLeague)
		league.POST("/remove-team/:leagueID/:teamID", init.LeagueCtrl.RemoveTeamFromLeague)
		league.POST("/advance-week/:leagueID", init.LeagueCtrl.AdvanceWeek)
		league.GET("/view-matches/:leagueID", init.LeagueCtrl.ViewMatchResults)
		league.POST("/edit-match/:matchID", init.LeagueCtrl.EditMatchResults)
		league.GET("/predict-champion/:leagueID", init.LeagueCtrl.PredictChampion)
		league.POST("/play-all-matches/:leagueID", init.LeagueCtrl.PlayAllMatches)
	}

	return router
}
