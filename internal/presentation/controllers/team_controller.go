package controllers

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// TeamController handles team-related requests
type TeamController struct {
	service services.TeamService
}

// NewTeamController creates a new TeamController
func NewTeamController(service services.TeamService) *TeamController {
	return &TeamController{service: service}
}

// AddTeam adds a new team to the database
// @Summary Add a new team
// @Tags Team
// @Accept json
// @Produce json
// @Param team body models.Team true "Team to add"
// @Success 200 {object} models.Team
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /teams [post]
func (ctrl *TeamController) AddTeam(c *gin.Context) {
	var team *models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.service.CreateTeam(team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, team)
}

// GetTeamByID retrieves a team by its ID
// @Summary Get a team by ID
// @Tags Team
// @Produce json
// @Param teamID path int true "Team ID"
// @Success 200 {object} models.Team
// @Failure 404 {object} gin.H
// @Router /teams/{teamID} [get]
func (ctrl *TeamController) GetTeamByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("teamID"))
	team, err := ctrl.service.GetTeamByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

// UpdateTeam updates an existing team
// @Summary Update an existing team
// @Tags Team
// @Accept json
// @Produce json
// @Param teamID path int true "Team ID"
// @Param team body models.Team true "Updated team"
// @Success 200 {object} models.Team
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /teams/{teamID} [put]
func (ctrl *TeamController) UpdateTeam(c *gin.Context) {
	var team *models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(c.Param("teamID"))
	team.ID = uint(id)
	if err := ctrl.service.UpdateTeam(team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, team)
}

// DeleteTeam deletes a team by its ID
// @Summary Delete a team by ID
// @Tags Team
// @Produce json
// @Param teamID path int true "Team ID"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /teams/{teamID} [delete]
func (ctrl *TeamController) DeleteTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("teamID"))
	if err := ctrl.service.DeleteTeam(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted"})
}

// GetAllTeams retrieves all teams
// @Summary Get all teams
// @Tags Team
// @Produce json
// @Success 200 {array} models.Team
// @Failure 500 {object} gin.H
// @Router /teams [get]
func (ctrl *TeamController) GetAllTeams(c *gin.Context) {
	teams, err := ctrl.service.GetAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}
