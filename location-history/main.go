package main

import (
	"context"
	"database/sql"
	"log"
	"math"
	"net"

	_ "github.com/mattn/go-sqlite3"
	pb "github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedLocationServiceServer
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./location_history.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func (s *server) UpdateLocation(ctx context.Context, req *pb.LocationUpdate) (*emptypb.Empty, error) {
	_, err := db.Exec("INSERT INTO location_history (username, latitude, longitude) VALUES (?, ?, ?)",
		req.Username, req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}
	log.Printf("Received location update: %v", req)
	return &emptypb.Empty{}, nil
}

func (s *server) GetDistance(ctx context.Context, req *pb.DistanceRequest) (*pb.DistanceResponse, error) {
	rows, err := db.Query("SELECT latitude, longitude FROM location_history WHERE username = ? AND timestamp BETWEEN ? AND ? ORDER BY timestamp",
		req.Username, req.StartTime.AsTime(), req.EndTime.AsTime())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalDistance float64
	var prevLat, prevLon float64
	first := true

	for rows.Next() {
		var lat, lon float64
		if err := rows.Scan(&lat, &lon); err != nil {
			return nil, err
		}

		if first {
			first = false
		} else {
			totalDistance += distance(prevLat, prevLon, lat, lon)
		}

		prevLat = lat
		prevLon = lon
	}

	return &pb.DistanceResponse{Distance: totalDistance}, nil
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
	initDB()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLocationServiceServer(s, &server{})
	reflection.Register(s)

	log.Println("Starting location history microservice on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
