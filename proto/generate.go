package proto

//go:generate protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import --proto_path=third_party --proto_path=. proto/v1/metrics.proto
