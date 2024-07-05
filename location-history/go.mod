module location-history

go 1.22.4

require (
	github.com/mattn/go-sqlite3 v1.14.11
	github.com/vzivanovic/GOLANG_FOR_STUDENTS/db v0.0.0
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

require github.com/stretchr/testify v1.9.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto => ../proto

replace github.com/vzivanovic/GOLANG_FOR_STUDENTS/db => ../db
