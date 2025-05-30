package main

import context "context"

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

func (s *ControlGrpcServer) WaitDeviceHandshakeReason(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func (s *ControlGrpcServer) SendDeviceHandshakeEvent(id string, deviceName string, code uint16) error {
	return nil
}
