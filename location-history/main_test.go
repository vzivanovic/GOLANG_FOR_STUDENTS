package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	pb "github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Ensure db is available in the test scope
var testDB *sql.DB

func setupTestDB() {
	var err error
	testDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS location_history (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT,
        latitude REAL,
        longitude REAL,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	testDB.Exec(createTableQuery)
}

func TestUpdateLocation(t *testing.T) {
	setupTestDB()

	s := &server{}
	s.db = testDB // assign testDB to the server's db field

	_, err := s.UpdateLocation(context.Background(), &pb.LocationUpdate{
		Username:  "testuser",
		Latitude:  37.7749,
		Longitude: -122.4194,
	})

	assert.NoError(t, err)

	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM location_history WHERE username = ?", "testuser").Scan(&count)

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGetDistance(t *testing.T) {
	setupTestDB()

	// Insert some test data
	testDB.Exec("INSERT INTO location_history (username, latitude, longitude, timestamp) VALUES (?, ?, ?, ?)",
		"testuser", 37.7749, -122.4194, "2023-01-01T00:00:00Z")
	testDB.Exec("INSERT INTO location_history (username, latitude, longitude, timestamp) VALUES (?, ?, ?, ?)",
		"testuser", 37.7750, -122.4195, "2023-01-01T01:00:00Z")

	s := &server{}
	s.db = testDB // assign testDB to the server's db field

	resp, err := s.GetDistance(context.Background(), &pb.DistanceRequest{
		Username:  "testuser",
		StartTime: timestamppb.New(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)),
		EndTime:   timestamppb.New(time.Date(2023, time.January, 1, 2, 0, 0, 0, time.UTC)),
	})

	assert.NoError(t, err)
	assert.Greater(t, resp.Distance, 0.0)
}
