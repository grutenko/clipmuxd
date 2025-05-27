package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"database/sql"

	"google.golang.org/grpc"

	_ "github.com/mattn/go-sqlite3"
)

func mustOpenStorage(file string) *Storage {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", file))
	if err != nil {
		panic(err)
	}
	return &Storage{db}
}

func main() {
	configFilename := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config := MustLoadConfig(*configFilename)
	storage := mustOpenStorage(config.DatabaseFile)
	defer storage.Conn.Close()

	ctx := AppCtx{
		Config:  config,
		Storage: storage,
	}

	handshakeListen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.HandshakePort))
	if err != nil {
		panic(err)
	}
	handshakeGrpcServer := grpc.NewServer()
	handshakeHandler := NewHandshakeGrpcServer(&ctx, nil)
	RegisterHandshakeServer(handshakeGrpcServer, handshakeHandler)

	go func() {
		log.Printf("Handshake server started :%d", config.HandshakePort)
		if err := handshakeGrpcServer.Serve(handshakeListen); err != nil {
			panic(err)
		}
	}()

	if err := os.RemoveAll(config.Socket); err != nil {
		log.Fatalf("failed to remove existing socket: %v", err)
	}
	controlListen, err := net.Listen("unix", config.Socket)
	if err != nil {
		panic(err)
	}
	controlGrpcServer := grpc.NewServer()
	controlHandler := NewControlGrpcServer(storage)
	RegisterControlServer(controlGrpcServer, controlHandler)

	go func() {
		log.Printf("Control server started %s", config.Socket)
		if err := controlGrpcServer.Serve(controlListen); err != nil {
			panic(err)
		}
	}()

	select {}
}
