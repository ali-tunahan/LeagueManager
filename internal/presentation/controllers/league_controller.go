package controllers

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LeagueController struct {
	leagueService services.LeagueService
	teamService   services.TeamService
}

func NewLeagueController(leagueService services.LeagueService, teamService services.TeamService) *LeagueController {
	return &LeagueController{
		leagueService: leagueService,
		teamService:   teamService,
	}
}

// CreateAndAdvanceLeague creates a league, adds 4 teams, and advances 5 weeks
func (lc *LeagueController) CreateAndAdvanceLeague(c *gin.Context) {
	// Create 4 teams
	teams := []models.Team{
		{Name: "Team A", AttackStrength: 80, DefenseStrength: 75},
		{Name: "Team B", AttackStrength: 70, DefenseStrength: 80},
		{Name: "Team C", AttackStrength: 65, DefenseStrength: 70},
		{Name: "Team D", AttackStrength: 60, DefenseStrength: 65},
	}

	for _, team := range teams {
		err := lc.teamService.CreateTeam(&team)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teams"})
			return
		}
	}

	// Retrieve the created teams
	allTeams, err := lc.teamService.GetAllTeams()
	if err != nil || len(allTeams) <= 4 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created teams", "teams": allTeams, "err": err})
		return
	}

	// Create a league
	league := &models.League{
		Name:        "Debug League",
		CurrentWeek: 0,
	}

	err = lc.leagueService.CreateLeague(league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint("Failed to create league err is: ", err)})
		return
	}

	count := 0
	// Add the teams to the league
	for _, team := range allTeams {
		err = lc.leagueService.AddTeamToLeague(league.ID, team.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add team to league", "err": err})
			return
		}
		count++
		if count == 4 {
			break
		}
	}

	err = lc.leagueService.StartLeague(league.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start league: " + err.Error()})
		return
	}

	// Advance the league by 5 weeks
	for i := 0; i < 5; i++ {
		err = lc.leagueService.AdvanceWeek(league.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to advance week: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "League created and advanced 5 weeks", "league_id": league.ID})
}

// GetLeagueByID retrieves a league by its ID
func (lc *LeagueController) GetLeagueByID(c *gin.Context) {
	leagueID := c.Param("leagueID")

	// cast to uint
	i, _ := strconv.Atoi(leagueID)
	leagueIDUint := uint(i)
	league, err := lc.leagueService.GetLeagueByID(leagueIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve league", "err": err})
		return
	}

	c.JSON(http.StatusOK, league)
}
