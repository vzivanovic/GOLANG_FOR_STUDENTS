module location-history

go 1.22.4

require (
	github.com/mattn/go-sqlite3 v1.14.11
	github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto v0.0.0
	google.golang.org/grpc v1.65.0  
	google.golang.org/protobuf v1.34.2 // direct
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect; updated version
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto => ../proto
