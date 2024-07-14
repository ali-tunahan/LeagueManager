package controllers

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"

	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// LeagueController handles league-related requests
type LeagueController struct {
	leagueService services.LeagueService
	teamService   services.TeamService
}

// NewLeagueController creates a new LeagueController
func NewLeagueController(leagueService services.LeagueService, teamService services.TeamService) *LeagueController {
	return &LeagueController{
		leagueService: leagueService,
		teamService:   teamService,
	}
}

// CreateLeague creates a league with no teams
// @Summary Create a league with no teams
// @Tags League
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /league/create [post]
func (lc *LeagueController) CreateLeague(c *gin.Context) {
	var league models.League
	if err := c.ShouldBindJSON(&league); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := lc.leagueService.CreateLeague(&league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create league"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "League created successfully", "league_id": league.ID})
}

// CreateAndInitializeLeague creates a league and initializes it with Premier League teams
// @Summary Create a league and initialize it with Premier League teams
// @Tags League
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /league/initialize [post]
func (lc *LeagueController) CreateAndInitializeLeague(c *gin.Context) {
	teams := []models.Team{
		{Name: "Arsenal", AttackStrength: 85, DefenseStrength: 80},
		{Name: "Chelsea", AttackStrength: 82, DefenseStrength: 78},
		{Name: "Liverpool", AttackStrength: 90, DefenseStrength: 85},
		{Name: "Manchester City", AttackStrength: 92, DefenseStrength: 88},
	}

	for _, team := range teams {
		err := lc.teamService.CreateTeam(&team)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teams"})
			return
		}
	}

	allTeams, err := lc.teamService.GetAllTeams()
	if err != nil || len(allTeams) < 4 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created teams"})
		return
	}

	// inside allTeam search for Arsenal, Chelsea, Liverpool, Manchester City and add them to the league
	league := &models.League{Name: "Premier League"}
	for _, team := range allTeams {
		if team.Name == "Arsenal" || team.Name == "Chelsea" || team.Name == "Liverpool" || team.Name == "Manchester City" {
			league.Teams = append(league.Teams, *team)
		}

	}

	err = lc.leagueService.CreateLeague(league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create league"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "League created and initialized successfully", "league_id": league.ID})
}

// AddTeamToLeague adds a team to a league
// @Summary Add a team to a league
// @Tags League
// @Accept json
// @Produce json
// @Param leagueID path int true "League ID"
// @Param teamID path int true "Team ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /league/add-team/{leagueID}/{teamID} [post]
func (lc *LeagueController) AddTeamToLeague(c *gin.Context) {
	leagueID, err := strconv.ParseUint(c.Param("leagueID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	err = lc.leagueService.AddTeamToLeague(uint(leagueID), uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add team to league: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team added to league successfully"})
}

// RemoveTeamFromLeague removes a team from a league
// @Summary Remove a team from a league
// @Tags League
// @Accept json
// @Produce json
// @Param leagueID path int true "League ID"
// @Param teamID path int true "Team ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /league/remove-team/{leagueID}/{teamID} [post]
func (lc *LeagueController) RemoveTeamFromLeague(c *gin.Context) {
	leagueID, err := strconv.ParseUint(c.Param("leagueID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	err = lc.leagueService.RemoveTeamFromLeague(uint(leagueID), uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove team from league: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team removed from league successfully"})
}

// AdvanceWeek advances the league by one week
// @Summary Advance the league by one week
// @Tags League
// @Accept json
// @Produce json
// @Param leagueID path int true "League ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /league/advance-week/{leagueID} [post]
func (lc *LeagueController) AdvanceWeek(c *gin.Context) {
	leagueID, err := strconv.ParseUint(c.Param("leagueID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}

	err = lc.leagueService.AdvanceWeek(uint(leagueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to advance week: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Week advanced successfully"})
}

// ViewMatchResults returns the match results for the current week
// @Summary View match results for the current week
// @Tags League
// @Accept json
// @Produce json
// @Param leagueID path int true "League ID"
// @Success 200 {object} []models.Match
// @Failure 500 {object} gin.H
// @Router /league/view-matches/{leagueID} [get]
func (lc *LeagueController) ViewMatchResults(c *gin.Context) {
	leagueID, err := strconv.ParseUint(c.Param("leagueID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}

	matches, err := lc.leagueService.ViewMatchResults(uint(leagueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to view match results: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, matches)
}

// EditMatchResults edits the results of a match
// @Summary Edit the results of a match
// @Tags League
// @Accept json
// @Produce json
// @Param matchID path int true "Match ID"
// @Param match body models.Match true "Updated Match"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /league/edit-match/{matchID} [post]
func (lc *LeagueController) EditMatchResults(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("matchID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	var updatedMatch models.Match
	if err := c.ShouldBindJSON(&updatedMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = lc.leagueService.EditMatchResults(uint(matchID), &updatedMatch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit match results: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match results edited successfully"})
}

// PredictChampion predicts the champion of the league
// @Summary Predict the champion of the league
// @Tags League
// @Accept json
// @Produce json
// @Param leagueID path int true "League ID"
// @Success 200 {object} []dto.TeamPrediction
// @Failure 500 {object} gin.H
// @Router /league/predict-champion/{leagueID} [get]
func (lc *LeagueController) PredictChampion(c *gin.Context) {
	leagueID, err := strconv.ParseUint(c.Param("leagueID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}

	predictions, err := lc.leagueService.PredictChampion(uint(leagueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to predict champion: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, predictions)
}

// PlayAllMatches plays all remaining matches in the league
// @Summary Play all remaining matches in the league
// @Tags League
// @Accept json
// @Produce json
// @Param leagueID path int true "League ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /league/play-all-matches/{leagueID} [post]
func (lc *LeagueController) PlayAllMatches(c *gin.Context) {
	leagueID, err := strconv.ParseUint(c.Param("leagueID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}

	err = lc.leagueService.PlayAllMatches(uint(leagueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to play all matches: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All matches played successfully"})
}
