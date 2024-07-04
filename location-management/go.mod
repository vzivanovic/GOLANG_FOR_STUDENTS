module location-management

go 1.22.4

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/vzivanovic/GOLANG_FOR_STUDENTS/db v0.0.0
	github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto v0.0.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require github.com/mattn/go-sqlite3 v1.14.11 // indirect

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect; updated version
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto => ../proto

replace github.com/vzivanovic/GOLANG_FOR_STUDENTS/db => ../db
