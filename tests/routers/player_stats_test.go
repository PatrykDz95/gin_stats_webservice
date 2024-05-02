package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin/database"
	"gin/entities"
	"gin/routers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup
	database.InitDB("my_test_database")
	// Run the tests
	code := m.Run()
	// Teardown
	err := database.DB.Where("1 = 1").Delete(&entities.PlayerStats{}).Error
	if err != nil {
		fmt.Printf("Failed to delete all player stats: %v\n", err)
	}
	// Teardown
	os.Exit(code)
}

func TestAddPlayerStats_ValidInput_CreatesPlayerStats(t *testing.T) {

	router := gin.Default()
	router.POST("/playerStats", routes.Add)

	playerStats := entities.PlayerStats{
		Name:     "John",
		Surname:  "Johnathan",
		Position: "Forward",
		Age:      25,
	}

	body, _ := json.Marshal(playerStats)
	req, _ := http.NewRequest("POST", "/playerStats", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestAddPlayerStats_InvalidInput_ReturnsBadRequest(t *testing.T) {
	router := gin.Default()
	router.POST("/playerStats", routes.Add)

	playerStats := entities.PlayerStats{
		Name:     "", // Invalid input
		Surname:  "Johnathan",
		Position: "Forward",
		Age:      25,
	}

	body, _ := json.Marshal(playerStats)
	req, _ := http.NewRequest("POST", "/playerStats", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetAll_ReturnsAllPlayers(t *testing.T) {
	router := gin.Default()
	router.GET("/playerStats", routes.GetAll)

	playerStats1 := entities.PlayerStats{Name: "John", Surname: "Doe", Age: 44, Position: "Forward"}
	playerStats2 := entities.PlayerStats{Name: "Joe", Surname: "Doe", Age: 33, Position: "Back"}
	playerStats3 := entities.PlayerStats{Name: "Mark", Surname: "Doe", Age: 23, Position: "Forward"}
	database.DB.Create(&playerStats1)
	database.DB.Create(&playerStats2)
	database.DB.Create(&playerStats3)

	req, _ := http.NewRequest("GET", "/playerStats", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetById_ReturnsPlayer_WhenPlayerExists(t *testing.T) {
	router := gin.Default()
	router.GET("/playerStats/:id", routes.GetById)

	playerStats := entities.PlayerStats{Name: "John", Surname: "Doe", Age: 30, Position: "Forward"}
	database.DB.Create(&playerStats)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/playerStats/%d", playerStats.ID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetById_ReturnsNotFound_WhenPlayerDoesNotExist(t *testing.T) {
	router := gin.Default()
	router.GET("/playerStats/:id", routes.GetById)

	req, _ := http.NewRequest("GET", "/playerStats/999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetByNameAndSurname_ReturnsPlayer_WhenPlayerExists(t *testing.T) {
	router := gin.Default()
	router.GET("/playerStats/:name/:surname", routes.GetByNameAndSurname)

	req, _ := http.NewRequest("GET", "/playerStats/John/Johnathan", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetByNameAndSurname_ReturnsNotFound_WhenPlayerDoesNotExist(t *testing.T) {
	router := gin.Default()
	router.GET("/playerStats/:name/:surname", routes.GetByNameAndSurname)

	req, _ := http.NewRequest("GET", "/playerStats/Nonexistent/Player", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUpdate_ReturnsBadRequest_WhenInvalidJSON(t *testing.T) {
	router := gin.Default()
	router.PUT("/playerStats/:id", routes.Update)

	req, _ := http.NewRequest("PUT", "/playerStats/1", strings.NewReader("invalid json"))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUpdate_ReturnsOK_WhenPlayerStatsUpdatedSuccessfully(t *testing.T) {
	router := gin.Default()
	router.PUT("/playerStats/:id", routes.Update)

	// Create a player stats record for testing
	playerStats := entities.PlayerStats{Name: "John", Surname: "Doe", Age: 30, Position: "Forward"}
	result := database.DB.Create(&playerStats)
	if result.Error != nil {
		t.Fatalf("Failed to create player stats for testing: %v", result.Error)
	}
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/playerStats/%d", playerStats.ID), strings.NewReader(`{"name": "Jane", "surname": "Doe", "age": 25, "position": "Guard"}`))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDelete_ReturnsNotFound_WhenPlayerStatsDoesNotExist(t *testing.T) {
	router := gin.Default()
	router.DELETE("/playerStats/:id", routes.Delete)

	req, _ := http.NewRequest("DELETE", "/playerStats/999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDelete_ReturnsOK_WhenPlayerStatsDeletedSuccessfully(t *testing.T) {
	router := gin.Default()
	router.DELETE("/playerStats/:id", routes.Delete)

	// Create a player stats record for testing
	playerStats := entities.PlayerStats{Name: "John", Surname: "Doe", Age: 30, Position: "Forward"}
	result := database.DB.Create(&playerStats)
	if result.Error != nil {
		t.Fatalf("Failed to create player stats for testing: %v", result.Error)
	}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/playerStats/%d", playerStats.ID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
