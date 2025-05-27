package main

// GRPC Server for UNIX sockets or Named Pipes for communication between clipmuxd and local UI

type ControlGrpcServer struct {
	Storage *Storage
	UnimplementedControlServer
}

func NewControlGrpcServer(storage *Storage) *ControlGrpcServer {
	return &ControlGrpcServer{
		Storage: storage,
	}
}
