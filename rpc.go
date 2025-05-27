package main

// GRPC Server for communicate between clipmuxd services

type CommonGrpcServer struct {
	SslClientKey  string
	SslClientCert string
	Storage       *Storage
	UnimplementedRpcServer
}

func NewCommonGrpcServer(sslClientKey, sslClientCert string, storage *Storage) *CommonGrpcServer {
	return &CommonGrpcServer{
		SslClientKey:  sslClientKey,
		SslClientCert: sslClientCert,
		Storage:       storage,
	}
}
