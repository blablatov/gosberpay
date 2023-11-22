module github.com/blablatov/mtls-grpc-gateway

go 1.18

replace github.com/blablatov/mtls-grpc-gateway/gw-mtls-proto => ./gw-mtls-proto

replace github.com/blablatov/mtls-grpc-gateway/gw-mockups => ./gw-mockups

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang/mock v1.1.1
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	golang.org/x/oauth2 v0.4.0
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.52.0-dev
)

require (
	cloud.google.com/go/compute/metadata v0.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
