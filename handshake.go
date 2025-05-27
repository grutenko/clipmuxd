package main

// GRPC server for handshake service

// HandshakeGrpcServer implements the HandshakeServiceServer interface.
type HandshakeGrpcServer struct {
	Ctx               *AppCtx
	ControlGrpcServer *ControlGrpcServer
	UnimplementedHandshakeServer
}

func NewHandshakeGrpcServer(ctx *AppCtx, control *ControlGrpcServer) HandshakeGrpcServer {
	return HandshakeGrpcServer{
		Ctx:               ctx,
		ControlGrpcServer: control,
	}
}
