package controllers

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TeamController struct {
	service services.TeamService
}

func NewTeamController(service services.TeamService) *TeamController {
	return &TeamController{service: service}
}

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

func (ctrl *TeamController) GetTeamByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("teamID"))
	team, err := ctrl.service.GetTeamByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

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

func (ctrl *TeamController) DeleteTeam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("teamID"))
	if err := ctrl.service.DeleteTeam(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted"})
}

func (ctrl *TeamController) GetAllTeams(c *gin.Context) {
	teams, err := ctrl.service.GetAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}
