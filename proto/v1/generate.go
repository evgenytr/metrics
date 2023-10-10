package api

//go:generate protoc --go_out=../../pkg/api/v1/ --go_opt=paths=source_relative --go-grpc_out=../../pkg/api/v1/ --go-grpc_opt=paths=source_relative --proto_path=../../third_party --proto_path=. metrics.proto
