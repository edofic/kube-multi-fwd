//go:generate protoc -I=. --go_out=. --go-grpc_out=. --go-grpc_opt=requireUnimplementedServers=false ./protocol.proto
package fwd
