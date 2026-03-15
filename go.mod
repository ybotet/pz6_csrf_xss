module github.com/ybotet/pz3_logr

go 1.25.1

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/mux v1.8.1
	github.com/sirupsen/logrus v1.9.4
	github.com/ybotet/pz3_logr/gen v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.79.1
)

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// Si tienes código generado en /gen, mantenemos el replace
replace github.com/ybotet/pz3_logr/gen => ./gen
