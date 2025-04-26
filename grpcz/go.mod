module github.com/hakadoriya/z.go/grpcz

go 1.23.4

replace github.com/hakadoriya/z.go => ../../z.go

require google.golang.org/grpc v1.71.1

require github.com/hakadoriya/z.go v0.0.0-00010101000000-000000000000

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)
