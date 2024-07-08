package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func setupTestDB() {
	testDB, _ = sql.Open("sqlite3", ":memory:")
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS user_locations (
        username TEXT PRIMARY KEY,
        latitude REAL,
        longitude REAL
    );
    `
	testDB.Exec(createTableQuery)
}
func teardownTestDB() {
	if testDB != nil {
		testDB.Close()
	}
}
func TestUpdateLocation(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	r := gin.Default()
	r.POST("/api/v1/location/update", func(c *gin.Context) {
		var req LocationUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update database
		if err := updateLocation(testDB, req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location in database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "location updated"})
	})
	// Test valid username
	w := httptest.NewRecorder()
	body, _ := json.Marshal(LocationUpdateRequest{
		Username:  "testuser",
		Latitude:  40.7749,
		Longitude: -120.4194,
	})
	req, _ := http.NewRequest("POST", "/api/v1/location/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "location updated")
	// Test invalid username
	w = httptest.NewRecorder()
	body, _ = json.Marshal(LocationUpdateRequest{
		Username:  "test@user",
		Latitude:  40.7749,
		Longitude: -120.4194,
	})
	req, _ = http.NewRequest("POST", "/api/v1/location/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")

	// Test invalid latitude
	w = httptest.NewRecorder()
	body, _ = json.Marshal(LocationUpdateRequest{
		Username:  "testuser",
		Latitude:  100.0, // Invalid latitude
		Longitude: -120.4194,
	})
	req, _ = http.NewRequest("POST", "/api/v1/location/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")

	// Test invalid longitude
	w = httptest.NewRecorder()
	body, _ = json.Marshal(LocationUpdateRequest{
		Username:  "testuser",
		Latitude:  40.7749,
		Longitude: -200.0, // Invalid longitude
	})
	req, _ = http.NewRequest("POST", "/api/v1/location/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestSearchUsers(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	r := gin.Default()
	r.GET("/api/v1/location/search", func(c *gin.Context) {
		var req SearchRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res := searchUsers(testDB, req)
		c.JSON(http.StatusOK, res)
	})

	// Insert some test data
	testDB.Exec("INSERT INTO user_locations (username, latitude, longitude) VALUES (?, ?, ?)",
		"testuser", 37.7749, -122.4194)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/location/search?latitude=37.7749&longitude=-122.4194&radius=1&page=1&size=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
}
