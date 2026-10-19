module github.com/DoomGuy1818/rocket-science/payment

go 1.27.1

replace github.com/DoomGuy1818/rocket-science/shared => ../shared

require (
	github.com/DoomGuy1818/rocket-science/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.12.1
	google.golang.org/grpc v1.84.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.30.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
