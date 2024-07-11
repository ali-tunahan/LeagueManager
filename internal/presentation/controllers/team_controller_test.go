package controllers

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func setupRouter() *gin.Engine {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err := db.AutoMigrate(&models.Team{})
	if err != nil {
		return nil
	}
	repo := repositories.NewTeamRepository(db)
	service := services.NewTeamService(repo)
	controller := NewTeamController(service)

	r := gin.Default()
	api := r.Group("/api")
	{
		team := api.Group("/teams")
		team.GET("", controller.GetAllTeams)
		team.POST("", controller.AddTeam)
		team.GET("/:teamID", controller.GetTeamByID)
		team.PUT("/:teamID", controller.UpdateTeam)
		team.DELETE("/:teamID", controller.DeleteTeam)
	}

	return r
}

func TestTeamController(t *testing.T) {
	router := setupRouter()

	// Assert router is not nil since it implies error
	assert.NotNil(t, router)

	// Create Team
	w := httptest.NewRecorder()
	reqBody := `{"name": "Team A", "attack_strength": 80, "defense_strength": 70}`
	req, _ := http.NewRequest("POST", "/api/teams", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdTeam models.Team
	json.Unmarshal(w.Body.Bytes(), &createdTeam)
	assert.NotZero(t, createdTeam.ID)
	assert.Equal(t, "Team A", createdTeam.Name)

	// Get All Teams
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/teams", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var teams []models.Team
	json.Unmarshal(w.Body.Bytes(), &teams)
	assert.Len(t, teams, 1)

	// Get Team by ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/teams/"+strconv.Itoa(int(createdTeam.ID)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var team models.Team
	json.Unmarshal(w.Body.Bytes(), &team)
	assert.Equal(t, "Team A", team.Name)

	// Update Team
	w = httptest.NewRecorder()
	reqBody = `{"name": "Team A Updated", "attack_strength": 85, "defense_strength": 75}`
	req, _ = http.NewRequest("PUT", "/api/teams/"+strconv.Itoa(int(createdTeam.ID)), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get Updated Team by ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/teams/"+strconv.Itoa(int(createdTeam.ID)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &team)
	assert.Equal(t, "Team A Updated", team.Name)

	// Delete Team
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/teams/"+strconv.Itoa(int(createdTeam.ID)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Try to Get Deleted Team by ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/teams/"+strconv.Itoa(int(createdTeam.ID)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
