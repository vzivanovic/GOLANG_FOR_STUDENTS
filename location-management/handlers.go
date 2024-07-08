package main

import (
	"context"
	"database/sql"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto"
)

type LocationUpdateRequest struct {
	Username  string  `json:"username" binding:"required,min=4,max=16,alphanum"`
	Latitude  float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}

type SearchRequest struct {
	Latitude  float64 `form:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `form:"longitude" binding:"required,gte=-180,lte=180"`
	Radius    float64 `form:"radius" binding:"required"`
	Page      int     `form:"page" binding:"required"`
	Size      int     `form:"size" binding:"required"`
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

func updateLocation(db *sql.DB, req LocationUpdateRequest) error {
	_, err := db.Exec("INSERT OR REPLACE INTO user_locations (username, latitude, longitude) VALUES (?, ?, ?)",
		req.Username, req.Latitude, req.Longitude)
	return err
}

func searchUsers(db *sql.DB, req SearchRequest) SearchResponse {
	rows, err := db.Query("SELECT username, latitude, longitude FROM user_locations")
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

func UpdateLocationHandler(c *gin.Context, grpcHostname string, db *sql.DB) {
	var req LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := grpc.DialContext(context.Background(), grpcHostname+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to location history microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to location history microservice"})
		return
	}
	defer conn.Close()

	client := pb.NewLocationServiceClient(conn)

	_, err = client.UpdateLocation(context.Background(), &pb.LocationUpdate{
		Username:  req.Username,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		log.Printf("Failed to update location in microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location in microservice"})
		return
	}

	if err := updateLocation(db, req); err != nil {
		log.Printf("Failed to update location in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "location updated"})
}

func SearchUsersHandler(c *gin.Context, db *sql.DB) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := searchUsers(db, req)
	c.JSON(http.StatusOK, res)
}

func GetDistanceHandler(c *gin.Context, grpcHostname string, db *sql.DB) {
	var req DistanceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set EndTime to current time if not provided
	if req.EndTime.IsZero() {
		req.EndTime = time.Now()
	}

	conn, err := grpc.DialContext(context.Background(), grpcHostname+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to location history microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to location history microservice"})
		return
	}
	defer conn.Close()

	client := pb.NewLocationServiceClient(conn)

	res, err := client.GetDistance(context.Background(), &pb.DistanceRequest{
		Username:  req.Username,
		StartTime: timestamppb.New(req.StartTime),
		EndTime:   timestamppb.New(req.EndTime),
	})
	if err != nil {
		log.Printf("Failed to calculate distance in microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate distance in microservice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"distance": res.Distance})
}
