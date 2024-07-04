package main

import (
	"flag"
	"log"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vzivanovic/GOLANG_FOR_STUDENTS/db"
)

var grpcHostname string

type LocationUpdateRequest struct {
	Username  string  `json:"username" binding:"required,min=4,max=16,alphanum"`
	Latitude  float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}

type SearchRequest struct {
	Latitude  float64 `form:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `form:"longitude" binding:"required,gte=-180,lte=180"`
	Radius    float64 `form:"radius" binding:"required"`
	Page      int     `form:"page" binding:"required,default=1"`
	Size      int     `form:"size" binding:"required,default=10"`
}

type DistanceRequest struct {
	Username  string    `form:"username" binding:"required,min=4,max=16,alphanum"`
	StartTime time.Time `form:"start_time" binding:"required"`
	EndTime   time.Time `form:"end_time"`
}

type UserLocation struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SearchResponse struct {
	Users      []UserLocation `json:"users"`
	Page       int            `json:"page"`
	TotalPages int            `json:"total_pages"`
	TotalUsers int            `json:"total_users"`
}

func updateLocation(req LocationUpdateRequest) error {
	_, err := db.DB.Exec("INSERT OR REPLACE INTO user_locations (username, latitude, longitude) VALUES (?, ?, ?)",
		req.Username, req.Latitude, req.Longitude)
	return err
}

func searchUsers(req SearchRequest) SearchResponse {
	rows, err := db.DB.Query("SELECT username, latitude, longitude FROM user_locations")
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	var users []UserLocation
	for rows.Next() {
		var user UserLocation
		if err := rows.Scan(&user.Username, &user.Latitude, &user.Longitude); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		if distance(req.Latitude, req.Longitude, user.Latitude, user.Longitude) <= req.Radius {
			users = append(users, user)
		}
	}

	totalUsers := len(users)
	totalPages := int(math.Ceil(float64(totalUsers) / float64(req.Size)))
	start := (req.Page - 1) * req.Size
	end := start + req.Size
	if end > totalUsers {
		end = totalUsers
	}
	paginatedUsers := users[start:end]

	return SearchResponse{
		Users:      paginatedUsers,
		Page:       req.Page,
		TotalPages: totalPages,
		TotalUsers: totalUsers,
	}
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func main() {
	flag.StringVar(&grpcHostname, "grpc-hostname", "localhost", "gRPC server hostname")
	flag.Parse()

	db.InitLocationDB()

	r := gin.Default()

	r.POST("/api/v1/location/update", func(c *gin.Context) {
		UpdateLocationHandler(c, grpcHostname)
	})
	r.GET("/api/v1/location/search", SearchUsersHandler)
	r.GET("/api/v1/location/distance", func(c *gin.Context) {
		GetDistanceHandler(c, grpcHostname)
	})

	r.Run(":8080")
}
