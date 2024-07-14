package controllers_test

import (
	"LeagueManager/internal/application/services"
	"LeagueManager/internal/domain/models"
	"LeagueManager/internal/domain/repositories"
	"LeagueManager/internal/presentation/controllers"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func setupTest() (*gorm.DB, *gin.Engine) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	db.AutoMigrate(&models.Team{}, &models.League{}, &models.Match{}, &models.Standing{})

	teamRepo := repositories.NewTeamRepository(db)
	leagueRepo := repositories.NewLeagueRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	standingRepo := repositories.NewStandingRepository(db)

	leagueService := services.NewLeagueService(leagueRepo, teamRepo, matchRepo, standingRepo)
	teamService := services.NewTeamService(teamRepo, leagueRepo)

	leagueController := controllers.NewLeagueController(leagueService, teamService)
	teamController := controllers.NewTeamController(teamService)

	router := gin.Default()

	api := router.Group("/api")
	{
		team := api.Group("/teams")
		team.GET("", teamController.GetAllTeams)
		team.POST("", teamController.AddTeam)
		team.GET("/:teamID", teamController.GetTeamByID)
		team.PUT("/:teamID", teamController.UpdateTeam)
		team.DELETE("/:teamID", teamController.DeleteTeam)

		league := api.Group("/leagues")
		league.POST("/create", leagueController.CreateLeague)
		league.POST("/initialize", leagueController.CreateAndInitializeLeague)
		league.POST("/add-team/:leagueID/:teamID", leagueController.AddTeamToLeague)
		league.POST("/remove-team/:leagueID/:teamID", leagueController.RemoveTeamFromLeague)
		league.POST("/advance-week/:leagueID", leagueController.AdvanceWeek)
		league.GET("/view-matches/:leagueID", leagueController.ViewMatchResults)
		league.POST("/edit-match/:matchID", leagueController.EditMatchResults)
		league.GET("/predict-champion/:leagueID", leagueController.PredictChampion)
		league.POST("/play-all-matches/:leagueID", leagueController.PlayAllMatches)
		league.POST("/start/:leagueID", leagueController.StartLeague)
	}

	return db, router
}

func createLeague(t *testing.T, router *gin.Engine) uint {
	w := httptest.NewRecorder()
	reqBody := `{"name":"Test League"}`
	req, _ := http.NewRequest("POST", "/api/leagues/create", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return uint(response["league_id"].(float64))
}

func createTeam(t *testing.T, router *gin.Engine, teamName string) uint {
	w := httptest.NewRecorder()
	reqBody := `{"name":"` + teamName + `","attack_strength":80,"defense_strength":75}`
	req, _ := http.NewRequest("POST", "/api/teams", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Team
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response.ID
}

func TestCreateLeague(t *testing.T) {
	_, router := setupTest()

	leagueID := createLeague(t, router)
	assert.NotZero(t, leagueID)
}

func TestCreateAndInitializeLeague(t *testing.T) {
	_, router := setupTest()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/leagues/initialize", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "League created and initialized successfully", response["message"])

	leagueID := uint(response["league_id"].(float64))
	assert.NotZero(t, leagueID)
}

func TestAddTeamToLeague(t *testing.T) {
	_, router := setupTest()

	leagueID := createLeague(t, router)
	teamID := createTeam(t, router, "Team E")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Team added to league successfully", response["message"])
}

func TestRemoveTeamFromLeague(t *testing.T) {
	_, router := setupTest()

	leagueID := createLeague(t, router)
	teamID := createTeam(t, router, "Team E")

	// Add team to league
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Remove team from league
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/remove-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Team removed from league successfully", response["message"])
}

func TestAdvanceWeek(t *testing.T) {
	_, router := setupTest()

	leagueID := createLeague(t, router)

	// Add 4 teams
	teamID1 := createTeam(t, router, "Team A")
	teamID2 := createTeam(t, router, "Team B")
	teamID3 := createTeam(t, router, "Team C")
	teamID4 := createTeam(t, router, "Team D")

	// Add teams to league
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID1)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID2)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID3)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID4)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/start/"+strconv.Itoa(int(leagueID)), nil)
	router.ServeHTTP(w, req)
	// print the response
	t.Log(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/advance-week/"+strconv.Itoa(int(leagueID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Week advanced successfully", response["message"])
}

func TestViewMatchResults(t *testing.T) {
	_, router := setupTest()

	leagueID := createLeague(t, router)

	// Add 4 teams
	teamID1 := createTeam(t, router, "Team A")
	teamID2 := createTeam(t, router, "Team B")
	teamID3 := createTeam(t, router, "Team C")
	teamID4 := createTeam(t, router, "Team D")

	// Add teams to league
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID1)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID2)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID3)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID4)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Start league
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/start/"+strconv.Itoa(int(leagueID)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Advance week
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/advance-week/"+strconv.Itoa(int(leagueID)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// View match results
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/leagues/view-matches/"+strconv.Itoa(int(leagueID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var matches []models.Match
	err := json.Unmarshal(w.Body.Bytes(), &matches)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(matches))
}

func TestStartLeague(t *testing.T) {
	_, router := setupTest()

	leagueID := createLeague(t, router)
	teamID1 := createTeam(t, router, "Team A")
	teamID2 := createTeam(t, router, "Team B")

	// Add teams to league
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID1)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID2)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Start league should fail do not have 4 teams
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/start/"+strconv.Itoa(int(leagueID)), nil)
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)

	// Add 2 more teams
	teamID3 := createTeam(t, router, "Team C")
	teamID4 := createTeam(t, router, "Team D")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID3)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/add-team/"+strconv.Itoa(int(leagueID))+"/"+strconv.Itoa(int(teamID4)), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Start league
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/leagues/start/"+strconv.Itoa(int(leagueID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response2 map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response2)
	assert.NoError(t, err)

	assert.Equal(t, "League started successfully", response2["message"])
}
